package backup

import (
	"duchatelle.io/dphoto/dphoto/scanner"
)

func PassThroughDownload(media scanner.FoundMedia) (scanner.FoundMedia, error) {
	return media, nil
}
