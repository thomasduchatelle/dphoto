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

func TestAlbumView_AlbumDeleted(t *testing.T) {
	tonyOwner := ownermodel.Owner("tony")
	ownerUserId := usermodel.UserId("ironman@avenger.hero")
	visitor1UserId := usermodel.UserId("pepper@stark.com")
	visitor2UserId := usermodel.UserId("wanda@avenger.hero")

	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	mar24 := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)

	deletedAlbumId := catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("deleted")}
	survivingAlbumId := catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("surviving")}
	destinationAlbumId := catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("destination")}

	deletedAlbumRow := func(availability Availability) UserAlbumSummary {
		return UserAlbumSummary{
			AlbumSummary: AlbumSummary{AlbumId: deletedAlbumId, Name: "Deleted", Start: jan24, End: feb24, MediaCount: 1},
			Availability: availability,
		}
	}
	survivingAlbumRow := UserAlbumSummary{
		AlbumSummary: AlbumSummary{AlbumId: survivingAlbumId, Name: "Surviving", Start: mar24, End: apr24, MediaCount: 5},
		Availability: OwnerAvailability(ownerUserId),
	}
	destinationAlbumRow := UserAlbumSummary{
		AlbumSummary: AlbumSummary{AlbumId: destinationAlbumId, Name: "Destination", Start: feb24, End: mar24, MediaCount: 2},
		Availability: OwnerAvailability(ownerUserId),
	}

	mediaId1 := catalog.MediaId("media-1")
	mediaId2 := catalog.MediaId("media-2")
	mediaId3 := catalog.MediaId("media-3")

	type fields struct {
		Repository                    *AlbumSummaryInMemoryRepository
		MediaCounterPort              MediaCounterPort
		ListUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
	}
	tests := []struct {
		name          string
		fields        fields
		event         catalog.AlbumDeleted
		expectStored  []UserAlbumSummary
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "it should remove the owner row",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{deletedAlbumRow(OwnerAvailability(ownerUserId))},
				},
				MediaCounterPort:              MediaCounterPortFake(nil),
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(nil),
			},
			event:        catalog.AlbumDeleted{DeletedAlbumId: deletedAlbumId, TransferredMedias: catalog.NewTransferredMedias()},
			expectStored: nil,
			wantErr:      assert.NoError,
		},
		{
			name: "it should remove all visitor rows",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						deletedAlbumRow(OwnerAvailability(ownerUserId)),
						deletedAlbumRow(VisitorAvailability(visitor1UserId)),
						deletedAlbumRow(VisitorAvailability(visitor2UserId)),
					},
				},
				MediaCounterPort:              MediaCounterPortFake(nil),
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(nil),
			},
			event:        catalog.AlbumDeleted{DeletedAlbumId: deletedAlbumId, TransferredMedias: catalog.NewTransferredMedias()},
			expectStored: nil,
			wantErr:      assert.NoError,
		},
		{
			name: "it should leave other albums untouched",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						deletedAlbumRow(OwnerAvailability(ownerUserId)),
						survivingAlbumRow,
					},
				},
				MediaCounterPort:              MediaCounterPortFake(nil),
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(nil),
			},
			event:        catalog.AlbumDeleted{DeletedAlbumId: deletedAlbumId, TransferredMedias: catalog.NewTransferredMedias()},
			expectStored: []UserAlbumSummary{survivingAlbumRow},
			wantErr:      assert.NoError,
		},
		{
			name: "it should recount transferred destination albums",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						deletedAlbumRow(OwnerAvailability(ownerUserId)),
						destinationAlbumRow,
					},
				},
				MediaCounterPort:              MediaCounterPortFake{destinationAlbumId: 5},
				ListUserWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(map[catalog.AlbumId][]Availability{destinationAlbumId: {OwnerAvailability(ownerUserId)}}),
			},
			event: catalog.AlbumDeleted{
				DeletedAlbumId: deletedAlbumId,
				TransferredMedias: catalog.TransferredMedias{
					Transfers: map[catalog.AlbumId][]catalog.MediaId{destinationAlbumId: {mediaId1, mediaId2, mediaId3}},
				},
			},
			expectStored: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: destinationAlbumId, Name: "Destination", Start: feb24, End: mar24, MediaCount: 5},
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
				tt.fields.MediaCounterPort,
				FindAlbumsByIdsFunc(func(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error) { return nil, nil }),
				tt.fields.ListUserWhoCanAccessAlbumPort,
			)

			err := albumView.AlbumDeleted(context.Background(), tt.event)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectStored, tt.fields.Repository.Summaries)
			}
		})
	}
}
