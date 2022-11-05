// Package catalogacl is a layer on top of catalog (same business model) which apply ACL rules before performing any action
package catalogacl

import (
	"github.com/thomasduchatelle/dphoto/domain/catalog"
)

type ListAlbumsFilter struct {
	OnlyDirectlyOwned bool // OnlyDirectlyOwned provides a sub-view where only resources directly owned by user are displayed and accessible
}

type View interface {
	// ListAlbums returns albums visible by the user (owned by current user, and shared to him)
	ListAlbums(filter ListAlbumsFilter) ([]*AlbumInView, error)

	// ListMediasFromAlbum returns medias contained on the album if user is allowed
	ListMediasFromAlbum(owner, album string) (*catalog.MediaPage, error)

	// CanReadMedia returns an error if the user is not allowed
	CanReadMedia(owner string, id string) error
}

type AlbumInView struct {
	*catalog.Album
	SharedTo []string // SharedTo is the list of emails to which this album is shared
	SharedBy string   // SharedBy is the email from whom the album has been shared
}

func NewUserView(user string) View {
	return &viewimpl{
		userEmail: user,
	}
}

type viewimpl struct {
	userEmail string

	ownerOf []string // ownerOf usually contains the userEmail: it's the representation of the user as owner IDs
}

func (v *viewimpl) CanReadMedia(owner string, id string) error {
	//TODO implement me
	panic("implement me")
}
