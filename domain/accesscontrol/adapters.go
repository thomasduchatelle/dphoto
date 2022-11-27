package accesscontrol

import (
	"github.com/thomasduchatelle/dphoto/domain/access"
	"time"
)

type ScopesReader interface {
	// ListUserScopes returns all access of a certain type that have been granted to a user
	ListUserScopes(email string, types ...ScopeType) ([]*Scope, error)
}

type ReverseScopesReader interface {
	// ListOwnerScopes is a reverse query to find to whom has been shared owner resources.
	ListOwnerScopes(owner string, types ...ScopeType) ([]*Scope, error)
}

// GrantRepository stores the permissions attached to a user.
type GrantRepository interface {
	// IsGranted takes a resource with the form 'drn:<type>:<owner>:<id>'
	IsGranted(email string, resource access.Permission) (bool, error)

	// Grants allows a user to access a resource that doesn't belong to him
	Grants(email string, resource access.Permission) error

	// Revokes revokes access to a resource for a given users
	Revokes(email string, resource access.Permission) error

	// CountUserPermissions counts number of grants of a certain resource
	CountUserPermissions(email string, resourceType ScopeType) (int, time.Time, time.Time, error)

	// ListResourcesPermissionsByOwner lists permissions about resources owned buy given owners that have been grated to others
	ListResourcesPermissionsByOwner(owners []string, types ...ScopeType) ([]*Scope, error)
}
