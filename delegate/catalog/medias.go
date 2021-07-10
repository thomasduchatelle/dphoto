package catalog

// ListMedias return a page of medias within an album
func ListMedias(folderName string, request PageRequest) (*MediaPage, error) {
	return Repository.FindMedias(normaliseFolderName(folderName), FindMediaFilter{
		PageRequest: request,
	})
}

// InsertMedias stores metadata and location of photo and videos
func InsertMedias(medias []CreateMediaRequest) error {
	return Repository.InsertMedias(medias)
}

// FindSignatures returns a list of the medias already known ; they can't be duplicated
func FindSignatures(signatures []*MediaSignature) ([]*MediaSignature, error) {
	return Repository.FindExistingSignatures(signatures)
}
