package catalogdynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"testing"
	"time"
)

func TestRepository_TransferMediasFromRecords(t *testing.T) {
	dyn := dynamotestutils.NewTestContext(context.Background(), t)
	album01 := catalog.AlbumId{
		Owner:      "ironman",
		FolderName: catalog.NewFolderName("/my-album-01"),
	}
	album02 := catalog.AlbumId{
		Owner:      "ironman",
		FolderName: catalog.NewFolderName("/my-album-02"),
	}
	nothingTransferred := catalog.NewTransferredMedias()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	middle := time.Date(2024, 1, 2, 12, 2, 42, 0, time.UTC)
	end := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	media01Id := catalog.MediaId("media-01")

	type args struct {
		records catalog.MediaTransferRecords
	}
	tests := []struct {
		name    string
		args    args
		before  []map[string]types.AttributeValue
		after   []map[string]types.AttributeValue
		want    catalog.TransferredMedias
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should not transfer any media if none are found",
			args: args{
				records: catalog.MediaTransferRecords{
					album01: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{album02},
							Start:      start,
							End:        end,
						},
					},
				},
			},
			want:    nothingTransferred,
			wantErr: assert.NoError,
		},
		{
			name: "it should transfer a media from 1 album to the other",
			args: args{
				records: catalog.MediaTransferRecords{
					album01: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{album02},
							Start:      start,
							End:        end,
						},
					},
				},
			},
			before: []map[string]types.AttributeValue{
				mediaAttributeMap(album02, middle, media01Id),
			},
			after: []map[string]types.AttributeValue{
				mediaAttributeMap(album01, middle, media01Id),
			},
			want: catalog.TransferredMedias{
				Transfers: map[catalog.AlbumId][]catalog.MediaId{
					album01: {media01Id},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not transfer from the wrong album",
			args: args{
				records: catalog.MediaTransferRecords{
					album01: []catalog.MediaSelector{
						{
							FromAlbums: []catalog.AlbumId{album01},
							Start:      start,
							End:        end,
						},
					},
				},
			},
			before: []map[string]types.AttributeValue{
				mediaAttributeMap(album02, middle, media01Id),
			},
			after: []map[string]types.AttributeValue{
				mediaAttributeMap(album02, middle, media01Id),
			},
			want:    nothingTransferred,
			wantErr: assert.NoError,
		},
		{
			name: "edge case - it should not fail if no album is selected",
			args: args{
				records: catalog.MediaTransferRecords{
					album01: []catalog.MediaSelector{
						{
							FromAlbums: nil,
							Start:      start,
							End:        end,
						},
					},
				},
			},
			want:    nothingTransferred,
			wantErr: assert.NoError,
		},
		{
			name: "edge case - it should not fail if no transfer records are requested",
			args: args{
				records: nil,
			},
			want:    nothingTransferred,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn := dyn.Subtest(t)

			err := dyn.WithDbContent(dyn.Ctx, tt.before)
			if !assert.NoError(t, err) {
				return
			}

			r := &Repository{
				client: dyn.Client,
				table:  dyn.Table,
			}

			got, err := r.TransferMediasFromRecords(dyn.Ctx, tt.args.records)

			if !tt.wantErr(t, err, fmt.Sprintf("TransferMediasFromRecords(%v)", tt.args.records)) {
				return
			}
			assert.Equalf(t, tt.want, got, "TransferMediasFromRecords(%v)", tt.args.records)

			gotAfter, err := dyn.Got(dyn.Ctx)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.after, gotAfter)
			}
		})
	}
}

func mediaAttributeMap(albumId catalog.AlbumId, dateTime time.Time, mediaId catalog.MediaId) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"PK":           &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#MEDIA#%s", albumId.Owner, mediaId)},
		"SK":           &types.AttributeValueMemberS{Value: "#METADATA"},
		"AlbumIndexPK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s", albumId.Owner, albumId.FolderName)},
		"AlbumIndexSK": &types.AttributeValueMemberS{Value: fmt.Sprintf("MEDIA#%s#%s", dateTime.Format(time.DateTime), mediaId)},
		"Id":           &types.AttributeValueMemberS{Value: string(mediaId)},
		"dateTime":     &types.AttributeValueMemberS{Value: dateTime.Format(time.RFC3339Nano)},
	}
}
