package backup

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

type RunnerAnalyser interface {
	Analyse(found FoundMedia, progressChannel chan *ProgressEvent) (*AnalysedMedia, error)
}

type RunnerCataloger interface {
	Catalog(ctx context.Context, medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error)
}

type RunnerUploader interface {
	Upload(buffer []*BackingUpMediaRequest, progressChannel chan *ProgressEvent) error
}

type runnerPublisher func(chan FoundMedia, chan *ProgressEvent) error
type RunnerAnalyserFunc func(found FoundMedia, progressChannel chan *ProgressEvent) (*AnalysedMedia, error)
type RunnerCatalogerFunc func(ctx context.Context, medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error)
type runnerUniqueFilter func(medias *BackingUpMediaRequest, progressChannel chan *ProgressEvent) bool
type RunnerUploaderFunc func(buffer []*BackingUpMediaRequest, progressChannel chan *ProgressEvent) error

func (r RunnerCatalogerFunc) Catalog(ctx context.Context, medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error) {
	return r(ctx, medias, progressChannel)
}

func (r RunnerUploaderFunc) Upload(buffer []*BackingUpMediaRequest, progressChannel chan *ProgressEvent) error {
	return r(buffer, progressChannel)
}

type runner struct {
	MDC                  *log.Entry         // MDC is log.WithFields({}) that contains Mapped Diagnostic Context
	Publisher            runnerPublisher    // Publisher is pushing files that have been found in the Volume into a channel
	Analyser             RunnerAnalyser     // Analyser is extracting metadata from the file
	Cataloger            RunnerCataloger    // Cataloger is assigning the media to an album and filtering out media already backed up
	UniqueFilter         runnerUniqueFilter // UniqueFilter is removing duplicates from the source Volume
	Uploader             RunnerUploader     // Uploader is storing the media in the archive, and registering it in the catalog
	ConcurrentAnalyser   int                // ConcurrentAnalyser is the number of goroutines that analyse the medias
	ConcurrentCataloguer int                // ConcurrentAnalyser is the number of goroutines that analyse the medias
	ConcurrentUploader   int                // ConcurrentUploader is the number of goroutine that upload files online
	BatchSize            int                // BatchSize is the size of the buffer for the uploader
	SkipRejects          bool
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

	log.WithError(err).Errorln("logging error and cancelling process")
	r.errorsMutex.Lock()
	defer r.errorsMutex.Unlock()

	r.cancel()

	r.errors = append(r.errors, err)
}

// start initialises the channels and publish files in the source channel
func (r *runner) start(ctx context.Context, sizeHint int) (chan *ProgressEvent, chan []error) {
	r.BatchSize = defaultValue(r.BatchSize, 1)

	r.progressEventChannel = make(chan *ProgressEvent, sizeHint*5)
	var cancellableCtx context.Context
	cancellableCtx, r.cancel = context.WithCancel(ctx)

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

	r.analyseMedias(cancellableCtx, foundChannel, analysedChannel)
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
		return errors.Wrapf(reportedErrors[0], "Backup failed, %d error(s) reported before shutdown. First one encountered", len(reportedErrors))
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

func (r *runner) analyseMedias(ctx context.Context, foundChannel chan FoundMedia, analysedChannel chan *AnalysedMedia) {
	r.startsInParallel(defaultValue(r.ConcurrentAnalyser, 1), func() {
		for {
			select {
			case <-ctx.Done():
				return
			case media, more := <-foundChannel:
				if more {
					r.MDC.Debugf("Runner > analysing %s", media)
					analysed, err := r.Analyser.Analyse(media, r.progressEventChannel)
					if err != nil && r.SkipRejects {
						r.MDC.Infof("silently skip %s: %s", media, err.Error()) // not strictly correct as the file won't be counted
					} else if err != nil {
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
	r.startsInParallel(defaultValue(r.ConcurrentCataloguer, 1), func() {
		for {
			select {
			case buffer, more := <-analysedChannel:
				if more && r.hasAnyErrors() == 0 {
					catalogedMedias, err := r.Cataloger.Catalog(context.TODO(), buffer, r.progressEventChannel)
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
	r.startsInParallel(defaultValue(r.ConcurrentUploader, 1), func() {
		for {
			select {
			case buffer, more := <-bufferedChannel:
				if more && r.hasAnyErrors() == 0 {
					err := r.Uploader.Upload(buffer, r.progressEventChannel)
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

func (r RunnerAnalyserFunc) Analyse(found FoundMedia, progressChannel chan *ProgressEvent) (*AnalysedMedia, error) {
	return r(found, progressChannel)
}
