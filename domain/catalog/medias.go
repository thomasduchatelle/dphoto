package catalog

// ListMedias return a page of medias within an album
func ListMedias(owner string, folderName string, request PageRequest) (*MediaPage, error) {
	return Repository.FindMedias(owner, normaliseFolderName(folderName), FindMediaFilter{
		PageRequest: request,
	})
}

// InsertMedias stores metadata and location of photo and videos
func InsertMedias(owner string, medias []CreateMediaRequest) error {
	for _, m := range medias {
		m.Location.FolderName = normaliseFolderName(m.Location.FolderName)
	}
	return Repository.InsertMedias(owner, medias)
}

// FindSignatures returns a list of the medias already known ; they can't be duplicated
func FindSignatures(owner string, signatures []*MediaSignature) ([]*MediaSignature, error) {
	return Repository.FindExistingSignatures(owner, signatures)
}

func GetMediaLocations(owner string, signature MediaSignature) ([]*MediaLocation, error) {
	return Repository.FindMediaLocations(owner, signature)
}
