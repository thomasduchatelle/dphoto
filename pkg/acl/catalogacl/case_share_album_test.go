package catalogacl_test

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"testing"
	"time"
)

func TestShareAlbumCase_ShareAlbumWith(t *testing.T) {
	theDate := time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC)
	aclcore.TimeFunc = func() time.Time {
		return theDate
	}

	type args struct {
		owner      ownermodel.Owner
		folderName catalog.FolderName
		userEmail  usermodel.UserId
	}
	const owner = ownermodel.Owner("tony@stark.com")
	folderName := catalog.NewFolderName("/weddings")
	albumId := catalog.AlbumId{Owner: owner, FolderName: folderName}
	const userEmail = usermodel.UserId("pepper@stark.com")

	tests := []struct {
		name         string
		fields       func(t *testing.T) (aclcore.ScopeWriter, catalogacl.FindAlbumPort)
		args         args
		wantObserved map[catalog.AlbumId][]usermodel.UserId
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "it should create the ACL rule when the album exists",
			fields: func(t *testing.T) (aclcore.ScopeWriter, catalogacl.FindAlbumPort) {
				catalogMock := mocks.NewFindAlbumPort(t)
				catalogMock.EXPECT().FindAlbum(mock.Anything, albumId).Return(&catalog.Album{
					AlbumId: albumId,
				}, nil)

				scopeWriter := mocks.NewScopeWriter(t)
				scopeWriter.On("SaveIfNewScope", aclcore.Scope{
					Type:          aclcore.AlbumVisitorScope,
					GrantedAt:     theDate,
					GrantedTo:     userEmail,
					ResourceOwner: owner,
					ResourceId:    folderName.String(),
				}).Return(nil)

				return scopeWriter, catalogMock
			},
			args: args{owner, folderName, userEmail},
			wantObserved: map[catalog.AlbumId][]usermodel.UserId{
				albumId: {userEmail},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return an error if the album doesn't exists",
			fields: func(t *testing.T) (aclcore.ScopeWriter, catalogacl.FindAlbumPort) {
				catalogMock := mocks.NewFindAlbumPort(t)
				catalogMock.EXPECT().FindAlbum(mock.Anything, albumId).Return(nil, catalog.AlbumNotFoundErr)

				return mocks.NewScopeWriter(t), catalogMock
			},
			args:         args{owner, folderName, userEmail},
			wantObserved: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i)
			},
		},
		{
			name: "it should passthroughs an other error",
			fields: func(t *testing.T) (aclcore.ScopeWriter, catalogacl.FindAlbumPort) {
				catalogMock := mocks.NewFindAlbumPort(t)
				catalogMock.EXPECT().FindAlbum(mock.Anything, albumId).Return(nil, errors.New("TEST Something else"))

				return mocks.NewScopeWriter(t), catalogMock
			},
			args:         args{owner, folderName, userEmail},
			wantObserved: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && assert.Contains(t, err.Error(), "TEST Something else", i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := new(AlbumSharedObserverFake)

			scopeWriter, catalogPort := tt.fields(t)
			s := &catalogacl.ShareAlbumCase{
				ScopeWriter:   scopeWriter,
				FindAlbumPort: catalogPort,
				Observers:     []catalogacl.AlbumSharedObserver{observer},
			}

			err := s.ShareAlbumWith(context.TODO(), catalog.AlbumId{Owner: tt.args.owner, FolderName: tt.args.folderName}, tt.args.userEmail)
			if !tt.wantErr(t, err, fmt.Sprintf("ShareAlbumWith(%v, %v, %v)", tt.args.owner, tt.args.folderName, tt.args.userEmail)) {
				return
			}

			assert.Equalf(t, tt.wantObserved, observer.Shared, "Shared=%+v", observer.Shared)
		})
	}
}

type AlbumSharedObserverFake struct {
	Shared map[catalog.AlbumId][]usermodel.UserId
}

func (a *AlbumSharedObserverFake) AlbumShared(ctx context.Context, albumId catalog.AlbumId, userEmail usermodel.UserId) error {
	if a.Shared == nil {
		a.Shared = make(map[catalog.AlbumId][]usermodel.UserId)
	}

	previous, _ := a.Shared[albumId]
	a.Shared[albumId] = append(previous, userEmail)
	return nil
}
