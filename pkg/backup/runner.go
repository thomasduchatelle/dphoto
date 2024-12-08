package backup

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

type RunnerUploader interface {
	Upload(buffer []*BackingUpMediaRequest, progressChannel chan *progressEvent) error
}

type runnerPublisher func(chan FoundMedia, chan *progressEvent) error
type runnerUniqueFilter func(medias *BackingUpMediaRequest, progressChannel chan *progressEvent) bool
type RunnerUploaderFunc func(buffer []*BackingUpMediaRequest, progressChannel chan *progressEvent) error

func (r RunnerUploaderFunc) Upload(buffer []*BackingUpMediaRequest, progressChannel chan *progressEvent) error {
	return r(buffer, progressChannel)
}

type runner struct {
	MDC                  *log.Entry         // MDC is log.WithFields({}) that contains Mapped Diagnostic Context
	Options              Options            // Options contains the configuration for each step of the backup process. [migration to CataloguerFilterObserver Pattern]
	Publisher            runnerPublisher    // Publisher is pushing files that have been found in the Volume into a channel
	Analyser             Analyser           // Analyser is extracting metadata from the file
	CatalogReferencer    Cataloguer         // CatalogReferencer is assigning the media to an album and filtering out media already backed up
	UniqueFilter         runnerUniqueFilter // UniqueFilter is removing duplicates from the source Volume
	Uploader             RunnerUploader     // Uploader is storing the media in the archive, and registering it in the catalog
	Listeners            []interface{}      // Listeners is a list of observers that are notified of the progress of the backup process, they need to implement the appropriate interfaces
	progressEventChannel chan *progressEvent
	channelPublisher     *ChannelPublisher

	interrupterObserver    Interrupter             // AnalyserInterrupterObserver is temporary while refactoring to observer pattern
	errorCollectorObserver IErrorCollectorObserver // errorCollector is temporary while refactoring to observer pattern
}

func (r *runner) appendError(err error) {
	if err == nil {
		return
	}

	// Deprecated - use observers directly
	r.errorCollectorObserver.appendError(err)
	r.interrupterObserver.Cancel()
}

// start initialises the channels and publish files in the source channel
func (r *runner) start(ctx context.Context, sizeHint int) (chan *progressEvent, chan []error) {
	progressObserver := NewProgressObserver(sizeHint)
	r.channelPublisher = NewAsyncPublisher(sizeHint, r.batchSize())

	r.Analyser = r.Options.GetAnalyserDecorator().Decorate(r.Analyser, progressObserver)

	var interruptableContext context.Context
	r.interrupterObserver, interruptableContext = NewInterrupterObserver(ctx, r.Options)
	r.errorCollectorObserver = newErrorCollector()

	observer, err := r.CreateObserver(progressObserver)
	if err != nil {
		panic(err)
	}

	r.progressEventChannel = progressObserver.EventChannel
	r.analyseMedias(interruptableContext, r.channelPublisher.FoundChannel, observer, r.channelPublisher.AnalysedMediaChannelCloser)
	r.bufferAnalysedMedias(r.channelPublisher.AnalysedMediaChannel, r.channelPublisher.BufferedAnalysedChannel)
	r.catalogueMedias(interruptableContext, r.channelPublisher.BufferedAnalysedChannel, observer, r.channelPublisher.CataloguedChannelCloser)
	r.bufferUniqueCataloguedMedias(r.channelPublisher.CataloguedChannel, r.channelPublisher.BufferedCataloguedChannel)
	r.uploadMedias(interruptableContext, r.channelPublisher.BufferedCataloguedChannel, r.channelPublisher.CompletionChannel, r.errorCollectorObserver)

	r.startPublishing(r.channelPublisher.FoundChannel)
	return r.progressEventChannel, r.channelPublisher.CompletionChannel
}

