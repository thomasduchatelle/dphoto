package access

import "time"

// IsGranted takes a resource with the form 'drn:<type>:<owner>:<id>'
func IsGranted(email string, resource Resource) (bool, error) {
	return false, nil
}

// Grants allows a user to access a resource that doesn't belong to him
func Grants(email string, resource Resource) error {
	return nil
}

// Revokes revokes access to a resource for a given users
func Revokes(email string, resource Resource) error {
	return nil
}

// ListGrants returns all access of a certain type that have been granted to a user
func ListGrants(email string, resourceType ...ResourceType) ([]*Resource, error) {
	return nil, nil
}

// CountGrants counts number of grants of a certain resource
func CountGrants(email string, resourceType ResourceType) (int, time.Time, time.Time, error) {
	return 0, time.Time{}, time.Time{}, nil
}
