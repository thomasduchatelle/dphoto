package backup

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

type runnerPublisher func(chan FoundMedia, chan *ProgressEvent) error
type runnerAnalyser func(found FoundMedia, progressChannel chan *ProgressEvent) (*AnalysedMedia, error)
type runnerCataloger func(medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error)
type runnerUniqueFilter func(medias *BackingUpMediaRequest, progressChannel chan *ProgressEvent) bool
type runnerUploader func(buffer []*BackingUpMediaRequest, progressChannel chan *ProgressEvent) error

type runner struct {
	MDC                  *log.Entry // MDC is log.WithFields({}) that contains Mapped Diagnostic Context
	Publisher            runnerPublisher
	Analyser             runnerAnalyser
	Cataloger            runnerCataloger
	UniqueFilter         runnerUniqueFilter
	Uploader             runnerUploader
	ConcurrentAnalyser   int // ConcurrentAnalyser is the number of goroutines that analyse the medias
	ConcurrentCataloguer int // ConcurrentAnalyser is the number of goroutines that analyse the medias
	ConcurrentUploader   int // ConcurrentUploader is the number of goroutine that upload files online
	BatchSize            int // BatchSize is the size of the buffer for the uploader
	context              context.Context
	cancel               context.CancelFunc
	progressEventChannel chan *ProgressEvent
	errorsMutex          sync.Mutex
	errors               []error
	completionChannel    chan []error
}

func (r *runner) appendError(err error) {
	if err == nil {
		return
	}

	log.WithError(err).Errorln("cancelling queued processes")
	r.errorsMutex.Lock()
	defer r.errorsMutex.Unlock()

	r.cancel()

	r.errors = append(r.errors, err)
}

// start initialises the channels and publish files in the source channel
func (r *runner) start(ctx context.Context, sizeHint int) (chan *ProgressEvent, chan []error) {
	r.progressEventChannel = make(chan *ProgressEvent, sizeHint*5)
	r.context, r.cancel = context.WithCancel(ctx)

	bufferedChannelSize := 1 + sizeHint/r.BatchSize
	if sizeHint == 0 {
		bufferedChannelSize = 0
	}

	foundChannel := make(chan FoundMedia, sizeHint)
	analysedChannel := make(chan *AnalysedMedia, sizeHint)
	bufferedAnalysedChannel := make(chan []*AnalysedMedia, bufferedChannelSize)
	cataloguedChannel := make(chan *BackingUpMediaRequest, bufferedChannelSize)
	bufferedCataloguedChannel := make(chan []*BackingUpMediaRequest, bufferedChannelSize)
	r.completionChannel = make(chan []error, 1)

	r.analyseMedias(foundChannel, analysedChannel)
	r.bufferAnalysedMedias(analysedChannel, bufferedAnalysedChannel)
	r.catalogueMedias(bufferedAnalysedChannel, cataloguedChannel)
	r.bufferUniqueCataloguedMedias(cataloguedChannel, bufferedCataloguedChannel)
	r.uploadMedias(bufferedCataloguedChannel, r.completionChannel)

	r.startPublishing(foundChannel)
	return r.progressEventChannel, r.completionChannel
}

// waitToFinish is blocking until runner completes (or fails), returned completion channel should not be consumed.
func (r *runner) waitToFinish() error {

	reportedErrors := <-r.completionChannel

	for i, err := range reportedErrors {
		r.MDC.WithError(err).Errorf("Error %d/%d: %s", i+1, len(reportedErrors), err.Error())
	}

	if len(reportedErrors) > 0 {
		return errors.Wrapf(reportedErrors[0], "Backup failed, %d errors reported until shutdown.", len(reportedErrors))
	}

	return nil
}