func (r *runner) CreateObserver(progressObserver *ProgressObserver) (*CompositeRunnerObserver, error) {
	observer := &CompositeRunnerObserver{
		Observers: []interface{}{
			r.errorCollectorObserver,
			r.interrupterObserver,
			progressObserver,
		},
	}

	if r.Options.RejectDir != "" {
		rejectsObserver, err := NewCopyRejectsObserver(r.Options.RejectDir)
		if err != nil {
			return nil, err
		}
		observer.Observers = append(observer.Observers, rejectsObserver)
	}

	observer.Observers = append(observer.Observers, r.channelPublisher)

	if len(r.Listeners) > 0 {
		observer.Observers = append(observer.Observers, r.Listeners...)
	}

	return observer, nil
}

// waitToFinish is blocking until runner completes (or fails), returned completion channel should not be consumed.
func (r *runner) waitToFinish() error {

	reportedErrors := <-r.channelPublisher.CompletionChannel

	for i, err := range reportedErrors {
		r.MDC.WithError(err).Errorf("Error %d/%d: %s", i+1, len(reportedErrors), err.Error())
	}

	if !r.Options.SkipRejects && len(reportedErrors) > 0 {
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
	analyserObserver := newAnalyserObserverChain(r.Options, observer)

	r.startsInParallel(defaultValue(r.Options.ConcurrencyParameters.ConcurrentAnalyserRoutines, 1), func() {
		for {
			select {
			case <-ctx.Done():
				return
			case media, more := <-foundChannel:
				if more {
					err := r.Analyser.Analyse(ctx, media, analyserObserver, analyserObserver)
					r.appendError(errors.Wrap(err, "error in analyser"))

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

		buffer := make([]*AnalysedMedia, 0, r.batchSize())
		for media := range readyToBackupChannel {
			buffer = append(buffer, media)
			if len(buffer) == cap(buffer) {
				bufferedChannel <- buffer
				buffer = make([]*AnalysedMedia, 0, r.batchSize())
			}
		}

		if len(buffer) > 0 {
			bufferedChannel <- buffer
		}
	}()
}

func (r *runner) catalogueMedias(ctx context.Context, analysedChannel chan []*AnalysedMedia, observer *CompositeRunnerObserver, closer func()) {
	cataloguer := &CataloguerWithFilters{
		Delegate:                  r.CatalogReferencer,
		CataloguerFilters:         postCatalogFiltersList(r.Options),
		CatalogReferencerObserver: observer,
		CataloguerFilterObserver:  observer,
	}

	r.startsInParallel(1, func() { // FIXME it was 'defaultValue(r.Options.ConcurrencyParameters.ConcurrentCataloguerRoutines, 1)' but needed to be removed due to concurrency race.
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
					err := cataloguer.Catalog(ctx, buffer) // FIXME the 'observer' has been removed.
					r.appendError(errors.Wrap(err, "error in cataloguer"))
				}
			}
		}
	}, closer)
}

func (r *runner) bufferUniqueCataloguedMedias(backingUpMediaChannel chan *BackingUpMediaRequest, bufferedChannel chan []*BackingUpMediaRequest) {
	go func() {
		defer close(bufferedChannel)

		buffer := make([]*BackingUpMediaRequest, 0, r.batchSize())
		for analysed := range backingUpMediaChannel {
			if r.UniqueFilter(analysed, r.progressEventChannel) {

				buffer = append(buffer, analysed)
				if len(buffer) >= cap(buffer) {
					bufferedChannel <- buffer
					buffer = make([]*BackingUpMediaRequest, 0, r.batchSize())
				}
			}
		}

		if len(buffer) > 0 {
			bufferedChannel <- buffer
		}
	}()
}

func (r *runner) uploadMedias(ctx context.Context, bufferedChannel chan []*BackingUpMediaRequest, completionChannel chan []error, collector IErrorCollectorObserver) {
	r.startsInParallel(defaultValue(r.Options.ConcurrencyParameters.ConcurrentUploaderRoutines, 1), func() {
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

func (r *runner) batchSize() int {
	return defaultValue(r.Options.BatchSize, 1)
}

func (r *runner) startPublishing(foundChannel chan FoundMedia) {
	go func() {
		err := r.Publisher(foundChannel, r.progressEventChannel)
		r.appendError(err)
		close(foundChannel)
	}()
}
