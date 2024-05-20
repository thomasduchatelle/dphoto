package catalogviews

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
	"time"
)

func TestNewAlbumViewAcceptance(t *testing.T) {
	tonyOwner := ownermodel.Owner("tony")
	pepperOwner := ownermodel.Owner("pepper")
	userId2 := usermodel.UserId("user-id-02")
	ironmanCurrentUser := usermodel.CurrentUser{
		UserId: "ironman@avenger.hero",
		Owner:  &tonyOwner,
	}
	noFilter := ListAlbumsFilter{}
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	mar24 := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)
	ownedAlbum1 := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: *ironmanCurrentUser.Owner, FolderName: catalog.NewFolderName("album-1")},
		Start:   jan24,
		End:     feb24,
	}
	ownedAlbum2 := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: *ironmanCurrentUser.Owner, FolderName: catalog.NewFolderName("album-2")},
		Start:   feb24,
		End:     mar24,
	}
	sharedAlbum3 := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: pepperOwner, FolderName: catalog.NewFolderName("album-3")},
		Start:   mar24,
		End:     apr24,
	}

	type fields struct {
		FindAlbumByOwnerPort        FindAlbumByOwnerPort
		GetAlbumSharingGridPort     GetAlbumSharingGridPort
		FindAlbumsByIdsPort         FindAlbumsByIdsPort
		SharedWithUserPort          SharedWithUserPort
		GetAvailabilitiesByUserPort GetAvailabilitiesByUserPort
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
			name: "it should aggregate the list of owned albums with the list of shared albums",
			fields: fields{
				FindAlbumByOwnerPort:    stubFindAlbumByOwnerPort(&ownedAlbum1, &ownedAlbum2),
				GetAlbumSharingGridPort: stubGetAlbumSharingGridPort(ownedAlbum1.AlbumId, userId2),
				FindAlbumsByIdsPort:     stubFindAlbumsByIdsPort(sharedAlbum3),
				SharedWithUserPort:      stubSharedWithUserPort(sharedAlbum3),
				GetAvailabilitiesByUserPort: stubGetAvailabilitiesByUserPortFromCount(map[catalog.AlbumId]int{
					ownedAlbum1.AlbumId:  1,
					ownedAlbum2.AlbumId:  2,
					sharedAlbum3.AlbumId: 3,
				}),
			},
			args: args{
				user:   ironmanCurrentUser,
				filter: noFilter,
			},
			want: []*VisibleAlbum{
				{
					Album:      sharedAlbum3,
					MediaCount: 3,
				},
				{
					Album:              ownedAlbum2,
					MediaCount:         2,
					OwnedByCurrentUser: true,
				},
				{
					Album:              ownedAlbum1,
					MediaCount:         1,
					Visitors:           []usermodel.UserId{userId2},
					OwnedByCurrentUser: true,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should only get the list of owned albums",
			fields: fields{
				FindAlbumByOwnerPort:    stubFindAlbumByOwnerPort(&ownedAlbum1, &ownedAlbum2),
				GetAlbumSharingGridPort: stubGetAlbumSharingGridPort(ownedAlbum1.AlbumId, userId2),
				FindAlbumsByIdsPort:     stubFindAlbumsByIdsPort(sharedAlbum3),
				SharedWithUserPort:      stubSharedWithUserPort(sharedAlbum3),
				GetAvailabilitiesByUserPort: stubGetAvailabilitiesByUserPortFromCount(map[catalog.AlbumId]int{
					ownedAlbum1.AlbumId:  1,
					ownedAlbum2.AlbumId:  2,
					sharedAlbum3.AlbumId: 3,
				}),
			},
			args: args{
				user:   ironmanCurrentUser,
				filter: ListAlbumsFilter{OnlyDirectlyOwned: true},
			},
			want: []*VisibleAlbum{
				{
					Album:              ownedAlbum2,
					MediaCount:         2,
					OwnedByCurrentUser: true,
				},
				{
					Album:              ownedAlbum1,
					MediaCount:         1,
					Visitors:           []usermodel.UserId{userId2},
					OwnedByCurrentUser: true,
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			albumView := NewAlbumView(
				tt.fields.FindAlbumByOwnerPort,
				tt.fields.GetAlbumSharingGridPort,
				tt.fields.FindAlbumsByIdsPort,
				tt.fields.SharedWithUserPort,
				tt.fields.GetAvailabilitiesByUserPort,
			)

			got, err := albumView.ListAlbums(context.Background(), tt.args.user, tt.args.filter)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAlbumView_ListAlbums(t *testing.T) {
	owner := ownermodel.Owner("tony")
	ironmanCurrentUser := usermodel.CurrentUser{
		UserId: "ironman@avenger.hero",
		Owner:  &owner,
	}
	noFilter := ListAlbumsFilter{}
	jan24 := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb24 := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	mar24 := time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
	janToFeb24 := &VisibleAlbum{
		Album: catalog.Album{
			AlbumId: catalog.AlbumId{Owner: *ironmanCurrentUser.Owner, FolderName: catalog.NewFolderName("jan24-feb24")},
			Start:   jan24,
			End:     feb24,
		},
		MediaCount: 42,
	}
	febToMar24 := &VisibleAlbum{
		Album: catalog.Album{
			AlbumId: catalog.AlbumId{Owner: *ironmanCurrentUser.Owner, FolderName: catalog.NewFolderName("feb24-mar24")},
			Start:   feb24,
			End:     mar24,
		},
		MediaCount: 42,
	}
	janToMar24 := &VisibleAlbum{
		Album: catalog.Album{
			AlbumId: catalog.AlbumId{Owner: *ironmanCurrentUser.Owner, FolderName: catalog.NewFolderName("jan24-mar24")},
			Start:   jan24,
			End:     mar24,
		},
		MediaCount: 42,
	}
	testError := errors.New("TEST error on the provider")

	type fields struct {
		Providers []ListAlbumsProvider
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
			name: "it should not fail if there is no provider",
			fields: fields{
				Providers: nil,
			},
			args: args{
				user:   ironmanCurrentUser,
				filter: noFilter,
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "it should return the albums from the different providers",
			fields: fields{
				Providers: []ListAlbumsProvider{
					stubListAlbumsProvider(ironmanCurrentUser, noFilter, febToMar24),
					stubListAlbumsProvider(ironmanCurrentUser, noFilter, janToFeb24),
				},
			},
			args: args{
				user:   ironmanCurrentUser,
				filter: noFilter,
			},
			want:    []*VisibleAlbum{febToMar24, janToFeb24},
			wantErr: assert.NoError,
		},
		{
			name: "it should order the albums by reverse chronological order on start date",
			fields: fields{
				Providers: []ListAlbumsProvider{
					stubListAlbumsProvider(ironmanCurrentUser, noFilter, janToFeb24),
					stubListAlbumsProvider(ironmanCurrentUser, noFilter, febToMar24),
				},
			},
			args: args{
				user:   ironmanCurrentUser,
				filter: noFilter,
			},
			want:    []*VisibleAlbum{febToMar24, janToFeb24},
			wantErr: assert.NoError,
		},
		{
			name: "it should order the albums by reverse chronological order on end date (if start date the same)",
			fields: fields{
				Providers: []ListAlbumsProvider{
					stubListAlbumsProvider(ironmanCurrentUser, noFilter, janToMar24, janToFeb24),
				},
			},
			args: args{
				user:   ironmanCurrentUser,
				filter: noFilter,
			},
			want:    []*VisibleAlbum{janToMar24, janToFeb24},
			wantErr: assert.NoError,
		},
		{
			name: "it should order the albums by reverse chronological order on end date (if start date the same)",
			fields: fields{
				Providers: []ListAlbumsProvider{
					stubListAlbumsProviderWithError(testError),
				},
			},
			args: args{
				user:   ironmanCurrentUser,
				filter: noFilter,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, testError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &AlbumView{
				Providers: tt.fields.Providers,
			}

			got, err := v.ListAlbums(context.Background(), tt.args.user, tt.args.filter)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func stubSharedWithUserPort(sharedAlbum3 catalog.Album) SharedWithUserFunc {
	return func(ctx context.Context, userId usermodel.UserId) ([]catalog.AlbumId, error) {
		return []catalog.AlbumId{sharedAlbum3.AlbumId}, nil
	}
}

func stubFindAlbumsByIdsPort(sharedAlbum3 catalog.Album) FindAlbumsByIdsFunc {
	return func(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error) {
		if len(ids) > 0 && ids[0] == sharedAlbum3.AlbumId {
			return []*catalog.Album{&sharedAlbum3}, nil
		}
		return nil, nil
	}
}

func stubListAlbumsProviderWithError(testError error) ListAlbumsProviderFunc {
	return func(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error) {
		return nil, testError
	}
}

func stubListAlbumsProvider(expectedUser usermodel.CurrentUser, expectedFilter ListAlbumsFilter, albums ...*VisibleAlbum) ListAlbumsProvider {
	return ListAlbumsProviderFunc(func(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error) {
		if expectedUser.UserId != user.UserId {
			return nil, nil
		}
		if expectedFilter.OnlyDirectlyOwned != filter.OnlyDirectlyOwned {
			return nil, nil
		}

		return albums, nil
	})
}

func stubFindAlbumByOwnerPort(albums ...*catalog.Album) FindAlbumByOwnerFunc {
	return func(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
		return albums, nil
	}
}

func stubGetAlbumSharingGridPort(albumId catalog.AlbumId, userId2 usermodel.UserId) GetAlbumSharingGridFunc {
	return func(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) {
		if owner == albumId.Owner {
			return map[catalog.AlbumId][]usermodel.UserId{
				albumId: {userId2},
			}, nil
		}
		return nil, nil
	}
}

func stubMediaCounterPort(m map[catalog.AlbumId]int) MediaCounterPort {
	return MediaCounterFunc(func(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
		return m, nil
	})
}

func stubGetAvailabilitiesByUserPortFromCount(m map[catalog.AlbumId]int) GetAvailabilitiesByUserPort {
	return GetAvailabilitiesByUserFunc(func(ctx context.Context, userId usermodel.UserId) ([]AlbumSize, error) {
		var availabilities []AlbumSize
		for albumId, count := range m {
			availabilities = append(availabilities, AlbumSize{
				AlbumId:    albumId,
				MediaCount: count,
			})
		}
		return availabilities, nil
	})
}
