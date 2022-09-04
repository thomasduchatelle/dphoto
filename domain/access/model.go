package access

import "time"

const (
	Viewer Role = "viewer" // Viewer is an account allowed to access medias shared with him but not backup his own
	Owner2 Role = "owner"  // Owner is an account super seeding Viewer and who's allowed to back up his media, shared them
	Admin  Role = "admin"  // Admin is an account allowed to manage users

	SoftwareResource ResourceType = "software"
	OwnerResource    ResourceType = "owner"
	AlbumResource    ResourceType = "album"
	MediaResource    ResourceType = "media"
)

// Role is defining a set of permissions a user has
type Role string

// ResourceType is defining a set of permissions a user has
type ResourceType string

// Resource is representing an album, a media, or something else
type Resource struct {
	Type  ResourceType
	Owner string
	Id    string
	Name  string // Name is optional and is not used to identify a resource.
}

type Grant struct {
	Resource  Resource
	GrantedAt time.Time
}

func newSoftwareResource(role string) Resource {
	return Resource{
		Type: SoftwareResource,
		Id:   role,
	}
}
