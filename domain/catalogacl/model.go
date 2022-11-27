// Package catalogacl is a layer on top of catalog (same business model) which apply ACL rules before performing any action
package catalogacl

import (
	"github.com/thomasduchatelle/dphoto/domain/catalog"
)

type View struct {
	UserEmail     string               // UserEmail from which the catalog is seen
	AccessControl AccessControlAdapter // AccessControl is a proxy to the accesscontrol domain
}

type AccessControlAdapter interface {
	Owner() (string, error)
	SharedWithUserAlbum() ([]catalog.AlbumId, error)
	SharedByUserGrid(owner string) (map[string][]string, error)
	CanListMediasFromAlbum(owner string, folderName string) error
	CanReadMedia(owner string, id string) error
}

type ListAlbumsFilter struct {
	OnlyDirectlyOwned bool // OnlyDirectlyOwned provides a sub-view where only resources directly owned by user are displayed and accessible
}

type AlbumInView struct {
	*catalog.Album
	SharedTo []string // SharedTo is the list of emails to which this album is shared
}

// Sharing is caring.
type Sharing struct {
	MainUser   string // MainUser is the user having the "owner:main" grant
	Owner      string
	FolderName string
}
