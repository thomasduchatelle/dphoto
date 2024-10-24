package backup

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
	progressEventChannel chan *ProgressEvent
	channelPublisher     *ChannelPublisher

	interrupterObserver    Interrupter             // AnalyserInterrupterObserver is temporary while refactoring to observer pattern
	errorCollectorObserver IErrorCollectorObserver // ErrorCollectorObserver is temporary while refactoring to observer pattern
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
func (r *runner) start(ctx context.Context, sizeHint int) (chan *ProgressEvent, chan []error) {
	progressObserver := NewProgressObserver(sizeHint)
	r.channelPublisher = NewAsyncPublisher(sizeHint, r.batchSize())

	r.Analyser = r.Options.GetAnalyserDecorator().Decorate(r.Analyser, progressObserver)

	var interruptableContext context.Context
	r.interrupterObserver, interruptableContext = NewInterrupterObserver(ctx, r.Options)
	r.errorCollectorObserver = NewErrorCollectorObserver()

	observer, err := r.CreateObserver(progressObserver)
	if err != nil {
		panic(err)
	}

	r.progressEventChannel = progressObserver.EventChannel
	r.analyseMedias(interruptableContext, r.channelPublisher.FoundChannel, observer, r.channelPublisher.AnalysedMediaChannelCloser)
	r.bufferAnalysedMedias(r.channelPublisher.AnalysedMediaChannel, r.channelPublisher.BufferedAnalysedChannel)
	r.catalogueMedias(interruptableContext, r.channelPublisher.BufferedAnalysedChannel, r.channelPublisher.CataloguedChannel)
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
	analyser := &AnalyserAsyncWrapper{
		Analyser: r.Analyser,
	}
	r.startsInParallel(defaultValue(r.Options.ConcurrentAnalyser, 1), func() {
		for {
			select {
			case <-ctx.Done():
				return
			case media, more := <-foundChannel:
				if more {
					analyser.Analyse(ctx, media, observer, observer)

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

func (r *runner) catalogueMedias(ctx context.Context, analysedChannel chan []*AnalysedMedia, requestsChannel chan *BackingUpMediaRequest) {
	r.startsInParallel(defaultValue(r.Options.ConcurrentCataloguer, 1), func() {
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
	r.startsInParallel(defaultValue(r.Options.ConcurrentUploader, 1), func() {
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
