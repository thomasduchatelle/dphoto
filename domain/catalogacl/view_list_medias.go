package catalogacl

import (
	"github.com/thomasduchatelle/dphoto/domain/catalog"
)

// ListMediasFromAlbum returns medias contained on the album if user is allowed
func (v *View) ListMediasFromAlbum(owner, folderName string) (*catalog.MediaPage, error) {
	err := v.AccessControl.CanListMediasFromAlbum(owner, folderName)
	if err != nil {
		return nil, err
	}

	return catalog.ListMedias(owner, folderName, catalog.PageRequest{})
}
