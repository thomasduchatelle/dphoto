package accesscontrol

import (
	"github.com/thomasduchatelle/dphoto/domain/access"
	"time"
)

// IsGranted takes a resource with the form 'drn:<type>:<owner>:<id>'
func IsGranted(email string, resource access.Permission) (bool, error) {
	return false, nil
}

// Grants allows a user to access a resource that doesn't belong to him
func Grants(email string, resource access.Permission) error {
	return nil
}

// Revokes revokes access to a resource for a given users
func Revokes(email string, resource access.Permission) error {
	return nil
}

// ListUserPermissions returns all access of a certain type that have been granted to a user
func ListUserPermissions(email string, types ...PermissionType) ([]*Permission, error) {
	return nil, nil
}

// CountUserPermissions counts number of grants of a certain resource
func CountUserPermissions(email string, resourceType PermissionType) (int, time.Time, time.Time, error) {
	return 0, time.Time{}, time.Time{}, nil
}

// ListResourcesPermissionsByOwner lists permissions about resources owned buy given owners that have been grated to others
func ListResourcesPermissionsByOwner(owners []string, types ...PermissionType) ([]*Permission, error) {
	return nil, nil
}
