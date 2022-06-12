package backup

import (
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"time"
)

func GetPreAuthorisedUrl(owner string, locations []*catalog.MediaLocation, expires time.Duration) (string, error) {
	// note - it is assumed the first location is the right one!
	return Storage.ContentSignedUrl(owner, locations[0].FolderName, locations[0].Filename, expires)
}
