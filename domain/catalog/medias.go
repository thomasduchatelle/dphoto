package catalog

// ListMedias return a page of medias within an album
func ListMedias(owner string, folderName string, request PageRequest) (*MediaPage, error) {
	medias, err := repositoryPort.FindMedias(NewFindMediaRequest(owner).WithAlbum(normaliseFolderName(folderName)))
	return &MediaPage{
		Content: medias,
	}, err
}

// InsertMedias stores metadata and location of photo and videos
func InsertMedias(owner string, medias []CreateMediaRequest) error {
	return repositoryPort.InsertMedias(owner, medias)
}

// FindSignatures returns a list of the medias already known ; they can't be duplicated
func FindSignatures(owner string, signatures []*MediaSignature) ([]*MediaSignature, error) {
	return repositoryPort.FindExistingSignatures(owner, signatures)
}

// AssignIdsToNewMedias filters out signatures that are already known and compute a unique ID for the others.
func AssignIdsToNewMedias(owner string, signatures []*MediaSignature) (map[MediaSignature]string, error) {
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

// FindMediaOwnership returns the folderName containing the media, or NotFoundError.
func FindMediaOwnership(owner, mediaId string) (string, error) {
	return repositoryPort.FindMediaCurrentAlbum(owner, mediaId)
}
