// Package backup is providing commands to inspect a file system (hard-drive, USB, Android, S3) and backup medias to a remote DPhoto storage.
package backup

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

type SourceVolume interface {
	String() string
	FindMedias() ([]FoundMedia, error)
}

// Backup is analysing each media and is backing it up if not already in the catalog.
func Backup(owner string, volume SourceVolume, optionsSlice ...Options) (CompletionReport, error) {
	unsafeChar := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	backupId := fmt.Sprintf("%s_%s", strings.Trim(unsafeChar.ReplaceAllString(volume.String(), "_"), "_"), time.Now().Format("20060102_150405"))
	mdc := log.WithFields(log.Fields{
		"BackupId": backupId,
		"Volume":   volume.String(),
	})

	options := readOptions(optionsSlice)

	cataloger, err := chooseCataloger(owner, options)
	if err != nil {
		return nil, err
	}

	publisher, hintSize, err := newPublisher(volume)

	run := runner{
		MDC:                  mdc,
		Publisher:            publisher,
		Analyser:             newBackupAnalyseMedia(),
		Cataloger:            cataloger,
		UniqueFilter:         newUniqueFilter(),
		Uploader:             newBackupUploader(owner),
		ConcurrentAnalyser:   ConcurrentAnalyser,
		ConcurrentCataloguer: ConcurrentCataloguer,
		ConcurrentUploader:   ConcurrentUploader,
		BatchSize:            BatchSize,
	}

	progressChannel, doneChannel := run.start(context.TODO(), hintSize)
	backupReport := NewTracker(progressChannel, []interface{}{options.Listener})

	reportedErrors := <-doneChannel
	backupReport.WaitToComplete()

	for i, err := range reportedErrors {
		mdc.WithError(err).Errorf("Error %d/%d: %s", i+1, len(reportedErrors), err.Error())
	}

	if len(reportedErrors) > 0 {
		return nil, errors.Wrapf(reportedErrors[0], "Backup failed, %d errors reported until shutdown.", len(reportedErrors))
	}

	mdc.Infoln("Backup completed.")
	return backupReport, nil
}

func chooseCataloger(owner string, options Options) (runnerCataloger, error) {
	if len(options.RestrictedAlbumFolderName) > 0 {
		return newReadOnlyCataloger(owner, options.RestrictedAlbumFolderName)
	}

	return newCreatorCataloger(owner)
}
