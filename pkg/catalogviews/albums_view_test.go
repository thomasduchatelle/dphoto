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

// The removed TestNewAlbumViewAcceptance covered the owned+shared provider aggregation,
// a split that no longer exists: ListAlbums is served from a single ListSummariesForUser
// query. The five TestAlbumView_ListAlbums scenarios below replace it.

func TestAlbumView_ListAlbums(t *testing.T) {
	tonyOwner := ownermodel.Owner("tony")
	pepperOwner := ownermodel.Owner("pepper")
	visitorUserId := usermodel.UserId("pepper@stark.com")
	ironmanCurrentUser := usermodel.CurrentUser{
		UserId: "ironman@avenger.hero",
		Owner:  &tonyOwner,
	}
	visitorOnlyUser := usermodel.CurrentUser{
		UserId: "wanda@avenger.hero",
	}
	noFilter := ListAlbumsFilter{}
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	mar24 := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)
	ownedAlbum1Id := catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("album-1")}
	ownedAlbum2Id := catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("album-2")}
	sharedAlbum3Id := catalog.AlbumId{Owner: pepperOwner, FolderName: catalog.NewFolderName("album-3")}

	ownerSummariesRepository := func() *AlbumSummaryInMemoryRepository {
		return &AlbumSummaryInMemoryRepository{
			Summaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: ownedAlbum1Id, Name: "Album One", Start: jan24, End: feb24, MediaCount: 1},
					Availability: OwnerAvailability(ironmanCurrentUser.UserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: ownedAlbum2Id, Name: "Album Two", Start: feb24, End: mar24, MediaCount: 2},
					Availability: OwnerAvailability(ironmanCurrentUser.UserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: sharedAlbum3Id, Name: "Album Three", Start: mar24, End: apr24, MediaCount: 3},
					Availability: VisitorAvailability(ironmanCurrentUser.UserId),
				},
			},
		}
	}

	ownerSharingGrid := GetAlbumSharingGridFunc(func(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) {
		if owner == tonyOwner {
			return map[catalog.AlbumId][]usermodel.UserId{
				ownedAlbum1Id: {visitorUserId},
			}, nil
		}
		return nil, nil
	})

	emptySharingGrid := GetAlbumSharingGridFunc(func(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) {
		return nil, nil
	})

	type fields struct {
		Repository              AlbumSummaryRepository
		GetAlbumSharingGridPort GetAlbumSharingGridPort
	}
	type args struct {
		user   usermodel.CurrentUser
		filter ListAlbumsFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*VisibleAlbum
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should serve all display fields from the view",
			fields: fields{
				Repository:              ownerSummariesRepository(),
				GetAlbumSharingGridPort: emptySharingGrid,
			},
			args: args{user: ironmanCurrentUser, filter: noFilter},
			want: []*VisibleAlbum{
				{
					Album:      catalog.Album{AlbumId: sharedAlbum3Id, Name: "Album Three", Start: mar24, End: apr24},
					MediaCount: 3,
				},
				{
					Album:              catalog.Album{AlbumId: ownedAlbum2Id, Name: "Album Two", Start: feb24, End: mar24},
					MediaCount:         2,
					OwnedByCurrentUser: true,
				},
				{
					Album:              catalog.Album{AlbumId: ownedAlbum1Id, Name: "Album One", Start: jan24, End: feb24},
					MediaCount:         1,
					OwnedByCurrentUser: true,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should only return directly owned when filter is set",
			fields: fields{
				Repository:              ownerSummariesRepository(),
				GetAlbumSharingGridPort: emptySharingGrid,
			},
			args: args{user: ironmanCurrentUser, filter: ListAlbumsFilter{OnlyDirectlyOwned: true}},
			want: []*VisibleAlbum{
				{
					Album:              catalog.Album{AlbumId: ownedAlbum2Id, Name: "Album Two", Start: feb24, End: mar24},
					MediaCount:         2,
					OwnedByCurrentUser: true,
				},
				{
					Album:              catalog.Album{AlbumId: ownedAlbum1Id, Name: "Album One", Start: jan24, End: feb24},
					MediaCount:         1,
					OwnedByCurrentUser: true,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should decorate owned rows with the sharing grid",
			fields: fields{
				Repository:              ownerSummariesRepository(),
				GetAlbumSharingGridPort: ownerSharingGrid,
			},
			args: args{user: ironmanCurrentUser, filter: noFilter},
			want: []*VisibleAlbum{
				{
					Album:      catalog.Album{AlbumId: sharedAlbum3Id, Name: "Album Three", Start: mar24, End: apr24},
					MediaCount: 3,
				},
				{
					Album:              catalog.Album{AlbumId: ownedAlbum2Id, Name: "Album Two", Start: feb24, End: mar24},
					MediaCount:         2,
					OwnedByCurrentUser: true,
				},
				{
					Album:              catalog.Album{AlbumId: ownedAlbum1Id, Name: "Album One", Start: jan24, End: feb24},
					MediaCount:         1,
					Visitors:           []usermodel.UserId{visitorUserId},
					OwnedByCurrentUser: true,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return empty when the user has no rows",
			fields: fields{
				Repository:              &AlbumSummaryInMemoryRepository{},
				GetAlbumSharingGridPort: emptySharingGrid,
			},
			args:    args{user: visitorOnlyUser, filter: noFilter},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should sort by start desc then end desc",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{
							AlbumSummary: AlbumSummary{AlbumId: ownedAlbum1Id, Name: "Jan to Feb", Start: jan24, End: feb24},
							Availability: OwnerAvailability(ironmanCurrentUser.UserId),
						},
						{
							AlbumSummary: AlbumSummary{AlbumId: ownedAlbum2Id, Name: "Jan to Mar", Start: jan24, End: mar24},
							Availability: OwnerAvailability(ironmanCurrentUser.UserId),
						},
						{
							AlbumSummary: AlbumSummary{AlbumId: sharedAlbum3Id, Name: "Feb to Mar", Start: feb24, End: mar24},
							Availability: OwnerAvailability(ironmanCurrentUser.UserId),
						},
					},
				},
				GetAlbumSharingGridPort: emptySharingGrid,
			},
			args: args{user: ironmanCurrentUser, filter: noFilter},
			want: []*VisibleAlbum{
				{
					Album:              catalog.Album{AlbumId: sharedAlbum3Id, Name: "Feb to Mar", Start: feb24, End: mar24},
					OwnedByCurrentUser: true,
				},
				{
					Album:              catalog.Album{AlbumId: ownedAlbum2Id, Name: "Jan to Mar", Start: jan24, End: mar24},
					OwnedByCurrentUser: true,
				},
				{
					Album:              catalog.Album{AlbumId: ownedAlbum1Id, Name: "Jan to Feb", Start: jan24, End: feb24},
					OwnedByCurrentUser: true,
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			albumView := NewAlbumView(
				tt.fields.Repository,
				tt.fields.GetAlbumSharingGridPort,
				MediaCounterPortFake(nil),
				FindAlbumsByIdsFunc(func(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error) { return nil, nil }),
				stubListUserWhoCanAccessAlbumPort(nil),
			)

			got, err := albumView.ListAlbums(context.Background(), tt.args.user, tt.args.filter)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAlbumView_AlbumCreated(t *testing.T) {
	tonyOwner := ownermodel.Owner("tony")
	tonyUserId := usermodel.UserId("tony@stark.com")
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	mar24 := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)

	newAlbumId := catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("album-new")}
	sourceAlbumId := catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("album-source")}
	newAlbum := catalog.Album{AlbumId: newAlbumId, Name: "New Album", Start: feb24, End: mar24}

	type fields struct {
		Repository                     *AlbumSummaryInMemoryRepository
		MediaCounterPort               MediaCounterPort
		ListUsersWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
	}
	tests := []struct {
		name            string
		fields          fields
		event           catalog.AlbumCreated
		expectSummaries []UserAlbumSummary
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "it should make the album visible to the owner",
			fields: fields{
				Repository:       &AlbumSummaryInMemoryRepository{},
				MediaCounterPort: MediaCounterPortFake(nil),
				ListUsersWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(map[catalog.AlbumId][]Availability{
					newAlbumId: {OwnerAvailability(tonyUserId)},
				}),
			},
			event: catalog.AlbumCreated{
				CreatedAlbum:      newAlbum,
				TransferredMedias: catalog.NewTransferredMedias(),
			},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: newAlbumId, Name: "New Album", Start: feb24, End: mar24, MediaCount: 0},
					Availability: OwnerAvailability(tonyUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not shadow an existing row (PutSummaries overwrites Count with the transfer-derived value)",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{
							AlbumSummary: AlbumSummary{AlbumId: newAlbumId, Name: "Stale Name", Start: jan24, End: feb24, MediaCount: 5},
							Availability: OwnerAvailability(tonyUserId),
						},
					},
				},
				MediaCounterPort: MediaCounterPortFake(nil),
				ListUsersWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(map[catalog.AlbumId][]Availability{
					newAlbumId: {OwnerAvailability(tonyUserId)},
				}),
			},
			event: catalog.AlbumCreated{
				CreatedAlbum:      newAlbum,
				TransferredMedias: catalog.NewTransferredMedias(),
			},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: newAlbumId, Name: "New Album", Start: feb24, End: mar24, MediaCount: 0},
					Availability: OwnerAvailability(tonyUserId),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should recount transferred source albums and leave their display fields untouched",
			fields: fields{
				Repository: &AlbumSummaryInMemoryRepository{
					Summaries: []UserAlbumSummary{
						{
							AlbumSummary: AlbumSummary{AlbumId: sourceAlbumId, Name: "Source Album", Start: jan24, End: feb24, MediaCount: 10},
							Availability: OwnerAvailability(tonyUserId),
						},
					},
				},
				MediaCounterPort: MediaCounterPortFake(map[catalog.AlbumId]int{
					sourceAlbumId: 7,
				}),
				ListUsersWhoCanAccessAlbumPort: stubListUserWhoCanAccessAlbumPort(map[catalog.AlbumId][]Availability{
					newAlbumId:    {OwnerAvailability(tonyUserId)},
					sourceAlbumId: {OwnerAvailability(tonyUserId)},
				}),
			},
			event: catalog.AlbumCreated{
				CreatedAlbum: newAlbum,
				TransferredMedias: catalog.TransferredMedias{
					Transfers: map[catalog.AlbumId][]catalog.MediaId{
						newAlbumId: {"media-1", "media-2", "media-3"},
					},
					FromAlbums: []catalog.AlbumId{sourceAlbumId},
				},
			},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: sourceAlbumId, Name: "Source Album", Start: jan24, End: feb24, MediaCount: 7},
					Availability: OwnerAvailability(tonyUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: newAlbumId, Name: "New Album", Start: feb24, End: mar24, MediaCount: 3},
					Availability: OwnerAvailability(tonyUserId),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			albumView := NewAlbumView(
				tt.fields.Repository,
				GetAlbumSharingGridFunc(func(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) {
					return nil, nil
				}),
				tt.fields.MediaCounterPort,
				FindAlbumsByIdsFunc(func(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error) { return nil, nil }),
				tt.fields.ListUsersWhoCanAccessAlbumPort,
			)

			err := albumView.AlbumCreated(context.Background(), tt.event)
			if tt.wantErr(t, err) {
				assert.ElementsMatch(t, tt.expectSummaries, tt.fields.Repository.Summaries)
			}
		})
	}
}
