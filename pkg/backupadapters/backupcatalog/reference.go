package backupcatalog

import (
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

type Reference struct {
	MediaReference catalog.MediaFutureReference
	AlbumReference catalog.AlbumReference
}

func (r Reference) MediaId() string {
	return r.MediaReference.ProvisionalMediaId.Value()
}

func (r Reference) AlbumCreated() bool {
	if r.AlbumReference.AlbumId == nil {
		return false
	}

	return r.AlbumReference.AlbumJustCreated
}

func (r Reference) AlbumFolderName() string {
	if r.AlbumReference.AlbumId == nil {
		return ""
	}

	return r.AlbumReference.AlbumId.FolderName.String()
}

func (r Reference) Exists() bool {
	return r.MediaReference.AlreadyExists
}

func (r Reference) UniqueIdentifier() string {
	return r.MediaReference.Signature.Value()
}
