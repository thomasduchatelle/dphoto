package catalogaclview

import (
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// ListMediasFromAlbum returns medias contained on the album if user is allowed
func (v *View) ListMediasFromAlbum(albumId catalog.AlbumId) (*catalog.MediaPage, error) {
	err := v.CatalogRules.CanListMediasFromAlbum(albumId)
	if err != nil {
		return nil, err
	}

	return v.CatalogAdapter.ListMedias(albumId, catalog.PageRequest{})
}
