package scanner

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

// Source populates media channel with everything found on the volume.
type Source func(medias chan FoundMedia) (uint, uint, error)

// Filter returns true if file should be backed up
type Filter func(found FoundMedia) bool

// Analyser reads the header of the file to find metadata (EXIF, dimensions, ...)
type Analyser func(found FoundMedia) (*AnalysedMedia, error)

// Downloader downloads locally the file to avoid multi-reads and too high concurrency on slow media
type Downloader func(found FoundMedia) (FoundMedia, error)

// Uploader backups media on an online storage (and update the indexes)
type Uploader func(buffer []*AnalysedMedia, progressChannel chan *ProgressEvent) error

// PreCompletion is called just before the run complete.
type PreCompletion func() error

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
	progressEventChannel chan *ProgressEvent
}

func Start(runner Runner) (chan *Report, chan *ProgressEvent) {
	runner.report = &Report{}
	runner.completionChannel = make(chan *Report, 1)
	runner.progressEventChannel = make(chan *ProgressEvent, runner.BufferSize)

	foundChannel := make(chan FoundMedia, runner.FoundMediaBufferSize)
	readyToAnalyseChannel := make(chan FoundMedia, runner.BufferSize)
	readyToBackupChannel := make(chan *AnalysedMedia, runner.BufferSize)

	runner.pipeFoundToReadyToAnalyseChannels(foundChannel, readyToAnalyseChannel)
	runner.pipeReadyToAnalyseToReadyToBackupChannels(readyToAnalyseChannel, readyToBackupChannel)
	runner.pipeReadyToBackupToCompletedChannels(readyToBackupChannel, runner.completionChannel, runner.progressEventChannel)

	runner.startScanning(foundChannel)

	return runner.completionChannel, runner.progressEventChannel
}

func (r *Runner) pipeFoundToReadyToAnalyseChannels(foundCh, downloadedCh chan FoundMedia) {
	group := sync.WaitGroup{}

	group.Add(r.ConcurrentDownloader)
	for i := 0; i < r.ConcurrentDownloader; i++ {
		go func() {
			defer group.Done()
			for found := range foundCh {
				if r.Filter(found) {
					r.MDC.Debugf("Runner > downloading %s", found)
					dl, err := r.Downloader(found)
					if err != nil {
						r.report.AppendError(err)
						return
					}
					downloadedCh <- dl
					r.progressEventChannel <- &ProgressEvent{Type: ProgressEventDownloaded, Count: 1, Size: found.SimpleSignature().Size}

				} else {
					r.progressEventChannel <- &ProgressEvent{Type: ProgressEventSkipped, Count: 1, Size: found.SimpleSignature().Size}
				}
			}
		}()
	}

	go func() {
		group.Wait()
		close(downloadedCh)
	}()
}

func (r *Runner) pipeReadyToAnalyseToReadyToBackupChannels(downloadedCh chan FoundMedia, readyToAnalyse chan *AnalysedMedia) {
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
				r.progressEventChannel <- &ProgressEvent{Type: ProgressEventAnalysed, Count: 1, Size: media.SimpleSignature().Size}
			}
		}()
	}

	go func() {
		group.Wait()
		close(readyToAnalyse)
	}()
}

func (r *Runner) pipeReadyToBackupToCompletedChannels(readyToBackupChannel chan *AnalysedMedia, completionChannel chan *Report, progressChannel chan *ProgressEvent) {
	group := sync.WaitGroup{}
	group.Add(r.ConcurrentUploader)

	for i := 0; i < r.ConcurrentUploader; i++ {
		go func() {
			defer group.Done()
			buffer := make([]*AnalysedMedia, 0, r.UploadBatchSize)

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

func (r *Runner) startScanning(channel chan FoundMedia) {
	go func() {
		defer close(channel)
		count, size, err := r.Source(channel)
		r.report.AppendError(err)
		r.progressEventChannel <- &ProgressEvent{Type: ProgressEventScanComplete, Count: count, Size: size}
	}()
}

// DummyProgressListener is logging each events from progress event channel ; it is mandatory to consume it
func DummyProgressListener(progressChannel chan *ProgressEvent) {
	go func() {
		for event := range progressChannel {
			log.Debugf("progress: %+v", event)
		}
	}()
}
