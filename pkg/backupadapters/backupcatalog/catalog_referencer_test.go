package backupcatalog

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestCatalogReferencerAdapter_Reference(t *testing.T) {
	owner1 := ownermodel.Owner("owner1")
	jan24 := time.Date(2021, time.January, 24, 0, 0, 0, 0, time.UTC)

	type fields struct {
		Owner                   ownermodel.Owner
		FindSignaturePort       InsertMediaSimulator
		StatefulAlbumReferencer StatefulAlbumReferencer
	}
	type args struct {
		medias []*backup.AnalysedMedia
	}
	signature1 := catalog.MediaSignature{
		SignatureSha256: "sha256-1",
		SignatureSize:   12,
	}
	mediaReference1 := catalog.MediaFutureReference{
		Signature:          signature1,
		ProvisionalMediaId: "uuid-01",
		AlreadyExists:      false,
	}
	mediaReference1Exists := catalog.MediaFutureReference{
		Signature:          signature1,
		ProvisionalMediaId: "sha256-1#12",
		AlreadyExists:      true,
	}
	analysedMedia1 := backup.AnalysedMedia{
		Sha256Hash: signature1.SignatureSha256,
		FoundMedia: backup.NewInMemoryMedia("file1", jan24, []byte(strings.Repeat("a", signature1.SignatureSize))),
		Details:    &backup.MediaDetails{DateTime: jan24},
	}
	analysedMedia2 := backup.AnalysedMedia{
		Sha256Hash: signature1.SignatureSha256,
		FoundMedia: backup.NewInMemoryMedia("file2", jan24, []byte(strings.Repeat("a", signature1.SignatureSize))),
		Details:    &backup.MediaDetails{DateTime: jan24},
	}
	albumId1 := catalog.AlbumId{Owner: owner1, FolderName: catalog.NewFolderName("/jan24")}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []backup.BackingUpMediaRequest
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return an empty map when no media are provided",
			fields: fields{
				Owner: owner1,
				FindSignaturePort: FindSignaturePortFake{
					owner1: {},
				},
				StatefulAlbumReferencer: make(StatefulAlbumReferencerFake),
			},
			args: args{
				medias: nil,
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should return a map with the reference of a media that has been found, and the album is not found and not created",
			fields: fields{
				Owner: owner1,
				FindSignaturePort: FindSignaturePortFake{
					owner1: {
						mediaReference1Exists,
					},
				},
				StatefulAlbumReferencer: StatefulAlbumReferencerFake{
					jan24: {
						AlbumId:          nil,
						AlbumJustCreated: false,
					},
				},
			},
			args: args{
				medias: []*backup.AnalysedMedia{
					&analysedMedia1,
				},
			},
			want: []backup.BackingUpMediaRequest{
				{
					AnalysedMedia: &analysedMedia1,
					CatalogReference: ReferenceSnapshot{
						ExistsValue:           true,
						AlbumCreatedValue:     false,
						AlbumFolderNameValue:  "",
						UniqueIdentifierValue: mediaReference1Exists.Signature.Value(),
						MediaIdValue:          mediaReference1Exists.ProvisionalMediaId.Value(),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should returns twice the analysed media even if it shares the same signature",
			fields: fields{
				Owner: owner1,
				FindSignaturePort: FindSignaturePortFake{
					owner1: {
						mediaReference1Exists,
					},
				},
				StatefulAlbumReferencer: StatefulAlbumReferencerFake{
					jan24: {
						AlbumId:          nil,
						AlbumJustCreated: false,
					},
				},
			},
			args: args{
				medias: []*backup.AnalysedMedia{
					&analysedMedia1,
					&analysedMedia2,
				},
			},
			want: []backup.BackingUpMediaRequest{
				{
					AnalysedMedia: &analysedMedia1,
					CatalogReference: ReferenceSnapshot{
						ExistsValue:           true,
						AlbumCreatedValue:     false,
						AlbumFolderNameValue:  "",
						UniqueIdentifierValue: mediaReference1Exists.Signature.Value(),
						MediaIdValue:          mediaReference1Exists.ProvisionalMediaId.Value(),
					},
				},
				{
					AnalysedMedia: &analysedMedia2,
					CatalogReference: ReferenceSnapshot{
						ExistsValue:           true,
						AlbumCreatedValue:     false,
						AlbumFolderNameValue:  "",
						UniqueIdentifierValue: mediaReference1Exists.Signature.Value(),
						MediaIdValue:          mediaReference1Exists.ProvisionalMediaId.Value(),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return a map with the reference of a media that has not be found, and the album has been created",
			fields: fields{
				Owner: owner1,
				FindSignaturePort: FindSignaturePortFake{
					owner1: {
						mediaReference1,
					},
				},
				StatefulAlbumReferencer: StatefulAlbumReferencerFake{
					jan24: {
						AlbumId:          &albumId1,
						AlbumJustCreated: true,
					},
				},
			},
			args: args{
				medias: []*backup.AnalysedMedia{
					&analysedMedia1,
				},
			},
			want: []backup.BackingUpMediaRequest{
				{
					AnalysedMedia: &analysedMedia1,
					CatalogReference: ReferenceSnapshot{
						ExistsValue:           false,
						AlbumCreatedValue:     true,
						AlbumFolderNameValue:  albumId1.FolderName.String(),
						UniqueIdentifierValue: mediaReference1Exists.Signature.Value(),
						MediaIdValue:          mediaReference1Exists.ProvisionalMediaId.Value(),
					},
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := new(CataloguerObserverFake)

			c := &CatalogReferencerAdapter{
				Owner:                   tt.fields.Owner,
				InsertMediaSimulator:    tt.fields.FindSignaturePort,
				StatefulAlbumReferencer: tt.fields.StatefulAlbumReferencer,
			}
			err := c.Reference(context.Background(), tt.args.medias, observer)
			if !tt.wantErr(t, err) {
				return
			}

			var reallyGot []backup.BackingUpMediaRequest
			for _, request := range observer.Got {
				reallyGot = append(reallyGot, backup.BackingUpMediaRequest{
					AnalysedMedia: request.AnalysedMedia,
					CatalogReference: ReferenceSnapshot{
						ExistsValue:           request.CatalogReference.Exists(),
						AlbumCreatedValue:     request.CatalogReference.AlbumCreated(),
						AlbumFolderNameValue:  request.CatalogReference.AlbumFolderName(),
						UniqueIdentifierValue: mediaReference1Exists.Signature.Value(),
						MediaIdValue:          mediaReference1Exists.ProvisionalMediaId.Value(),
					},
				})
			}

			assert.Equalf(t, tt.want, reallyGot, "Catalog() = %v, want %v", observer.Got, tt.want)
		})
	}
}

type FindSignaturePortFake map[ownermodel.Owner][]catalog.MediaFutureReference

func (f FindSignaturePortFake) SimulateInsertingMedia(ctx context.Context, owner ownermodel.Owner, signatures []catalog.MediaSignature) ([]catalog.MediaFutureReference, error) {
	signaturesCopy := make(map[catalog.MediaSignature]any)
	for _, signature := range signatures {
		if _, duplicate := signaturesCopy[signature]; duplicate {
			return nil, errors.Errorf("invalid request, duplicate signature %s", signature)
		}
		signaturesCopy[signature] = nil
	}

	var foundSignatures []catalog.MediaFutureReference
	for _, reference := range f[owner] {
		if slices.Contains(signatures, reference.Signature) {
			foundSignatures = append(foundSignatures, reference)
		}
	}

	return foundSignatures, nil
}

type ReferenceSnapshot struct {
	ExistsValue           bool
	AlbumCreatedValue     bool
	AlbumFolderNameValue  string
	UniqueIdentifierValue string
	MediaIdValue          string
}

func (r ReferenceSnapshot) MediaId() string {
	return r.MediaIdValue
}

func (r ReferenceSnapshot) Exists() bool {
	return r.ExistsValue
}

func (r ReferenceSnapshot) AlbumCreated() bool {
	return r.AlbumCreatedValue
}

func (r ReferenceSnapshot) AlbumFolderName() string {
	return r.AlbumFolderNameValue
}

func (r ReferenceSnapshot) UniqueIdentifier() string {
	return r.UniqueIdentifierValue
}

type StatefulAlbumReferencerFake map[time.Time]catalog.AlbumReference

func (s StatefulAlbumReferencerFake) FindReference(ctx context.Context, mediaTime time.Time) (catalog.AlbumReference, error) {
	return s[mediaTime], nil
}

type CataloguerObserverFake struct {
	Got []backup.BackingUpMediaRequest
}

func (c *CataloguerObserverFake) OnMediaCatalogued(ctx context.Context, requests []backup.BackingUpMediaRequest) error {
	c.Got = append(c.Got, requests...)
	return nil
}
