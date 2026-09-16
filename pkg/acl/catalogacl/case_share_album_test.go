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
	weddingsAlbum := catalog.Album{
		AlbumId: albumId,
		Name:    "Weddings",
		Start:   time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
		End:     time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC),
	}

	expectedScopeId := aclcore.ScopeId{
		Type:          aclcore.AlbumVisitorScope,
		GrantedTo:     userEmail,
		ResourceOwner: owner,
		ResourceId:    folderName.String(),
	}
	expectedScope := aclcore.Scope{
		Type:          aclcore.AlbumVisitorScope,
		GrantedAt:     theDate,
		GrantedTo:     userEmail,
		ResourceOwner: owner,
		ResourceId:    folderName.String(),
	}

	findAlbumPortWithWeddings := func() *FindAlbumPortInMemory {
		return NewFindAlbumPortInMemory(&weddingsAlbum)
	}

	type fields struct {
		ScopeRepository *ScopeRepositoryInMemory
		FindAlbumPort   *FindAlbumPortInMemory
		Observer        *AlbumSharedObserverFake
	}
	type args struct {
		owner      ownermodel.Owner
		folderName catalog.FolderName
		userEmail  usermodel.UserId
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectScopes []*aclcore.Scope
		expectShared map[usermodel.UserId][]catalog.Album
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "it should create the ACL rule and notify observers with the full album when the album exists",
			fields: fields{
				ScopeRepository: NewScopeRepositoryInMemory(),
				FindAlbumPort:   findAlbumPortWithWeddings(),
				Observer:        &AlbumSharedObserverFake{},
			},
			args:         args{owner, folderName, userEmail},
			expectScopes: []*aclcore.Scope{&expectedScope},
			expectShared: map[usermodel.UserId][]catalog.Album{
				userEmail: {weddingsAlbum},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return an error if the album doesn't exists",
			fields: fields{
				ScopeRepository: NewScopeRepositoryInMemory(),
				FindAlbumPort:   NewFindAlbumPortInMemory(),
				Observer:        &AlbumSharedObserverFake{},
			},
			args:         args{owner, folderName, userEmail},
			expectScopes: nil,
			expectShared: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &catalogacl.ShareAlbumCase{
				ScopeWriter:   tt.fields.ScopeRepository,
				FindAlbumPort: tt.fields.FindAlbumPort,
				Observers:     []catalogacl.AlbumSharedObserver{tt.fields.Observer},
			}

			err := s.ShareAlbumWith(context.TODO(), catalog.AlbumId{Owner: tt.args.owner, FolderName: tt.args.folderName}, tt.args.userEmail)
			if !tt.wantErr(t, err, fmt.Sprintf("ShareAlbumWith(%v, %v, %v)", tt.args.owner, tt.args.folderName, tt.args.userEmail)) {
				return
			}

			storedScopes, findErr := tt.fields.ScopeRepository.FindScopesById(expectedScopeId)
			assert.NoError(t, findErr)
			assert.Equal(t, tt.expectScopes, storedScopes, "stored scopes")
			assert.Equal(t, tt.expectShared, tt.fields.Observer.Shared, "Shared=%+v", tt.fields.Observer.Shared)
		})
	}
}

type AlbumSharedObserverFake struct {
	Shared map[usermodel.UserId][]catalog.Album
}

func (a *AlbumSharedObserverFake) AlbumShared(ctx context.Context, album catalog.Album, userEmail usermodel.UserId) error {
	if a.Shared == nil {
		a.Shared = make(map[usermodel.UserId][]catalog.Album)
	}

	a.Shared[userEmail] = append(a.Shared[userEmail], album)
	return nil
}
