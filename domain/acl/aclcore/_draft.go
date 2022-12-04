package aclcore

// GrantRepository stores the permissions attached to a user.
type GrantRepository interface {
	// IsGranted takes a resource with the form 'drn:<type>:<owner>:<id>'
	IsGranted(email string, resource Permission) (bool, error)

	// Grants allows a user to access a resource that doesn't belong to him
	Grants(email string, resource access.Permission) error

	// Revokes revokes access to a resource for a given users
	Revokes(email string, resource access.Permission) error

	// CountUserPermissions counts number of grants of a certain resource
	CountUserPermissions(email string, resourceType ScopeType) (int, time.Time, time.Time, error)

	// ListResourcesPermissionsByOwner lists permissions about resources owned buy given owners that have been grated to others
	ListResourcesPermissionsByOwner(owners []string, types ...ScopeType) ([]*Scope, error)
}
