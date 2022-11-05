package accesscontrol

import "time"

const (
	ApiRole   PermissionType = "api"
	OwnerRole PermissionType = "owner"
	AlbumRole PermissionType = "album"
	MediaRole PermissionType = "media"
)

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
	GrantedTo     string         // GrantedTo is the consumer to which this has been shared with
}
