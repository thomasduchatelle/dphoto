package backup

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

func newLogger(volumeName string) *logger {
	unsafeChar := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	scanId := fmt.Sprintf("%s_%s", strings.Trim(unsafeChar.ReplaceAllString(volumeName, "_"), "_"), time.Now().Format("20060102_150405"))

	return &logger{
		mdc: log.WithFields(log.Fields{
			"ScanId": scanId,
			"Volume": volumeName,
		}),
	}
}

type logger struct {
	mdc *log.Entry
}

func (l *logger) OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error {
	l.mdc.WithFields(log.Fields{
		"Media":     media.FoundMedia.String(),
		"Reference": reference.AlbumFolderName(),
		"Cause":     cause,
	}).Info("Media filtered out")
	return nil
}

func (l *logger) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	l.mdc.WithFields(log.Fields{
		"Media": found.String(),
		"Cause": cause,
	}).Info("Media rejected")
	return nil
}

func (l *logger) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	l.mdc.WithFields(log.Fields{
		"Media": media.FoundMedia.String(),
	}).Infof("Media analysed %s", media.Details.DateTime.Format(time.DateTime))
	return nil
}

func (l *logger) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	for _, request := range requests {
		l.mdc.Infof("Media catalogued %s <- %s", request.CatalogReference.AlbumFolderName(), request.AnalysedMedia.FoundMedia)
	}

	return nil
}
