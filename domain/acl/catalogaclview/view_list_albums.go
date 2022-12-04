package catalogaclview

import (
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"sort"
)

// ListAlbums returns albums visible by the user (owned by current user, and shared to him)
func (v *View) ListAlbums(filter ListAlbumsFilter) ([]*AlbumInView, error) {
	view, err := v.listOwnedAlbums()
	if err != nil {
		return nil, err
	}

	if !filter.OnlyDirectlyOwned {
		albums, err := v.listSharedWithUserAlbums()
		if err != nil {
			return nil, err
		}

		view = append(view, albums...)
		sort.Slice(view, func(i, j int) bool {
			return view[i].Start.Before(view[j].Start)
		})
	}

	return view, nil
}

func (v *View) listOwnedAlbums() ([]*AlbumInView, error) {
	owner, err := v.CatalogRules.Owner()
	if err != nil || owner == "" {
		return nil, err
	}

	ownedAlbums, err := v.CatalogAdapter.FindAllAlbums(owner)
	if err != nil {
		return nil, err
	}

	sharing, err := v.CatalogRules.SharedByUserGrid(owner)

	var view []*AlbumInView
	for _, album := range ownedAlbums {
		sharedTo, _ := sharing[album.FolderName]
		view = append(view, &AlbumInView{
			Album:    album,
			SharedTo: sharedTo,
		})
	}

	return view, err
}

func (v *View) listSharedWithUserAlbums() ([]*AlbumInView, error) {
	shares, err := v.CatalogRules.SharedWithUserAlbum()
	if err != nil || len(shares) == 0 {
		return nil, err
	}

	ids := make([]catalog.AlbumId, len(shares), len(shares))
	for i, alb := range shares {
		ids[i] = catalog.AlbumId{
			Owner:      alb.Owner,
			FolderName: alb.FolderName,
		}
	}

	albums, err := v.CatalogAdapter.FindAlbums(ids)

	var view []*AlbumInView
	for _, album := range albums {
		view = append(view, &AlbumInView{
			Album: album,
		})
	}

	return view, err
}
