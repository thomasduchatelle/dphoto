// Package repositoryjson uses JSON files stored in local filesystem ; must be used with JSON recorder to re-build the database when requested
package repositoryjson

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
)

func New() (catalog.RepositoryAdapter, error) {
	// TODO Rebuild the in-memory database
	return &adapter{}, nil
}

type adapter struct {
	albums map[string]map[string]*catalog.Album
}

func (a *adapter) FindAllAlbums(owner string) ([]*catalog.Album, error) {
	var albums []*catalog.Album
	for _, album := range a.albumsForOwner(owner) {
		cp := *album
		albums = append(albums, &cp)
	}

	return albums, nil
}

func (a *adapter) InsertAlbum(album catalog.Album) error {
	if _, exists := a.albumsForOwner(album.Owner)[album.FolderName]; exists {
		return errors.Errorf("Album for %s with folder name %s already exists", album.Owner, album.FolderName)
	}

	a.albumsForOwner(album.Owner)[album.FolderName] = &album
	return nil
}

func (a *adapter) DeleteEmptyAlbum(owner string, folderName string) error {
	delete(a.albumsForOwner(owner), folderName)
	return nil
}

func (a *adapter) FindAlbum(owner string, folderName string) (*catalog.Album, error) {
	album, found := a.albumsForOwner(owner)[folderName]
	if !found {
		return nil, catalog.NotFoundError
	}

	cp := *album
	return &cp, nil
}

func (a *adapter) UpdateAlbum(album catalog.Album) error {
	a.albumsForOwner(album.Owner)[album.FolderName] = &album
	return nil
}

func (a *adapter) InsertMedias(owner string, media []catalog.CreateMediaRequest) error {
	//TODO implement me
	panic("implement me")
}

func (a *adapter) FindMedias(request *catalog.FindMediaRequest) (medias []*catalog.MediaMeta, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *adapter) FindMediaIds(request *catalog.FindMediaRequest) (ids []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *adapter) FindExistingSignatures(owner string, signatures []*catalog.MediaSignature) ([]*catalog.MediaSignature, error) {
	//TODO implement me
	panic("implement me")
}

func (a *adapter) TransferMedias(owner string, mediaIds []string, newFolderName string) error {
	//TODO implement me
	panic("implement me")
}

func (a *adapter) albumsForOwner(owner string) map[string]*catalog.Album {
	albums, ok := a.albums[owner]
	if !ok {
		albums = make(map[string]*catalog.Album)
		a.albums[owner] = albums
	}

	return albums
}
