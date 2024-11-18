package backup

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type BatchScanner struct {
	CataloguerFactory CataloguerFactory
	DetailsReaders    DetailsReaderAdapter
}

func (s *BatchScanner) Scan(ctx context.Context, owner ownermodel.Owner, volume SourceVolume, optionSlice ...Options) ([]*ScannedFolder, error) {

	options := ReduceOptions(optionSlice...)

	launcher, monitor, err := s.newScanner(ctx, options, volume.String(), owner)
	if err != nil {
		return nil, err
	}

	err, _ = <-launcher.process(ctx, volume)

	return monitor.getReport(ctx, err)
}

func (s *BatchScanner) newScanner(ctx context.Context, options Options, volumeName string, owner ownermodel.Owner) (analyserLauncher, *reportGetter, error) {
	tracker := newTrackerV2(options)
	report := newScanReport()
	scanLogger := newLogger(volumeName)

	// TODO The flushableCollector won't be required once the controller manage the buffering (and flushing of the buffer)
	flusher := new(flushableCollector)
	flusher.append(tracker)

	monitoring := &scanListeners{
		scanCompleteObserver:      tracker,
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

	controller := newMultiThreadedController(options.ConcurrencyParameters, monitoring, options.GetBatchSize(), flusher)

	cataloguer, err := s.CataloguerFactory.NewDryRunCataloguer(ctx, owner)
	if err != nil {
		return nil, nil, err
	}

	launcher, err := newScanningChain(ctx, controller, scanningOptions{
		Options:    options,
		cataloguer: cataloguer,
		analyser:   options.GetAnalyserDecorator().Decorate(newDefaultAnalyser(s.DetailsReaders)),
	})
	return launcher, &reportGetter{
		flusher: flusher,
		report:  report,
	}, err
}

type reportGetter struct {
	flusher *flushableCollector
	report  *scanReport
}

func (g *reportGetter) getReport(ctx context.Context, previousErr error) ([]*ScannedFolder, error) {
	flushErr := g.flusher.Flush(ctx)

	if previousErr != nil && flushErr != nil {
		log.WithError(flushErr).Error("Failed to flush (silenced to not override cause error)")
		return nil, previousErr

	} else if flushErr != nil {
		return nil, flushErr
	}

	return g.report.collect(), previousErr
}
