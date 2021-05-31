package downloader

import "duchatelle.io/dphoto/dphoto/backup/model"

func PassThroughDownload(media model.FoundMedia) (model.FoundMedia, error) {
	return media, nil
}
