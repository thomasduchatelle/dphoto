package backup

import "C"
import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

type RunnerCataloger interface {
	Catalog(ctx context.Context, medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error)
}

type RunnerUploader interface {
	Upload(buffer []*BackingUpMediaRequest, progressChannel chan *ProgressEvent) error
}

type runnerPublisher func(chan FoundMedia, chan *ProgressEvent) error
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
	Options              Options            // Options contains the configuration for each step of the backup process. [migration to Observer Pattern]
	Publisher            runnerPublisher    // Publisher is pushing files that have been found in the Volume into a channel
	Analyser             Analyser           // Analyser is extracting metadata from the file
	Cataloger            RunnerCataloger    // Cataloger is assigning the media to an album and filtering out media already backed up
	UniqueFilter         runnerUniqueFilter // UniqueFilter is removing duplicates from the source Volume
	Uploader             RunnerUploader     // Uploader is storing the media in the archive, and registering it in the catalog
	ConcurrentCataloguer int                // ConcurrentAnalyser is the number of goroutines that analyse the medias [DEPRECATED - part of the 'runner/orchestrator' pattern, use Options to migrate to the 'Observer' pattern]
	ConcurrentUploader   int                // ConcurrentUploader is the number of goroutine that upload files online [DEPRECATED - part of the 'runner/orchestrator' pattern, use Options to migrate to the 'Observer' pattern]
	BatchSize            int                // BatchSize is the size of the buffer for the uploader	 [DEPRECATED - part of the 'runner/orchestrator' pattern, use Options to migrate to the 'Observer' pattern]
	SkipRejects          bool               // SkipRejects [DEPRECATED - part of the 'runner/orchestrator' pattern, use Options to migrate to the 'Observer' pattern]
	progressEventChannel chan *ProgressEvent
	errors               []error

	completionChannel      chan []error
	interrupterObserver    *InterrupterObserver    // InterrupterObserver is temporary while refactoring to observer pattern
	errorCollectorObserver *ErrorCollectorObserver // ErrorCollectorObserver is temporary while refactoring to observer pattern
}

func (r *runner) appendError(err error) {
	if err == nil {
		return
	}

	// Deprecated - use observers directly
	r.errorCollectorObserver.appendError(err)
	r.interrupterObserver.cancel()
}

// start initialises the channels and publish files in the source channel
func (r *runner) start(ctx context.Context, sizeHint int) (chan *ProgressEvent, chan []error) {
	progressObserver := NewProgressObserver(sizeHint)
	r.progressEventChannel = progressObserver.EventChannel

	r.BatchSize = defaultValue(r.BatchSize, 1)
	r.Analyser = r.Options.GetAnalyserDecorator().Decorate(r.Analyser, progressObserver)

	var interruptableContext context.Context
	r.interrupterObserver, interruptableContext = NewInterrupterObserver(ctx)
	r.errorCollectorObserver = NewErrorCollectorObserver()

	bufferedChannelSize := 1 + sizeHint/r.BatchSize
	if sizeHint == 0 {
		bufferedChannelSize = 0
	}

	foundChannel := make(chan FoundMedia, sizeHint)
	analysedMediaChannelPublisher := NewAsyncPublisher(sizeHint)
	bufferedAnalysedChannel := make(chan []*AnalysedMedia, bufferedChannelSize)
	cataloguedChannel := make(chan *BackingUpMediaRequest, bufferedChannelSize)
	bufferedCataloguedChannel := make(chan []*BackingUpMediaRequest, bufferedChannelSize)
	r.completionChannel = make(chan []error, 1)

	observer := &CompositeRunnerObserver{
		Observers: []interface{}{
			r.errorCollectorObserver,
			r.interrupterObserver,
			progressObserver,
			analysedMediaChannelPublisher,
		},
	}

	r.analyseMedias(interruptableContext, foundChannel, observer, analysedMediaChannelPublisher.Close)
	r.bufferAnalysedMedias(analysedMediaChannelPublisher.AnalysedMediaChannel, bufferedAnalysedChannel)
	r.catalogueMedias(interruptableContext, bufferedAnalysedChannel, cataloguedChannel)
	r.bufferUniqueCataloguedMedias(cataloguedChannel, bufferedCataloguedChannel)
	r.uploadMedias(interruptableContext, bufferedCataloguedChannel, r.completionChannel, r.errorCollectorObserver)

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

func (r *runner) analyseMedias(ctx context.Context, foundChannel chan FoundMedia, observer *CompositeRunnerObserver, closer func()) {
	r.startsInParallel(defaultValue(r.Options.ConcurrentAnalyser, 1), func() {
		for {
			select {
			case <-ctx.Done():
				return
			case media, more := <-foundChannel:
				if more {
					r.Analyser.Analyse(media, observer, observer) // TODO manage the SkipRejects mode

				} else {
					return
				}
			}
		}
	}, closer)
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

func (r *runner) catalogueMedias(ctx context.Context, analysedChannel chan []*AnalysedMedia, requestsChannel chan *BackingUpMediaRequest) {
	r.startsInParallel(defaultValue(r.ConcurrentCataloguer, 1), func() {
		interrupted := false
		for {
			select {
			case <-ctx.Done():
				interrupted = true
			case buffer, more := <-analysedChannel:
				if !more {
					return
				}
				if !interrupted {
					catalogedMedias, err := r.Cataloger.Catalog(context.TODO(), buffer, r.progressEventChannel)
					r.appendError(errors.Wrap(err, "error in cataloguer"))

					for _, media := range catalogedMedias {
						requestsChannel <- media
					}
				}
			}
		}
	}, func() {
		close(requestsChannel)
	})
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

func (r *runner) uploadMedias(ctx context.Context, bufferedChannel chan []*BackingUpMediaRequest, completionChannel chan []error, collector *ErrorCollectorObserver) {
	r.startsInParallel(defaultValue(r.ConcurrentUploader, 1), func() {
		interrupted := false
		for {
			select {
			case <-ctx.Done():
				interrupted = true

			case buffer, more := <-bufferedChannel:
				if more && !interrupted {
					err := r.Uploader.Upload(buffer, r.progressEventChannel)
					r.appendError(errors.Wrap(err, "error in uploader"))
				} else if !more {
					return
				}
			}
		}
	}, func() {
		completionChannel <- collector.Errors()
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
