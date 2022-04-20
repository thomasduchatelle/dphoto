package catalog

import "github.com/thomasduchatelle/dphoto/domain/catalogmodel"

// ListMedias return a page of medias within an album
func ListMedias(owner string, folderName string, request catalogmodel.PageRequest) (*catalogmodel.MediaPage, error) {
	return Repository.FindMedias(owner, normaliseFolderName(folderName), catalogmodel.FindMediaFilter{
		PageRequest: request,
	})
}

// InsertMedias stores metadata and location of photo and videos
func InsertMedias(owner string, medias []catalogmodel.CreateMediaRequest) error {
	for _, m := range medias {
		m.Location.FolderName = normaliseFolderName(m.Location.FolderName)
	}
	return Repository.InsertMedias(owner, medias)
}

// FindSignatures returns a list of the medias already known ; they can't be duplicated
func FindSignatures(owner string, signatures []*catalogmodel.MediaSignature) ([]*catalogmodel.MediaSignature, error) {
	return Repository.FindExistingSignatures(owner, signatures)
}
