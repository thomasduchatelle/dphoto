package catalogviews

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type LoggingPutSummariesObserver struct {
}

func (l *LoggingPutSummariesObserver) PutSummaries(ctx context.Context, summaries []AlbumSummaryForUsers) error {
	log.Infof("Updating album summaries for %d albums", len(summaries))
	for _, summary := range summaries {
		log.Infof("Album summary: %s", summary)
	}
	return nil
}