func (r *runner) startsInParallel(parallel int, consume func(), closeChannel func()) {
	group := sync.WaitGroup{}

	group.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer group.Done()

			consume()
		}()
	}

	go func() {
		group.Wait()
		closeChannel()
	}()
}
func (r *runner) analyseMedias(foundChannel chan FoundMedia, analysedChannel chan *AnalysedMedia) {
	r.startsInParallel(r.ConcurrentAnalyser, func() {
		for {
			select {
			case <-r.context.Done():
				return
			case media, more := <-foundChannel:
				if more {
					r.MDC.Debugf("Runner > analysing %s", media)
					analysed, err := r.Analyser(media, r.progressEventChannel)
					if err != nil {
						r.appendError(errors.Wrap(err, "error in analyser"))
					} else {
						analysedChannel <- analysed
					}

				} else {
					return
				}
			}
		}
	}, func() {
		close(analysedChannel)
	})
}

func (r *runner) bufferAnalysedMedias(readyToBackupChannel chan *AnalysedMedia, bufferedChannel chan []*AnalysedMedia) {
	go func() {
		defer close(bufferedChannel)

		buffer := make([]*AnalysedMedia, 0, r.BatchSize)
		for media := range readyToBackupChannel {
			buffer = append(buffer, media)
			if len(buffer) == cap(buffer) {
				bufferedChannel <- buffer
				buffer = make([]*AnalysedMedia, 0, r.BatchSize)
			}
		}

		if len(buffer) > 0 {
			bufferedChannel <- buffer
		}
	}()
}

func (r *runner) catalogueMedias(analysedChannel chan []*AnalysedMedia, requestsChannel chan *BackingUpMediaRequest) {
	r.startsInParallel(r.ConcurrentCataloguer, func() {
		for {
			select {
			case <-r.context.Done():
				return

			case buffer, more := <-analysedChannel:
				if more && r.hasAnyErrors() == 0 {
					catalogedMedias, err := r.Cataloger(buffer, r.progressEventChannel)
					r.appendError(errors.Wrap(err, "error in cataloguer"))

					for _, media := range catalogedMedias {
						requestsChannel <- media
					}
				} else if !more {
					return
				}
			}
		}
	}, func() {
		close(requestsChannel)
	})
}

func (r *runner) hasAnyErrors() int {
	r.errorsMutex.Lock()
	defer r.errorsMutex.Unlock()

	return len(r.errors)
}

func (r *runner) bufferUniqueCataloguedMedias(backingUpMediaChannel chan *BackingUpMediaRequest, bufferedChannel chan []*BackingUpMediaRequest) {
	go func() {
		defer close(bufferedChannel)

		buffer := make([]*BackingUpMediaRequest, 0, r.BatchSize)
		for analysed := range backingUpMediaChannel {
			if r.UniqueFilter(analysed, r.progressEventChannel) {

				buffer = append(buffer, analysed)
				if len(buffer) >= cap(buffer) {
					bufferedChannel <- buffer
					buffer = make([]*BackingUpMediaRequest, 0, r.BatchSize)
				}
			}
		}

		if len(buffer) > 0 {
			bufferedChannel <- buffer
		}
	}()
}

func (r *runner) uploadMedias(bufferedChannel chan []*BackingUpMediaRequest, completionChannel chan []error) {
	r.startsInParallel(r.ConcurrentUploader, func() {
		for {
			select {
			case <-r.context.Done():
				return

			case buffer, more := <-bufferedChannel:
				if more && r.hasAnyErrors() == 0 {
					err := r.Uploader(buffer, r.progressEventChannel)
					r.appendError(errors.Wrap(err, "error in uploader"))
				} else if !more {
					return
				}
			}
		}
	}, func() {
		r.errorsMutex.Lock()
		defer r.errorsMutex.Unlock()

		errs := make([]error, len(r.errors), len(r.errors))
		copy(errs, r.errors)

		completionChannel <- errs
		close(completionChannel)
		close(r.progressEventChannel)
	})
}

func (r *runner) startPublishing(foundChannel chan FoundMedia) {
	go func() {
		err := r.Publisher(foundChannel, r.progressEventChannel)
		r.appendError(err)
		close(foundChannel)
	}()
}
