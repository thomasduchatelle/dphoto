package runner

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	log "github.com/sirupsen/logrus"
	"sync"
)

// Source populates media channel with everything found on the volume.
type Source func(medias chan model.FoundMedia) error

// Filter returns true if file should be backed up
type Filter func(found model.FoundMedia) bool

// Analyser reads the header of the file to find metadata (EXIF, dimensions, ...)
type Analyser func(found model.FoundMedia) (*model.AnalysedMedia, error)

// Downloader downloads locally the file to avoid multi-reads and too high concurrency on slow media
type Downloader func(found model.FoundMedia) (model.FoundMedia, error)

// Uploader backups media on an online storage (and update the indexes)
type Uploader func(buffer []*model.AnalysedMedia) error

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
	BufferSize           int // BufferSize is the default size for channels
	ConcurrentDownloader int // ConcurrentDownloader is the number of goroutines that will filter and download files ; should be 1 with a pass-through downloaders
	ConcurrentAnalyser   int // ConcurrentAnalyser is the number of goroutines that analyse the medias
	ConcurrentUploader   int // ConcurrentUploader is the number of goroutine that upload files online
	UploadBatchSize      int // UploadBatchSize is the size of the buffer for the uploader
	report               *Report
	completionChannel    chan *Report
}

func Start(runner Runner) chan *Report {
	runner.report = &Report{}
	runner.completionChannel = make(chan *Report, 1)

	foundChannel := make(chan model.FoundMedia, runner.BufferSize)
	readyToAnalyseChannel := make(chan model.FoundMedia, runner.BufferSize)
	readyToBackupChannel := make(chan *model.AnalysedMedia, runner.BufferSize)

	runner.pipeFoundToReadyToAnalyseChannels(foundChannel, readyToAnalyseChannel)
	runner.pipeReadyToAnalyseToReadyToBackupChannels(readyToAnalyseChannel, readyToBackupChannel)
	runner.pipeReadyToBackupToCompletedChannels(readyToBackupChannel, runner.completionChannel)

	runner.startScanning(foundChannel)

	return runner.completionChannel
}

func (r *Runner) pipeFoundToReadyToAnalyseChannels(foundCh, downloadedCh chan model.FoundMedia) {
	group := sync.WaitGroup{}

	group.Add(r.ConcurrentDownloader)
	for i := 0; i < r.ConcurrentDownloader; i++ {
		go func() {
			defer group.Done()
			for found := range foundCh {
				if r.Filter(found) {
					dl, err := r.Downloader(found)
					if err != nil {
						r.report.AppendError(err)
						return
					}
					downloadedCh <- dl
				}
			}
		}()
	}

	go func() {
		group.Wait()
		close(downloadedCh)
	}()
}

func (r *Runner) pipeReadyToAnalyseToReadyToBackupChannels(downloadedCh chan model.FoundMedia, readyToAnalyse chan *model.AnalysedMedia) {
	group := sync.WaitGroup{}

	group.Add(r.ConcurrentAnalyser)
	for i := 0; i < r.ConcurrentAnalyser; i++ {
		go func() {
			defer group.Done()
			for media := range downloadedCh {
				analysed, err := r.Analyser(media)
				if err != nil {
					r.report.AppendError(err)
					return
				}
				readyToAnalyse <- analysed
			}
		}()
	}

	go func() {
		group.Wait()
		close(readyToAnalyse)
	}()
}

func (r *Runner) pipeReadyToBackupToCompletedChannels(readyToBackupChannel chan *model.AnalysedMedia, completionChannel chan *Report) {
	group := sync.WaitGroup{}
	group.Add(r.ConcurrentUploader)

	for i := 0; i < r.ConcurrentUploader; i++ {
		go func() {
			defer group.Done()
			buffer := make([]*model.AnalysedMedia, 0, r.UploadBatchSize)

			for media := range readyToBackupChannel {
				buffer = append(buffer, media)
				if len(buffer) == cap(buffer) {
					err := r.Uploader(buffer)
					if err != nil {
						r.report.AppendError(err)
						return
					}
				}
			}

			// flush buffer
			if len(buffer) > 0 {
				err := r.Uploader(buffer)
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
	}()
}

func (r *Runner) startScanning(channel chan model.FoundMedia) {
	go func() {
		defer close(channel)
		err := r.Source(channel)
		r.report.AppendError(err)
	}()
}
