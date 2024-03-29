package catalogacl_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mocks2 "github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"testing"
	"time"
)

func TestShareAlbumCase_ShareAlbumWith(t *testing.T) {
	theDate := time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC)
	aclcore.TimeFunc = func() time.Time {
		return theDate
	}

	type args struct {
		owner      string
		folderName string
		userEmail  string
	}
	const owner = "tony@stark.com"
	const folderName = "/weddings"
	const userEmail = "pepper@stark.com"

	tests := []struct {
		name    string
		fields  func(t *testing.T) (aclcore.ScopeWriter, catalogacl.ShareAlbumCatalogPort)
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should create the ACL rule when the album exists",
			fields: func(t *testing.T) (aclcore.ScopeWriter, catalogacl.ShareAlbumCatalogPort) {
				catalogMock := mocks2.NewShareAlbumCatalogPort(t)
				catalogMock.On("FindAlbum", owner, folderName).Return(&catalog.Album{
					Owner:      owner,
					FolderName: folderName,
				}, nil)

				scopeWriter := mocks2.NewScopeWriter(t)
				scopeWriter.On("SaveIfNewScope", aclcore.Scope{
					Type:          aclcore.AlbumVisitorScope,
					GrantedAt:     theDate,
					GrantedTo:     userEmail,
					ResourceOwner: owner,
					ResourceId:    folderName,
				}).Return(nil)

				return scopeWriter, catalogMock
			},
			args:    args{owner, folderName, userEmail},
			wantErr: assert.NoError,
		},
		{
			name: "it should return an error if the album doesn't exists",
			fields: func(t *testing.T) (aclcore.ScopeWriter, catalogacl.ShareAlbumCatalogPort) {
				catalogMock := mocks2.NewShareAlbumCatalogPort(t)
				catalogMock.On("FindAlbum", owner, folderName).Return(nil, catalog.NotFoundError)

				return mocks2.NewScopeWriter(t), catalogMock
			},
			args: args{owner, folderName, userEmail},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.NotFoundError, i)
			},
		},
		{
			name: "it should passthroughs an other error",
			fields: func(t *testing.T) (aclcore.ScopeWriter, catalogacl.ShareAlbumCatalogPort) {
				catalogMock := mocks2.NewShareAlbumCatalogPort(t)
				catalogMock.On("FindAlbum", owner, folderName).Return(nil, errors.New("TEST Something else"))

				return mocks2.NewScopeWriter(t), catalogMock
			},
			args: args{owner, folderName, userEmail},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && assert.Contains(t, err.Error(), "TEST Something else", i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scopeWriter, catalogPort := tt.fields(t)
			s := &catalogacl.ShareAlbumCase{
				ScopeWriter: scopeWriter,
				CatalogPort: catalogPort,
			}
			tt.wantErr(t, s.ShareAlbumWith(tt.args.owner, tt.args.folderName, tt.args.userEmail, aclcore.AlbumVisitorScope), fmt.Sprintf("ShareAlbumWith(%v, %v, %v)", tt.args.owner, tt.args.folderName, tt.args.userEmail))
		})
	}
}
