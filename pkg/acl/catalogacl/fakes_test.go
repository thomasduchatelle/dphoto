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

func NewFindAlbumPortInMemory(albums ...*catalog.Album) *FindAlbumPortInMemory {
	repo := &FindAlbumPortInMemory{Albums: make(map[catalog.AlbumId]*catalog.Album)}
	for _, album := range albums {
		repo.Albums[album.AlbumId] = album
	}
	return repo
}

func (f *FindAlbumPortInMemory) FindAlbum(ctx context.Context, albumId catalog.AlbumId) (*catalog.Album, error) {
	album, ok := f.Albums[albumId]
	if !ok {
		return nil, catalog.AlbumNotFoundErr
	}
	return album, nil
}

type ScopeRepositoryInMemory struct {
	Scopes map[usermodel.UserId][]*aclcore.Scope
}

func NewScopeRepositoryInMemory(scopes ...aclcore.Scope) *ScopeRepositoryInMemory {
	repo := &ScopeRepositoryInMemory{Scopes: make(map[usermodel.UserId][]*aclcore.Scope)}
	for _, scope := range scopes {
		s := scope
		repo.Scopes[s.GrantedTo] = append(repo.Scopes[s.GrantedTo], &s)
	}
	return repo
}

func (s *ScopeRepositoryInMemory) SaveIfNewScope(scope aclcore.Scope) error {
	for _, existing := range s.Scopes[scope.GrantedTo] {
		if existing.Id() == scope.Id() {
			return nil
		}
	}
	s.Scopes[scope.GrantedTo] = append(s.Scopes[scope.GrantedTo], &scope)
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
	var scopes []*aclcore.Scope
	for _, userScopes := range s.Scopes {
		for _, scope := range userScopes {
			if slices.Contains(ids, scope.Id()) {
				scopes = append(scopes, scope)
			}
		}
	}
	return scopes, nil
}

func (s *ScopeRepositoryInMemory) ListScopesByUser(ctx context.Context, userId usermodel.UserId, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	var scopes []*aclcore.Scope
	for _, scope := range s.Scopes[userId] {
		if slices.Contains(types, scope.Type) {
			scopes = append(scopes, scope)
		}
	}
	return scopes, nil
}
