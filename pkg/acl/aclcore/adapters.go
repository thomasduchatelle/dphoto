package aclcore

import (
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type ScopesReader interface {
	// ListUserScopes returns all access of a certain type that have been granted to a user
	ListUserScopes(email usermodel.UserId, types ...ScopeType) ([]*Scope, error)

	// FindScopesById returns scopes that have been granted (exists in DB)
	FindScopesById(ids ...ScopeId) ([]*Scope, error)
}

type ReverseScopesReader interface {
	// ListOwnerScopes is a reverse query to find to whom has been shared owner resources.
	ListOwnerScopes(owner ownermodel.Owner, types ...ScopeType) ([]*Scope, error)
}

type ScopeWriter interface {
	// DeleteScopes deletes the scope(s) if it exists, do nothing otherwise
	DeleteScopes(id ...ScopeId) error

	// SaveIfNewScope persists the scope if it doesn't exist yet, no error is returned if it already exists
	SaveIfNewScope(scope Scope) error
}

type RefreshTokenRepository interface {
	StoreRefreshToken(token string, spec RefreshTokenSpec) error

	FindRefreshToken(token string) (*RefreshTokenSpec, error)
	DeleteRefreshToken(token string) error

	// HouseKeepRefreshToken removes any token that have expired
	HouseKeepRefreshToken() (int, error)
}

type IdentityDetailsStore interface {
	StoreIdentity(identity Identity) error
	FindIdentity(email usermodel.UserId) (*Identity, error)
}
