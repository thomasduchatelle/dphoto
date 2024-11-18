package backup

import (
	"context"
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

	return monitor.report.collect(), err
}

func (s *BatchScanner) newScanner(ctx context.Context, options Options, volumeName string, owner ownermodel.Owner) (analyserLauncher, *reportGetter, error) {
	tracker := newTrackerV2(options)
	report := newScanReport()
	scanLogger := newLogger(volumeName)

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

	controller := newMultiThreadedController(options.ConcurrencyParameters, monitoring)
	controller.registerWrappers(tracker)

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
		report: report,
	}, err
}

type reportGetter struct {
	report *scanReport
}
