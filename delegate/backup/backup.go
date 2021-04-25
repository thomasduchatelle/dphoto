package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"duchatelle.io/dphoto/dphoto/backup/runner"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

func StartBackupRunner(volume model.VolumeToBackup) error {
	backupId := fmt.Sprintf("%s_%s", volume.UniqueId, time.Now().Format("20060102_150405"))
	mdc := log.WithFields(log.Fields{
		"BackupId":   backupId,
		"VolumeId":   volume.UniqueId,
		"VolumePath": volume.Path,
	})

	scanner, ok := ScannerAdapters[volume.Type]
	if !ok {
		return errors.Errorf("No scanner implementation provided for volume type %s", volume.Type)
	}

	mediaFilter, err := newMediaFilter(&volume)
	if err != nil {
		return err
	}

	uploader, err := NewUploader(new(CatalogProxy), OnlineStorageFactory())
	if err != nil {
		return err
	}

	r := runner.Runner{
		MDC: mdc,
		Source: func(medias chan model.FoundMedia) error {
			return scanner.FindMediaRecursively(volume, medias)
		},
		Filter:     mediaFilter.Filter,
		Downloader: DownloaderFactory().DownloadMedia,
		Analyser:   analyseMedia,
		Uploader:   uploader.Upload,
		PreCompletion: func() error {
			return mediaFilter.StoreState(backupId)
		},
		BufferSize:           UploadBatchSize,
		ConcurrentDownloader: DownloadThreadCount,
		ConcurrentAnalyser:   ImageReaderThreadCount,
		ConcurrentUploader:   UploadThreadCount,
		UploadBatchSize:      UploadBatchSize,
	}
	doneChannel := runner.Start(r)

	report := <-doneChannel
	for i, err := range report.Errors {
		mdc.WithError(err).Errorf("Error %d/%d: %s", i+1, len(report.Errors), err.Error())
	}

	if len(report.Errors) > 0 {
		return errors.Wrapf(report.Errors[0], "Backup failed, %d errors reported until shutdown.", len(report.Errors))
	}

	mdc.Infoln("Backup completed.")
	return nil
}
