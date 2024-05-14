package catalog

import (
	"context"
)

// ListMedias return a page of medias within an album
func ListMedias(albumId AlbumId, request PageRequest) (*MediaPage, error) {
	medias, err := repositoryPort.FindMedias(context.TODO(), NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
	return &MediaPage{
		Content: medias,
	}, err
}

// FindSignatures returns a list of the medias already known ; they can't be duplicated
func FindSignatures(owner Owner, signatures []*MediaSignature) ([]*MediaSignature, error) {
	return repositoryPort.FindExistingSignatures(context.TODO(), owner, signatures)
}

// FindMediaOwnership returns the folderName containing the media, or AlbumNotFoundError.
func FindMediaOwnership(owner Owner, mediaId MediaId) (*AlbumId, error) {
	return repositoryPort.FindMediaCurrentAlbum(context.TODO(), owner, mediaId)
}
