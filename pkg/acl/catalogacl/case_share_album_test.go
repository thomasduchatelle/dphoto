package catalogacl_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func TestShareAlbumCase_ShareAlbumWith(t *testing.T) {
	theDate := time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC)
	aclcore.TimeFunc = func() time.Time {
		return theDate
	}

	const owner = ownermodel.Owner("tony@stark.com")
	folderName := catalog.NewFolderName("/weddings")
	albumId := catalog.AlbumId{Owner: owner, FolderName: folderName}
	const userEmail = usermodel.UserId("pepper@stark.com")
	expectedScopeId := aclcore.ScopeId{
		Type:          aclcore.AlbumVisitorScope,
		GrantedTo:     userEmail,
		ResourceOwner: owner,
		ResourceId:    folderName.String(),
	}
	expectedScope := &aclcore.Scope{
		Type:          aclcore.AlbumVisitorScope,
		GrantedAt:     theDate,
		GrantedTo:     userEmail,
		ResourceOwner: owner,
		ResourceId:    folderName.String(),
	}

	albumExists := &FindAlbumPortInMemory{Albums: map[catalog.AlbumId]*catalog.Album{
		albumId: {AlbumId: albumId},
	}}

	type fields struct {
		ScopeWriter   *ScopeRepositoryInMemory
		FindAlbumPort catalogacl.FindAlbumPort
		Observers     []catalogacl.AlbumSharedObserver
	}
	type args struct {
		albumId   catalog.AlbumId
		userEmail usermodel.UserId
	}

	newFields := func(findAlbum catalogacl.FindAlbumPort) fields {
		return fields{
			ScopeWriter:   &ScopeRepositoryInMemory{},
			FindAlbumPort: findAlbum,
			Observers:     []catalogacl.AlbumSharedObserver{new(AlbumSharedObserverFake)},
		}
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		wantObserved map[catalog.AlbumId][]usermodel.UserId
		wantScopes   []*aclcore.Scope
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:         "it should create the ACL rule when the album exists",
			fields:       newFields(albumExists),
			args:         args{albumId: albumId, userEmail: userEmail},
			wantObserved: map[catalog.AlbumId][]usermodel.UserId{albumId: {userEmail}},
			wantScopes:   []*aclcore.Scope{expectedScope},
			wantErr:      assert.NoError,
		},
		{
			name:         "it should return an error if the album doesn't exists",
			fields:       newFields(&FindAlbumPortInMemory{}),
			args:         args{albumId: albumId, userEmail: userEmail},
			wantObserved: nil,
			wantScopes:   nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &catalogacl.ShareAlbumCase{
				ScopeWriter:   tt.fields.ScopeWriter,
				FindAlbumPort: tt.fields.FindAlbumPort,
				Observers:     tt.fields.Observers,
			}

			err := s.ShareAlbumWith(context.TODO(), tt.args.albumId, tt.args.userEmail)
			if !tt.wantErr(t, err, fmt.Sprintf("ShareAlbumWith(%v, %v)", tt.args.albumId, tt.args.userEmail)) {
				return
			}

			observer := tt.fields.Observers[0].(*AlbumSharedObserverFake)
			assert.Equalf(t, tt.wantObserved, observer.Shared, "Shared=%+v", observer.Shared)

			savedScopes, findErr := tt.fields.ScopeWriter.FindScopesById(expectedScopeId)
			assert.NoError(t, findErr)
			assert.Equal(t, tt.wantScopes, savedScopes)
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
