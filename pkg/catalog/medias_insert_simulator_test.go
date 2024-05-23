package catalog

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"testing"
)

func TestMediasInsertSimulator_SimulateInsertingMedia(t *testing.T) {
	owner1 := ownermodel.Owner("owner-1")
	signature1 := MediaSignature{SignatureSha256: "dbd318c1c462aee872f41109a4dfd3048871a03dedd0fe0e757ced57dad6f2d7", SignatureSize: 42}
	generatedMediaId := MediaId("29MYwcRiruhy9BEJpN_TBIhxoD3t0P4OdXztV9rW8tcq")
	existingMediaId := MediaId("existing-media-id")

	type fields struct {
		FindExistingSignaturePort FindExistingSignaturePort
	}
	type args struct {
		owner      ownermodel.Owner
		signatures []MediaSignature
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []MediaFutureReference
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return an empty list if no media are requested.",
			fields: fields{
				FindExistingSignaturePort: make(FindExistingSignaturePortFake),
			},
			args: args{
				owner:      owner1,
				signatures: nil,
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should return a not-found signature if media doesn't already exists.",
			fields: fields{
				FindExistingSignaturePort: FindExistingSignaturePortFake{
					owner1: nil,
				},
			},
			args: args{
				owner:      owner1,
				signatures: []MediaSignature{signature1},
			},
			want: []MediaFutureReference{
				{Signature: signature1, ProvisionalMediaId: generatedMediaId, AlreadyExists: false},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return a found signature if media already exists.",
			fields: fields{
				FindExistingSignaturePort: FindExistingSignaturePortFake{
					owner1: map[MediaSignature]MediaId{
						signature1: existingMediaId,
					},
				},
			},
			args: args{
				owner:      owner1,
				signatures: []MediaSignature{signature1},
			},
			want: []MediaFutureReference{
				{Signature: signature1, ProvisionalMediaId: existingMediaId, AlreadyExists: true},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return duplicates if duplicates signatures are requested (both backup adapter and dynamo adapter are deduplicating...)",
			fields: fields{
				FindExistingSignaturePort: FindExistingSignaturePortFake{
					owner1: map[MediaSignature]MediaId{
						signature1: existingMediaId,
					},
				},
			},
			args: args{
				owner:      owner1,
				signatures: []MediaSignature{signature1, signature1},
			},
			want: []MediaFutureReference{
				{Signature: signature1, ProvisionalMediaId: existingMediaId, AlreadyExists: true},
				{Signature: signature1, ProvisionalMediaId: existingMediaId, AlreadyExists: true},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MediasInsertSimulator{
				FindExistingSignaturePort: tt.fields.FindExistingSignaturePort,
			}

			got, err := m.SimulateInsertingMedia(context.Background(), tt.args.owner, tt.args.signatures)
			if !tt.wantErr(t, err, fmt.Sprintf("SimulateInsertingMedia(%v, %v, %v)", context.Background(), tt.args.owner, tt.args.signatures)) {
				return
			}
			assert.Equalf(t, tt.want, got, "SimulateInsertingMedia(%v, %v, %v)", context.Background(), tt.args.owner, tt.args.signatures)
		})
	}
}

type FindExistingSignaturePortFake map[ownermodel.Owner]map[MediaSignature]MediaId

func (f FindExistingSignaturePortFake) FindSignatures(ctx context.Context, owner ownermodel.Owner, signatures []MediaSignature) (map[MediaSignature]MediaId, error) {
	result := make(map[MediaSignature]MediaId)

	ownerMedias, stubbed := f[owner]
	if !stubbed {
		return result, nil
	}

	for _, signature := range signatures {
		if mediaId, exists := ownerMedias[signature]; exists {
			result[signature] = mediaId
		}
	}

	return result, nil
}
