package catalogacl

import (
	"github.com/thomasduchatelle/dphoto/domain/catalog"
)

type directlyOwnedView struct {
}

func (d directlyOwnedView) CanReadMedia(owner string, id string) error {
	//TODO implement me
	panic("implement me")
}

func (d directlyOwnedView) ListAlbums(email string) ([]*AlbumInView, error) {
	albums, err := catalog.ListAlbums(catalog.ListAlbumsInput{
		Owner: email,
	})

	var view []*AlbumInView
	for _, album := range albums {
		view = append(view, &AlbumInView{
			Album:    album,
			SharedTo: nil,
			SharedBy: "",
		})
	}
	return view, err
}

func (d directlyOwnedView) ListMediasFromAlbum(owner, album string) (*catalog.MediaPage, error) {
	return catalog.ListMedias(owner, album, catalog.PageRequest{})
}
