package catalogviews

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func TestAlbumView_AlbumShared(t *testing.T) {
	tonyOwner := ownermodel.Owner("tony")
	ownerUserId := usermodel.UserId("ironman@avenger.hero")
	visitorUserId := usermodel.UserId("pepper@stark.com")
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	weddingsId := catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("weddings")}
	weddingsAlbum := catalog.Album{
		AlbumId: weddingsId,
		Name:    "Weddings",
		Start:   jan24,
		End:     feb24,
	}

	type fields struct {
		Repository       *AlbumSummaryInMemoryRepository
		MediaCounterPort MediaCounterPort
	}
	type args struct {
		album  catalog.Album
		userId usermodel.UserId
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		expectSummaries   []UserAlbumSummary
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "it should add a visitor row with display fields and count",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{
							AlbumSummary: AlbumSummary{AlbumId: weddingsId, Name: "Weddings", Start: jan24, End: feb24, MediaCount: 7},
							Availability: OwnerAvailability(ownerUserId),
						},
					},
				},
				MediaCounterPort: MediaCounterPortFake{weddingsId: 7},
			},
			args: args{album: weddingsAlbum, userId: visitorUserId},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: weddingsId, Name: "Weddings", Start: jan24, End: feb24, MediaCount: 7},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: weddingsId, Name: "Weddings", Start: jan24, End: feb24, MediaCount: 7},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			albumView := NewAlbumView(
				tt.fields.Repository,
				GetAlbumSharingGridFunc(func(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) { return nil, nil }),
				tt.fields.MediaCounterPort,
				FindAlbumsByIdsFunc(func(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error) { return nil, nil }),
			)

			err := albumView.AlbumShared(context.Background(), tt.args.album, tt.args.userId)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.expectSummaries, tt.fields.Repository.Summaries)
			}
		})
	}
}

func TestAlbumView_AlbumUnshared(t *testing.T) {
	tonyOwner := ownermodel.Owner("tony")
	ownerUserId := usermodel.UserId("ironman@avenger.hero")
	visitorUserId := usermodel.UserId("pepper@stark.com")
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	weddingsId := catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("weddings")}

	type fields struct {
		Repository *AlbumSummaryInMemoryRepository
	}
	type args struct {
		albumId catalog.AlbumId
		userId  usermodel.UserId
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		expectSummaries []UserAlbumSummary
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "it should remove the visitor row only",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{
							AlbumSummary: AlbumSummary{AlbumId: weddingsId, Name: "Weddings", Start: jan24, End: feb24, MediaCount: 7},
							Availability: OwnerAvailability(ownerUserId),
						},
						{
							AlbumSummary: AlbumSummary{AlbumId: weddingsId, Name: "Weddings", Start: jan24, End: feb24, MediaCount: 7},
							Availability: VisitorAvailability(visitorUserId),
						},
					},
				},
			},
			args: args{albumId: weddingsId, userId: visitorUserId},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: weddingsId, Name: "Weddings", Start: jan24, End: feb24, MediaCount: 7},
					Availability: OwnerAvailability(ownerUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			albumView := NewAlbumView(
				tt.fields.Repository,
				GetAlbumSharingGridFunc(func(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) { return nil, nil }),
				MediaCounterPortFake(nil),
				FindAlbumsByIdsFunc(func(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error) { return nil, nil }),
			)

			err := albumView.AlbumUnshared(context.Background(), tt.args.albumId, tt.args.userId)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.expectSummaries, tt.fields.Repository.Summaries)
			}
		})
	}
}
