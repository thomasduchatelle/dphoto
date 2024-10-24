// Package backup is providing commands to inspect a file system (hard-drive, USB, Android, S3) and backup medias to a remote DPhoto storage.
package backup

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"regexp"
	"strings"
	"time"
)

type SourceVolume interface {
	String() string
	FindMedias() ([]FoundMedia, error)
}

// Backup is analysing each media and is backing it up if not already in the catalog.
func Backup(owner ownermodel.Owner, volume SourceVolume, optionsSlice ...Options) (CompletionReport, error) {
	unsafeChar := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	backupId := fmt.Sprintf("%s_%s", strings.Trim(unsafeChar.ReplaceAllString(volume.String(), "_"), "_"), time.Now().Format("20060102_150405"))
	mdc := log.WithFields(log.Fields{
		"BackupId": backupId,
		"Volume":   volume.String(),
	})

	options := ReduceOptions(optionsSlice...)

	cataloger, err := NewCataloguer(owner, options)
	if err != nil {
		return nil, err
	}

	publisher, hintSize, err := newPublisher(volume)

	run := runner{
		MDC:          mdc,
		Options:      options,
		Publisher:    publisher,
		Analyser:     getDefaultAnalyser(),
		Cataloger:    cataloger,
		UniqueFilter: newUniqueFilter(),
		Uploader:     &Uploader{Owner: owner, InsertMediaPort: insertMediaPort},
	}

	progressChannel, _ := run.start(context.TODO(), hintSize)
	backupReport := NewTracker(progressChannel, options.Listener)

	err = run.waitToFinish()
	backupReport.WaitToComplete()

	if err == nil {
		mdc.Infof("Backup completed, %d medias backed up.", backupReport.MediaCount())
	} else {
		mdc.WithError(err).Errorf("Backup faifed with err: %s", err.Error())
	}
	return backupReport, err
}

func defaultValue(value, fallback int) int {
	if value == 0 {
		return fallback
	}

	return value
}
