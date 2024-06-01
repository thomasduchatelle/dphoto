package catalog

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// ListMedias return a page of medias within an album
func ListMedias(albumId AlbumId, request PageRequest) (*MediaPage, error) {
	medias, err := repositoryPort.FindMedias(context.TODO(), NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
	return &MediaPage{
		Content: medias,
	}, err
}

// FindMediaOwnership returns the folderName containing the media, or AlbumNotFoundError.
func FindMediaOwnership(owner ownermodel.Owner, mediaId MediaId) (*AlbumId, error) {
	return repositoryPort.FindMediaCurrentAlbum(context.TODO(), owner, mediaId)
}
