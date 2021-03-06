package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/interactors"
	"duchatelle.io/dphoto/dphoto/backup/interactors/analyser"
	"duchatelle.io/dphoto/dphoto/backup/interactors/downloader"
	"duchatelle.io/dphoto/dphoto/backup/interactors/filter"
	"duchatelle.io/dphoto/dphoto/backup/interactors/runner"
	"duchatelle.io/dphoto/dphoto/backup/interactors/tracker"
	"duchatelle.io/dphoto/dphoto/backup/interactors/uploaders"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

// StartBackupRunner starts backup of given model.VolumeToBackup and returns when finished. Listeners will received
// progress updates.
func StartBackupRunner(volume model.VolumeToBackup, listeners ...interface{}) (model.BackupReport, error) {
	unsafeChar := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	backupId := fmt.Sprintf("%s_%s", strings.Trim(unsafeChar.ReplaceAllString(volume.UniqueId, "_"), "_"), time.Now().Format("20060102_150405"))
	mdc := log.WithFields(log.Fields{
		"BackupId":   backupId,
		"VolumeId":   volume.UniqueId,
		"VolumePath": volume.Path,
	})

	source, err := interactors.NewSource(volume, func(count, size uint) {
		mdc.Debugf("Source > volume scanning complete.")
	})
	if err != nil {
		return nil, err
	}

	mediaFilter, err := filter.NewMediaFilter(&volume)
	if err != nil {
		return nil, err
	}

	uploader, err := uploaders.NewUploader(new(uploaders.CatalogProxy), interactors.OnlineStoragePort)
	if err != nil {
		return nil, err
	}

	downloaderPort := interactors.DownloaderPort.DownloadMedia
	if volume.Local {
		downloaderPort = downloader.PassThroughDownload
	}

	r := runner.Runner{
		MDC:        mdc,
		Source:     source,
		Filter:     mediaFilter.Filter,
		Downloader: downloaderPort,
		Analyser:   analyser.AnalyseMedia,
		Uploader:   uploader.Upload,
		PreCompletion: func() error {
			return mediaFilter.StoreState(backupId)
		},
		FoundMediaBufferSize: scanBufferSize,
		BufferSize:           uploadBatchSize * uploadThreadCount * 2,
		ConcurrentDownloader: downloadThreadCount,
		ConcurrentAnalyser:   imageReaderThreadCount,
		ConcurrentUploader:   uploadThreadCount,
		UploadBatchSize:      uploadBatchSize,
	}
	doneChannel, progressChannel := runner.Start(r)
	backupReport := tracker.NewTracker(progressChannel, listeners)

	report := <-doneChannel
	backupReport.WaitToComplete()

	for i, err := range report.Errors {
		mdc.WithError(err).Errorf("Error %d/%d: %s", i+1, len(report.Errors), err.Error())
	}

	if len(report.Errors) > 0 {
		return nil, errors.Wrapf(report.Errors[0], "Backup failed, %d errors reported until shutdown.", len(report.Errors))
	}

	mdc.Infoln("Backup completed.")
	return backupReport, nil
}
