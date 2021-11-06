// Package downloader provides an alternative to a local temporary storage.
package downloader

import "github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"

func PassThroughDownload(media backupmodel.FoundMedia) (backupmodel.FoundMedia, error) {
	return media, nil
}
