// Package catalogaclview is a layer on top of catalog (same business model) which apply ACL rules before performing any action
package catalogaclview

import (
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

type View struct {
	UserEmail      string                  // UserEmail from which the catalog is seen
	CatalogRules   catalogacl.CatalogRules // CatalogRules is rules to use to authorise or deny accesses
	CatalogAdapter ACLViewCatalogAdapter   // CatalogAdapter is a proxy to the catalog domain
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

type ACLViewCatalogAdapter interface {
	FindAllAlbums(owner string) ([]*catalog.Album, error)
	FindAlbums(keys []catalog.AlbumId) ([]*catalog.Album, error)
	ListMedias(owner string, folderName string, request catalog.PageRequest) (*catalog.MediaPage, error)
}
