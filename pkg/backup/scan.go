package backup

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type BatchScanner struct {
}

func (s *BatchScanner) Scan(ctx context.Context, owner ownermodel.Owner, volume SourceVolume, optionSlice ...Options) ([]*ScannedFolder, error) {

	options := ReduceOptions(optionSlice...)

	run, monitor, err := newMultithreadedScanRunner(ctx, options, volume.String(), owner)
	if err != nil {
		return nil, err
	}

	err, _ = <-run.process(ctx, volume)

	return monitor.getReport(ctx, err)
}

func newMultithreadedScanRunner(ctx context.Context, options Options, volumeName string, owner ownermodel.Owner) (analyserLauncher, *reportGetter, error) {
	tracker := newTrackerV2(options)
	report := newScanReport()
	scanLogger := newLogger(volumeName)

	flusher := new(flushableCollector)
	flusher.append(tracker)

	monitoring := &scanMonitoring{
		PostAnalyserSuccess:       []AnalysedMediaObserver{scanLogger},
		PostAnalyserRejects:       []RejectedMediaObserver{scanLogger, tracker},
		PostAnalyserFilterRejects: []RejectedMediaObserver{scanLogger, tracker, report},
		PreCataloguerFilter:       []CatalogReferencerObserver{scanLogger},
		PostCatalogFiltersIn:      []CatalogReferencerObserver{tracker, report},
		PostCatalogFiltersOut:     []CataloguerFilterObserver{scanLogger, tracker},
	}
	if options.SkipRejects {
		monitoring.PostAnalyserRejects = append(monitoring.PostAnalyserRejects, report)
	}

	chain, err := newAnalyserChain(ctx, options, monitoring, flusher, owner)
	return &multiThreadedAnalyserLauncher{
			analyser:              options.GetAnalyserDecorator().Decorate(getDefaultAnalyser()),
			analyserObserverChain: chain,
			scanCompleteObserver:  tracker,
		}, &reportGetter{
			flusher: flusher,
			report:  report,
		}, err
}

type reportGetter struct {
	flusher *flushableCollector
	report  *scanReport
}

func (g *reportGetter) getReport(ctx context.Context, previousErr error) ([]*ScannedFolder, error) {
	flushErr := g.flusher.flush(ctx)

	if previousErr != nil && flushErr != nil {
		log.WithError(flushErr).Error("Failed to flush (silenced to not override cause error)")
		return nil, previousErr

	} else if flushErr != nil {
		return nil, flushErr
	}

	return g.report.collect(), previousErr
}

// newAnalyserChain creates the chain of handlers used to analyse and catalog the medias
func newAnalyserChain(ctx context.Context, options Options, monitoring scanMonitoringIntegrator, flusher *flushableCollector, owner ownermodel.Owner) (*analyserObserverChain, error) {
	batchSize := defaultValue(options.BatchSize, 1)

	referencer, err := referencerFactory.NewDryRunReferencer(ctx, owner)
	if err != nil {
		return nil, err
	}

	var postAnalyserRejects []RejectedMediaObserver
	if !options.SkipRejects {
		postAnalyserRejects = append(postAnalyserRejects, new(analyserFailsFastObserver))
	}

	// TODO Multi-threading
	// 0.  ANALYSER [Nx analysers]
	// 0.a  Analyser
	// 0.b  analyserNoDateTimeFilter
	// 1.  CATALOGUER [Mx cataloguers]
	// 1.a  bufferAnalysedMedia
	// 1.b  analyserToCatalogReferencer
	// 2.  REDUCER [1x reducer]
	// 2.a  applyFiltersOnCataloguer (postCatalogFiltersList)
	// 2.b  uniqueFilter

	return &analyserObserverChain{
		AnalysedMediaObservers: monitoring.AppendPostAnalyserSuccess(&analyserNoDateTimeFilter{
			analyserObserverChain{
				AnalysedMediaObservers: []AnalysedMediaObserver{
					bufferAnalysedMedia(ctx, batchSize, &analyserToCatalogReferencer{
						CatalogReferencer: referencer,
						CatalogReferencerObservers: monitoring.AppendPreCataloguerFilter(&applyFiltersOnCataloguer{
							CatalogReferencerObservers: monitoring.AppendPostCatalogFiltersIn(),
							CataloguerFilterObservers:  monitoring.AppendPostCatalogFiltersOut(),
							CataloguerFilters:          postCatalogFiltersList(options),
						}),
					}, flusher.append),
				},
				RejectedMediaObservers: monitoring.AppendPostAnalyserFilterRejects(),
			},
		}),
		RejectedMediaObservers: monitoring.AppendPostAnalyserRejects(postAnalyserRejects...),
	}, nil
}

func postCatalogFiltersList(options Options) []CataloguerFilter {
	filters := []CataloguerFilter{
		mustNotExists(),
		mustBeUniqueInVolume(),
	}

	if len(options.RestrictedAlbumFolderName) > 0 {
		var albumFolderNames []string
		for albumFolderName := range options.RestrictedAlbumFolderName {
			albumFolderNames = append(albumFolderNames, albumFolderName)
		}
		filters = append(filters, mustBeInAlbum(albumFolderNames...))
	}

	return filters
}

type scanCompleteObserver interface {
	OnScanComplete(ctx context.Context, count, size int) error
}

type multiThreadedAnalyserLauncher struct {
	analyser              Analyser
	analyserObserverChain *analyserObserverChain
	scanCompleteObserver  scanCompleteObserver
}

func (l *multiThreadedAnalyserLauncher) process(ctx context.Context, volume SourceVolume) chan error {
	channel := make(chan error, 1)
	go func() {
		defer close(channel)

		medias, err := volume.FindMedias()
		if err != nil {
			channel <- err
			return
		}
		err = l.scanCompleteObserver.OnScanComplete(ctx, len(medias), sizeOfAllMedias(medias))

		for _, media := range medias {
			if err = l.analyser.Analyse(ctx, media, l.analyserObserverChain, l.analyserObserverChain); err != nil { // TODO is that how it should be interrupted? It could work if there if there is no queue...
				break
			}
		}

		channel <- err
	}()

	return channel
}

func sizeOfAllMedias(medias []FoundMedia) int {
	size := 0
	for _, media := range medias {
		size += media.Size()
	}
	return size
}

type analyserToCatalogReferencer struct {
	CatalogReferencer          CatalogReferencer
	CatalogReferencerObservers []CatalogReferencerObserver
}

func (s *analyserToCatalogReferencer) callback(ctx context.Context, buffer []*AnalysedMedia) error {
	return s.CatalogReferencer.Reference(ctx, buffer, s)
}

func (s *analyserToCatalogReferencer) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	for _, observer := range s.CatalogReferencerObservers {
		if err := observer.OnMediaCatalogued(ctx, requests); err != nil {
			return err
		}
	}

	return nil
}

type analyserLauncher interface {
	process(ctx context.Context, volume SourceVolume) chan error
}
