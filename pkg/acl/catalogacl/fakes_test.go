package catalogacl_test

import (
	"context"
	"slices"

	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type FindAlbumPortInMemory struct {
	Albums map[catalog.AlbumId]*catalog.Album
}

func (f *FindAlbumPortInMemory) FindAlbum(ctx context.Context, albumId catalog.AlbumId) (*catalog.Album, error) {
	if album, ok := f.Albums[albumId]; ok {
		return album, nil
	}
	return nil, catalog.AlbumNotFoundErr
}

type ScopeRepositoryInMemory struct {
	Scopes map[usermodel.UserId][]*aclcore.Scope
}

func (s *ScopeRepositoryInMemory) SaveIfNewScope(scope aclcore.Scope) error {
	if s.Scopes == nil {
		s.Scopes = make(map[usermodel.UserId][]*aclcore.Scope)
	}
	for _, existing := range s.Scopes[scope.GrantedTo] {
		if existing.Id() == scope.Id() {
			return nil
		}
	}
	saved := scope
	s.Scopes[scope.GrantedTo] = append(s.Scopes[scope.GrantedTo], &saved)
	return nil
}

func (s *ScopeRepositoryInMemory) DeleteScopes(ids ...aclcore.ScopeId) error {
	for _, id := range ids {
		remaining := s.Scopes[id.GrantedTo][:0]
		for _, scope := range s.Scopes[id.GrantedTo] {
			if scope.Id() != id {
				remaining = append(remaining, scope)
			}
		}
		s.Scopes[id.GrantedTo] = remaining
	}
	return nil
}

func (s *ScopeRepositoryInMemory) FindScopesById(ids ...aclcore.ScopeId) ([]*aclcore.Scope, error) {
	var found []*aclcore.Scope
	for _, scopes := range s.Scopes {
		for _, scope := range scopes {
			if slices.Contains(ids, scope.Id()) {
				found = append(found, scope)
			}
		}
	}
	return found, nil
}

func (s *ScopeRepositoryInMemory) ListScopesByUser(ctx context.Context, email usermodel.UserId, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	var found []*aclcore.Scope
	for _, scope := range s.Scopes[email] {
		if slices.Contains(types, scope.Type) {
			found = append(found, scope)
		}
	}
	return found, nil
}
