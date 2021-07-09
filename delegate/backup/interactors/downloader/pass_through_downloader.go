// Package downloader provides an alternative to a local temporary storage.
package downloader

import "duchatelle.io/dphoto/dphoto/backup/backupmodel"

func PassThroughDownload(media backupmodel.FoundMedia) (backupmodel.FoundMedia, error) {
	return media, nil
}
