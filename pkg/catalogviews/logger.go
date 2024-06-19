package catalogviews

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type LoggingInsertAlbumSizeObserver struct {
}

func (l *LoggingInsertAlbumSizeObserver) InsertAlbumSize(ctx context.Context, albumSizes []MultiUserAlbumSize) error {
	log.Infof("Updating album sizes for %d albums", len(albumSizes))
	for _, albumSize := range albumSizes {
		log.Infof("Album size: %s", albumSize)
	}
	return nil
}
