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
			)

			got, err := albumView.ListAlbums(context.Background(), tt.args.user, tt.args.filter)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAlbumView_AlbumShared(t *testing.T) {
	tonyOwner := ownermodel.Owner("tony")
	visitorUserId := usermodel.UserId("pepper@stark.com")
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	albumId := catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("album-1")}
	album := catalog.Album{AlbumId: albumId, Name: "Album One", Start: jan24, End: feb24}

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
				Repository:       &AlbumSummaryInMemoryRepository{},
				MediaCounterPort: MediaCounterPortFake{albumId: 7},
			},
			args: args{album: album, userId: visitorUserId},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumId, Name: "Album One", Start: jan24, End: feb24, MediaCount: 7},
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
	ownerUserId := usermodel.UserId("tony@stark.com")
	visitorUserId := usermodel.UserId("pepper@stark.com")
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	albumId := catalog.AlbumId{Owner: tonyOwner, FolderName: catalog.NewFolderName("album-1")}

	ownerAndVisitorRepository := func() *AlbumSummaryInMemoryRepository {
		return &AlbumSummaryInMemoryRepository{
			Summaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumId, Name: "Album One", Start: jan24, End: feb24, MediaCount: 3},
					Availability: OwnerAvailability(ownerUserId),
				},
				{
					AlbumSummary: AlbumSummary{AlbumId: albumId, Name: "Album One", Start: jan24, End: feb24, MediaCount: 3},
					Availability: VisitorAvailability(visitorUserId),
				},
			},
		}
	}

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
			name:   "it should remove the visitor row only",
			fields: fields{Repository: ownerAndVisitorRepository()},
			args:   args{albumId: albumId, userId: visitorUserId},
			expectSummaries: []UserAlbumSummary{
				{
					AlbumSummary: AlbumSummary{AlbumId: albumId, Name: "Album One", Start: jan24, End: feb24, MediaCount: 3},
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
