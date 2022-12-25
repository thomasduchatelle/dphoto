package catalogaclview

import (
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// ListMediasFromAlbum returns medias contained on the album if user is allowed
func (v *View) ListMediasFromAlbum(owner, folderName string) (*catalog.MediaPage, error) {
	err := v.CatalogRules.CanListMediasFromAlbum(owner, folderName)
	if err != nil {
		return nil, err
	}

	return v.CatalogAdapter.ListMedias(owner, folderName, catalog.PageRequest{})
}
