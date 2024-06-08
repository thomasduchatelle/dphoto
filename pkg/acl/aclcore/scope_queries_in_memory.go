package aclcore

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"slices"
)

// ScopeReadRepositoryInMemory implements ScopeReadRepository and is used for testing and stubbing.
type ScopeReadRepositoryInMemory struct {
	Scopes []*Scope
}

func (s *ScopeReadRepositoryInMemory) FindScopesByIdCtx(ctx context.Context, ids ...ScopeId) ([]*Scope, error) {
	var scopes []*Scope
	for _, scope := range s.Scopes {
		if slices.Contains(ids, scope.Id()) {
			scopes = append(scopes, scope)
		}
	}

	return scopes, nil
}

func (s *ScopeReadRepositoryInMemory) ListScopesByOwner(ctx context.Context, owner ownermodel.Owner, types ...ScopeType) ([]*Scope, error) {
	var scopes []*Scope
	for _, scope := range s.Scopes {
		if scope.ResourceOwner == owner && slices.Contains(types, scope.Type) {
			scopes = append(scopes, scope)
		}
	}

	return scopes, nil
}

func (s *ScopeReadRepositoryInMemory) ListScopesByUser(ctx context.Context, userId usermodel.UserId, types ...ScopeType) ([]*Scope, error) {
	var scopes []*Scope
	for _, scope := range s.Scopes {
		if scope.GrantedTo == userId && slices.Contains(types, scope.Type) {
			scopes = append(scopes, scope)
		}
	}

	return scopes, nil
}
