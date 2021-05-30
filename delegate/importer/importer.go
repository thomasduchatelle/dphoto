package importer

import (
	"duchatelle.io/dphoto/dphoto/scanner"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

const (
	scanBufferSize = 1024 * 16
)

// StartImportStructure run through all medias and discover its folder structure to suggest creation of albums
func StartImportStructure(volume scanner.VolumeToBackup, listeners ...interface{}) (*Tracker, error) {
	unsafeChar := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	importId := fmt.Sprintf("%s_%s", strings.Trim(unsafeChar.ReplaceAllString(volume.UniqueId, "_"), "_"), time.Now().Format("20060102_150405"))
	mdc := log.WithFields(log.Fields{
		"ImportId":   importId,
		"VolumeId":   volume.UniqueId,
		"VolumePath": volume.Path,
	})

	source, ok := scanner.SourceAdapters[volume.Type]
	if !ok {
		return nil, errors.Errorf("No scanner implementation provided for volume type %s", volume.Type)
	}

	structure := new(StructureDiscovery)
	r := scanner.Runner{
		MDC: mdc,
		Source: func(medias chan scanner.FoundMedia) (uint, uint, error) {
			count, size, err := source.FindMediaRecursively(volume, medias)
			mdc.Debugf("Source > volume scanning complete.")
			return count, size, err
		},
		Filter: func(found scanner.FoundMedia) bool {
			return true
		},
		Downloader:           scanner.PassThroughDownload,
		Analyser:             scanner.AnalyseMedia,
		Uploader:             structure.StructureDiscovery,
		PreCompletion:        nil,
		FoundMediaBufferSize: scanBufferSize,
		BufferSize:           scanBufferSize,
		ConcurrentDownloader: 1,
		ConcurrentAnalyser:   2,
		ConcurrentUploader:   1,
		UploadBatchSize:      scanBufferSize,
	}
	doneChannel, progressChannel := scanner.Start(r)
	tracker := NewTracker(progressChannel, listeners)

	report := <-doneChannel
	<-tracker.Done

	for i, err := range report.Errors {
		mdc.WithError(err).Errorf("Error %d/%d: %s", i+1, len(report.Errors), err.Error())
	}

	if len(report.Errors) > 0 {
		return nil, errors.Wrapf(report.Errors[0], "Import failed, %d errors reported until shutdown.", len(report.Errors))
	}

	mdc.Infoln("Import completed.")
	return tracker, nil

}

type Tracker struct {
	Done                chan struct{} // Done is closed when all events have been processed.
}

// NewTracker creates the Tracker and start consuming (async)
func NewTracker(progressChannel chan *scanner.ProgressEvent) *Tracker {
	tracker := &Tracker{
		Done:          make(chan struct{}),
	}
	go func() {
		defer close(tracker.Done)
		tracker.consume(progressChannel)
	}()
	return tracker
}

func (t *Tracker) consume(progressChannel chan *scanner.ProgressEvent) {
	for event := range progressChannel {
		switch event.Type {
		case scanner.ProgressEventScanComplete:
		case scanner.ProgressEventUploaded:
			t.uploaded = t.uploaded.Add(event.Count, event.Size)

			typeCount, ok := t.detailedCount[event.Album]
			if !ok {
				typeCount = &TypeCounter{}
				t.detailedCount[event.Album] = typeCount
			}
			typeCount.incrementFoundCounter(event.MediaType, event.Count, event.Size)

			t.fireUploadedEvent()
		}
	}
}