package aclcore

type ScopesReader interface {
	// ListUserScopes returns all access of a certain type that have been granted to a user
	ListUserScopes(email string, types ...ScopeType) ([]*Scope, error)

	// FindScopesById returns scopes that have been granted (exists in DB)
	FindScopesById(ids ...ScopeId) ([]*Scope, error)
}

type ReverseScopesReader interface {
	// ListOwnerScopes is a reverse query to find to whom has been shared owner resources.
	ListOwnerScopes(owner string, types ...ScopeType) ([]*Scope, error)
}
