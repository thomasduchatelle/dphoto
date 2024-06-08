package aclcore

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type ScopeReadRepository interface {
	ListScopesByUser(ctx context.Context, email usermodel.UserId, types ...ScopeType) ([]*Scope, error)
	ListScopesByOwner(ctx context.Context, owner ownermodel.Owner, types ...ScopeType) ([]*Scope, error)
	FindScopesByIdCtx(ctx context.Context, ids ...ScopeId) ([]*Scope, error)
}

// ScopeQueries is a read-only repository to query scopes
type ScopeQueries struct {
	ScopeReadRepository
}
