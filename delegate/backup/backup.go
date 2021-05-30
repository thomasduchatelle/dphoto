package backup

import (
	"duchatelle.io/dphoto/dphoto/scanner"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

// StartBackupRunner starts backup of given scanner.VolumeToBackup and returns when finished. Listeners will received
// progress updates.
func StartBackupRunner(volume scanner.VolumeToBackup, listeners ...interface{}) (*Tracker, error) {
	unsafeChar := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	backupId := fmt.Sprintf("%s_%s", strings.Trim(unsafeChar.ReplaceAllString(volume.UniqueId, "_"), "_"), time.Now().Format("20060102_150405"))
	mdc := log.WithFields(log.Fields{
		"BackupId":   backupId,
		"VolumeId":   volume.UniqueId,
		"VolumePath": volume.Path,
	})

	source, ok := scanner.SourceAdapters[volume.Type]
	if !ok {
		return nil, errors.Errorf("No scanner implementation provided for volume type %s", volume.Type)
	}

	mediaFilter, err := newMediaFilter(&volume)
	if err != nil {
		return nil, err
	}

	uploader, err := NewUploader(new(CatalogProxy), OnlineStorage)
	if err != nil {
		return nil, err
	}

	downloader := Downloader.DownloadMedia
	if volume.Local {
		downloader = PassThroughDownload
	}

	r := scanner.Runner{
		MDC: mdc,
		Source: func(medias chan scanner.FoundMedia) (uint, uint, error) {
			count, size, err := source.FindMediaRecursively(volume, medias)
			mdc.Debugln("Source > volume scanning complete.")
			return count, size, err
		},
		Filter:     mediaFilter.Filter,
		Downloader: downloader,
		Analyser:   scanner.AnalyseMedia,
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
	doneChannel, progressChannel := scanner.Start(r)
	tracker := NewTracker(progressChannel, listeners)

	report := <-doneChannel
	<-tracker.Done

	for i, err := range report.Errors {
		mdc.WithError(err).Errorf("Error %d/%d: %s", i+1, len(report.Errors), err.Error())
	}

	if len(report.Errors) > 0 {
		return nil, errors.Wrapf(report.Errors[0], "Backup failed, %d errors reported until shutdown.", len(report.Errors))
	}

	mdc.Infoln("Backup completed.")
	return tracker, nil
}
