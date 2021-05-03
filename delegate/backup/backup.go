package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"duchatelle.io/dphoto/dphoto/backup/runner"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

func StartBackupRunner(volume model.VolumeToBackup) error {
	unsafeChar := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	backupId := fmt.Sprintf("%s_%s", strings.Trim(unsafeChar.ReplaceAllString(volume.UniqueId, "_"), "_"), time.Now().Format("20060102_150405"))
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

	uploader, err := NewUploader(new(CatalogProxy), OnlineStorage)
	if err != nil {
		return err
	}

	downloader := Downloader.DownloadMedia
	if volume.Local {
		downloader = PassThroughDownload
	}

	r := runner.Runner{
		MDC: mdc,
		Source: func(medias chan model.FoundMedia) error {
			err := scanner.FindMediaRecursively(volume, medias)
			mdc.Debugf("Source > Incoming volume scanning complete.")
			return err
		},
		Filter:     mediaFilter.Filter,
		Downloader: downloader,
		Analyser:   analyseMedia,
		Uploader:   uploader.Upload,
		PreCompletion: func() error {
			return mediaFilter.StoreState(backupId)
		},
		BufferSize:           uploadBatchSize,
		ConcurrentDownloader: downloadThreadCount,
		ConcurrentAnalyser:   imageReaderThreadCount,
		ConcurrentUploader:   uploadThreadCount,
		UploadBatchSize:      uploadBatchSize,
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
