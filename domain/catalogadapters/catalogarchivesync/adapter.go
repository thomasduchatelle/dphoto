// Package catalogarchivesync is calling archive domain directly
package catalogarchivesync

import (
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
)

func New() catalog.CArchiveAdapter {
	return new(adapter)
}

type adapter struct {
}

func (a *adapter) MoveMedias(owner string, ids []string, targetFolder string) error {
	return archive.Relocate(owner, ids, targetFolder)
}
