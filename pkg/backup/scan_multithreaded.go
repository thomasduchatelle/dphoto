package backup

import (
	"context"
	"sync"
)

type scanCompleteObserver interface {
	OnScanComplete(ctx context.Context, count, size int) error
}

type analyserLauncher interface {
	process(ctx context.Context, volume SourceVolume) chan error
}

func newMultiThreadedController(concurrencyParameters ConcurrencyParameters, monitoringIntegrator scanMonitoringIntegrator, size int, flusher *flushableCollector) *multiThreadedController {
	return &multiThreadedController{
		scanMonitoringIntegrator: monitoringIntegrator,
		analysedMedias:           make(chan *AnalysedMedia, 255),
		bufferedAnalysedMedias:   make(chan []*AnalysedMedia, 255),
		backingUpMediaRequests:   make(chan []BackingUpMediaRequest, 255),
		concurrencyParameters:    concurrencyParameters,
		batchSize:                size,
		flusher:                  flusher,
	}
}

func (m *multiThreadedController) Launcher(analyser Analyser, chain *analyserObserverChain, tracker scanCompleteObserver) scanningLauncher {
	return &multiThreadedControllerLauncher{
		analyser:              analyser,
		analyserObserverChain: chain,
		scanCompleteObserver:  tracker,

		parent:         m,
		errorCollector: newErrorCollector(),
		mediaChannels:  make(chan FoundMedia, 255),
		done:           make(chan error, 1),
	}
}

// multiThreadedController is leveraging GO channels to distribute the load over several threads.
//
// 0.  ANALYSER [Nx analysers]
// 0.a  Analyser
// 0.b  analyserNoDateTimeFilter
// 1.  CATALOGUER [Mx cataloguers]
// 1.a  bufferAnalysedMedia
// 1.b  analyserToCatalogReferencer
// 2.  REDUCER [1x reducer]
// 2.a  applyFiltersOnCataloguer (postCatalogFiltersList)
// 2.b  uniqueFilter
type multiThreadedController struct {
	scanMonitoringIntegrator
	analysedMedias             chan *AnalysedMedia
	bufferedAnalysedMedias     chan []*AnalysedMedia
	backingUpMediaRequests     chan []BackingUpMediaRequest
	catalogReferencerObservers []CatalogReferencerObserver
	concurrencyParameters      ConcurrencyParameters
	cataloguerAdapter          cataloguerAdapter

	flusher   *flushableCollector
	batchSize int
}

type cataloguerAdapter interface {
	OnBatchOfAnalysedMedia(ctx context.Context, batch []*AnalysedMedia) error
}

type multiThreadedControllerLauncher struct {
	parent                *multiThreadedController
	analyser              Analyser
	analyserObserverChain *analyserObserverChain
	scanCompleteObserver  scanCompleteObserver
	mediaChannels         chan FoundMedia
	errorCollector        *errorCollector // TODO privatise errorCollector and remove the "observer" aspect.
	done                  chan error
}

func (m *multiThreadedController) AppendPreCataloguerFilter(observers ...CatalogReferencerObserver) []CatalogReferencerObserver {
	m.catalogReferencerObservers = append(m.catalogReferencerObservers, m.scanMonitoringIntegrator.AppendPreCataloguerFilter(observers...)...)
	return []CatalogReferencerObserver{CatalogReferencerObserverFunc(func(ctx context.Context, requests []BackingUpMediaRequest) error {
		m.backingUpMediaRequests <- requests
		return nil
	})}
}

func (m *multiThreadedController) bufferAnalysedMedia(ctx context.Context, adapter cataloguerAdapter) AnalysedMediaObserver {
	m.cataloguerAdapter = adapter
	return AnalysedMediaObserverFunc(func(ctx context.Context, media *AnalysedMedia) error {
		m.analysedMedias <- media
		return nil
	})
}

