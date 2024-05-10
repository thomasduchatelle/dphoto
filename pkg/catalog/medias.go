package catalog

import context "context"

// ListMedias return a page of medias within an album
func ListMedias(albumId AlbumId, request PageRequest) (*MediaPage, error) {
	medias, err := repositoryPort.FindMedias(context.TODO(), NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
	return &MediaPage{
		Content: medias,
	}, err
}

// InsertMedias stores metadata and location of photo and videos
func InsertMedias(owner Owner, medias []CreateMediaRequest) error {
	return repositoryPort.InsertMedias(context.TODO(), owner, medias)
}

// FindSignatures returns a list of the medias already known ; they can't be duplicated
func FindSignatures(owner Owner, signatures []*MediaSignature) ([]*MediaSignature, error) {
	return repositoryPort.FindExistingSignatures(context.TODO(), owner, signatures)
}

// AssignIdsToNewMedias filters out signatures that are already known and compute a unique ID for the others.
func AssignIdsToNewMedias(owner Owner, signatures []*MediaSignature) (map[MediaSignature]string, error) {
	existingSignaturesSlice, err := FindSignatures(owner, signatures)
	if err != nil {
		return nil, err
	}

	existingSignatures := make(map[MediaSignature]interface{})
	for _, sign := range existingSignaturesSlice {
		existingSignatures[*sign] = nil
	}

	assignedIds := make(map[MediaSignature]string)
	for _, sign := range signatures {
		if _, exists := existingSignatures[*sign]; !exists {
			assignedIds[*sign], err = GenerateMediaId(*sign)
			if err != nil {
				return nil, err
			}
		}
	}

	return assignedIds, nil
}

// FindMediaOwnership returns the folderName containing the media, or AlbumNotFoundError.
func FindMediaOwnership(owner Owner, mediaId MediaId) (*AlbumId, error) {
	return repositoryPort.FindMediaCurrentAlbum(context.TODO(), owner, mediaId)
}
