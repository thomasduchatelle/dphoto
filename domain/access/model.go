package access

import "time"

const (
	//Viewer Permission = "viewer" // Viewer is an account allowed to access medias shared with him but not backup his own
	//Owner2 Permission = "owner"  // Owner is an account super seeding Viewer and who's allowed to back up his media, shared them
	//Admin  Permission = "admin"  // Admin is an account allowed to manage users

	ApiRole   PermissionType = "api"
	OwnerRole PermissionType = "owner"
	AlbumRole PermissionType = "album"
	MediaRole PermissionType = "media"
)

// Permission is defining a set of permissions a user has
//type Permission string

// PermissionType is defining a set of permissions a user has
type PermissionType string

// Permission is attached to a user (a consumer of the API) and define the role it has on resource basis
type Permission struct {
	Type          PermissionType // Type is mandatory, it defines what fields on this structure is used and allow to filter the results
	Role          string         // Role is mandatory, it defines what's the actor can perform
	ResourceOwner string         // ResourceOwner is optional but used for all catalog resources
	ResourceId    string         // ResourceId is optional but used for all catalog resources
	ResourceName  string         // ResourceName is optional and is not used to identify a resource.
	GrantedAt     time.Time      // GrantedAt is the date the role was created for the first time
}