func (l *multiThreadedControllerLauncher) process(ctxWithoutCancel context.Context, volume SourceVolume) chan error {
	ctx, cancelFunc := context.WithCancel(ctxWithoutCancel)
	l.errorCollector.registerErrorObserver(func(err error) {
		cancelFunc()
	})

	go l.forwardsBackingUpMediaRequestsToTheChain(ctx)

	startsInParallel(ctx, l.parent.concurrencyParameters.NumberOfConcurrentCataloguerRoutines(), l.forwardsBufferedAnalysedMedia, func() {
		close(l.parent.backingUpMediaRequests)
	})

	go l.buffersAnalysedMedias(ctx)

	startsInParallel(ctx, l.parent.concurrencyParameters.NumberOfConcurrentAnalyserRoutines(), l.forwardsMediaChannelToTheChain, func() {
		close(l.parent.analysedMedias)
	})

	go l.readVolumeAndPublishFoundMediasInChannel(ctx, volume)

	return l.done
}

func (l *multiThreadedControllerLauncher) buffersAnalysedMedias(ctx context.Context) {
	defer close(l.parent.bufferedAnalysedMedias)

	mediaBuffer := buffer[*AnalysedMedia]{
		consumer: func(ctx context.Context, buffer []*AnalysedMedia) error {
			l.parent.bufferedAnalysedMedias <- buffer
			return nil
		},
		content: make([]*AnalysedMedia, 0, l.parent.batchSize),
	}

	for {
		select {
		case media, more := <-l.parent.analysedMedias:
			if more {
				_ = mediaBuffer.Append(ctx, media)
			} else {
				_ = mediaBuffer.Flush(ctx)
				return
			}
		}
	}
}

func (l *multiThreadedControllerLauncher) forwardsBufferedAnalysedMedia(ctx context.Context) {
	for requests := range l.parent.bufferedAnalysedMedias {
		err := l.parent.cataloguerAdapter.OnBatchOfAnalysedMedia(ctx, requests)
		if err != nil {
			l.errorCollector.appendError(err)
		}
	}
}

func (l *multiThreadedControllerLauncher) forwardsBackingUpMediaRequestsToTheChain(ctx context.Context) {
	defer close(l.done)

	for requestBatch := range l.parent.backingUpMediaRequests {
		for _, observer := range l.parent.catalogReferencerObservers {
			err := observer.OnMediaCatalogued(ctx, requestBatch)
			if err != nil {
				l.errorCollector.appendError(err)
			}
		}
	}

	if err := l.errorCollector.collectError(); err != nil {
		l.done <- err
	}
}

func (l *multiThreadedControllerLauncher) forwardsMediaChannelToTheChain(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case media, more := <-l.mediaChannels:
			if !more {
				return
			}

			err := l.analyser.Analyse(ctx, media, l.analyserObserverChain, l.analyserObserverChain)
			if err != nil {
				l.errorCollector.appendError(err)
			}
		}
	}
}

func (l *multiThreadedControllerLauncher) readVolumeAndPublishFoundMediasInChannel(ctx context.Context, volume SourceVolume) {
	defer close(l.mediaChannels)

	medias, err := volume.FindMedias()
	if err != nil {
		l.errorCollector.appendError(err)
		return
	}
	err = l.scanCompleteObserver.OnScanComplete(ctx, len(medias), sizeOfAllMedias(medias))
	if err != nil {
		l.errorCollector.appendError(err)
		return
	}

	for _, media := range medias {
		l.mediaChannels <- media
	}
}

func sizeOfAllMedias(medias []FoundMedia) int {
	size := 0
	for _, media := range medias {
		size += media.Size()
	}
	return size
}

func startsInParallel(ctx context.Context, parallel int, consume func(ctx context.Context), closeChannel func()) {
	group := sync.WaitGroup{}

	group.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer group.Done()

			consume(ctx)
		}()
	}

	go func() {
		group.Wait()
		closeChannel()
	}()
}
