package runner

import (
	"github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"
	log "github.com/sirupsen/logrus"
	"sync"
)

// Runner is a workflow engine that filter, download, analyse, and upload medias using several goroutines.
// Workflow is stopped at the first error but might have several while channels are de-stacked.
type Runner struct {
	MDC                  *log.Entry // MDC is log.WithFields({}) that contains Mapped Diagnostic Context
	Source               Source
	Filter               Filter
	Downloader           Downloader
	Analyser             Analyser
	Uploader             Uploader
	PreCompletion        PreCompletion
	FoundMediaBufferSize int // FoundMediaBufferSize is the size of the scanner buffer, should be BIG in order to let the scan finish (and progress bars to give an accurate estimation)
	BufferSize           int // BufferSize is the default size for channels
	ConcurrentDownloader int // ConcurrentDownloader is the number of goroutines that will filter and download files ; should be 1 with a pass-through downloaders
	ConcurrentAnalyser   int // ConcurrentAnalyser is the number of goroutines that analyse the medias
	ConcurrentUploader   int // ConcurrentUploader is the number of goroutine that upload files online
	UploadBatchSize      int // UploadBatchSize is the size of the buffer for the uploader
	report               *Report
	completionChannel    chan *Report
	progressEventChannel chan *backupmodel.ProgressEvent
}

func Start(runner Runner) (chan *Report, chan *backupmodel.ProgressEvent) {
	runner.report = &Report{}
	runner.completionChannel = make(chan *Report, 1)
	runner.progressEventChannel = make(chan *backupmodel.ProgressEvent, runner.BufferSize)

	foundChannel := make(chan backupmodel.FoundMedia, runner.FoundMediaBufferSize)
	readyToAnalyseChannel := make(chan backupmodel.FoundMedia, runner.BufferSize)
	readyToBackupChannel := make(chan *backupmodel.AnalysedMedia, runner.BufferSize)

	runner.pipeFoundToReadyToAnalyseChannels(foundChannel, readyToAnalyseChannel)
	runner.pipeReadyToAnalyseToReadyToBackupChannels(readyToAnalyseChannel, readyToBackupChannel)
	runner.pipeReadyToBackupToCompletedChannels(readyToBackupChannel, runner.completionChannel, runner.progressEventChannel)

	runner.startScanning(foundChannel)

	return runner.completionChannel, runner.progressEventChannel
}

func (r *Runner) pipeFoundToReadyToAnalyseChannels(foundCh, downloadedCh chan backupmodel.FoundMedia) {
	group := sync.WaitGroup{}

	group.Add(r.ConcurrentDownloader)
	for i := 0; i < r.ConcurrentDownloader; i++ {
		go func() {
			defer group.Done()
			for found := range foundCh {
				r.MDC.Debugf("Runner > downloading %s", found)
				dl, err := r.Downloader(found)
				if err != nil {
					r.report.AppendError(err)
					return
				}
				downloadedCh <- dl
				r.progressEventChannel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventDownloaded, Count: 1, Size: found.SimpleSignature().Size}
			}
		}()
	}

	go func() {
		group.Wait()
		close(downloadedCh)
	}()
}

func (r *Runner) pipeReadyToAnalyseToReadyToBackupChannels(downloadedCh chan backupmodel.FoundMedia, readyToAnalyse chan *backupmodel.AnalysedMedia) {
	group := sync.WaitGroup{}

	group.Add(r.ConcurrentAnalyser)
	for i := 0; i < r.ConcurrentAnalyser; i++ {
		go func() {
			defer group.Done()
			for media := range downloadedCh {
				r.MDC.Debugf("Runner > analysing %s", media)
				analysed, err := r.Analyser(media)
				if err != nil {
					r.report.AppendError(err)
					return
				}
				readyToAnalyse <- analysed
				r.progressEventChannel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventAnalysed, Count: 1, Size: media.SimpleSignature().Size}
			}
		}()
	}

	go func() {
		group.Wait()
		close(readyToAnalyse)
	}()
}

func (r *Runner) pipeReadyToBackupToCompletedChannels(readyToBackupChannel chan *backupmodel.AnalysedMedia, completionChannel chan *Report, progressChannel chan *backupmodel.ProgressEvent) {
	group := sync.WaitGroup{}
	group.Add(r.ConcurrentUploader)

	for i := 0; i < r.ConcurrentUploader; i++ {
		go func() {
			defer group.Done()
			buffer := make([]*backupmodel.AnalysedMedia, 0, r.UploadBatchSize)

			for media := range readyToBackupChannel {
				buffer = append(buffer, media)
				if len(buffer) == cap(buffer) {
					err := r.Uploader(buffer, r.progressEventChannel)
					if err != nil {
						r.report.AppendError(err)
						return
					}
					buffer = buffer[:0]
				}
			}

			// flush buffer
			if len(buffer) > 0 {
				err := r.Uploader(buffer, r.progressEventChannel)
				if err != nil {
					r.report.AppendError(err)
					return
				}
			}
		}()
	}

	go func() {
		group.Wait()

		if len(r.report.Errors) == 0 {
			r.report.AppendError(r.PreCompletion())
		}

		completionChannel <- r.report
		close(completionChannel)
		close(progressChannel)
	}()
}

func (r *Runner) startScanning(channel chan backupmodel.FoundMedia) {
	bufferChannel := make(chan backupmodel.FoundMedia, 255)
	var allMedias []backupmodel.FoundMedia
	var size uint = 0

	go func() {
		defer close(channel)
		for m := range bufferChannel {
			if r.Filter(m) {
				allMedias = append(allMedias, m)
				size += m.SimpleSignature().Size
			} else {
				r.progressEventChannel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventSkipped, Count: 1, Size: m.SimpleSignature().Size}
			}
		}

		r.progressEventChannel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventScanComplete, Count: uint(len(allMedias)), Size: size}

		for _, m := range allMedias {
			channel <- m
		}
	}()

	go func() {
		defer close(bufferChannel)
		_, _, err := r.Source(bufferChannel)
		r.report.AppendError(err)
	}()
}

// DummyProgressListener is logging each events from progress event channel ; it is mandatory to consume it
func DummyProgressListener(progressChannel chan *backupmodel.ProgressEvent) {
	go func() {
		for event := range progressChannel {
			log.Debugf("progress: %+v", event)
		}
	}()
}
