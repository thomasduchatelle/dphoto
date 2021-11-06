// Package downloader provides an alternative to a local temporary storage.
package downloader

import "github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"

func PassThroughDownload(media backupmodel.FoundMedia) (backupmodel.FoundMedia, error) {
	return media, nil
}
