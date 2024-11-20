package backup

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/backup/chain"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"slices"
)

type BatchScanner struct {
	CataloguerFactory CataloguerFactory
	DetailsReaders    []DetailsReader
}

func (s *BatchScanner) Scan(ctx context.Context, owner ownermodel.Owner, volume SourceVolume, optionSlice ...Options) ([]*ScannedFolder, error) {

	options := ReduceOptions(optionSlice...)

	launcher, reportBuilder, err := s.prepareVolumeScan(ctx, options, volume.String(), owner)
	if err != nil {
		return nil, err
	}

	err = <-launcher.Process(ctx, volume)

	return reportBuilder.build(), err
}

// scanListeners list the listeners that will be notified during the scan process.
type scanConfiguration struct {
	Analyser                 Analyser
	Cataloguer               Cataloguer
	ScanCompleteObserver     scanCompleteObserver
	PostAnalyserRejects      []RejectedMediaObserver
	PostCatalogFiltersIn     []CatalogReferencerObserver
	PostCataloguerFiltersOut []CataloguerFilterObserver
	Wrappers                 []chain.CloserFunc
}

func (s *BatchScanner) prepareVolumeScan(ctx context.Context, options Options, volumeName string, owner ownermodel.Owner) (analyserLauncher, *scanReportBuilder, error) {
	tracker, _ := newTrackerV2(options)
	reportBuilder := newScanReportBuilder()
	scanLogger := newLogger(volumeName)

	cataloguer, err := s.CataloguerFactory.NewOwnerScopedCataloguer(ctx, owner)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create cataloguer for owner %s", owner)
	}

	config := &scanConfiguration{
		Analyser:                 options.GetAnalyserDecorator().Decorate(newDefaultAnalyser(s.DetailsReaders...), tracker),
		Cataloguer:               cataloguer,
		ScanCompleteObserver:     tracker,
		PostAnalyserRejects:      []RejectedMediaObserver{scanLogger, tracker},
		PostCatalogFiltersIn:     []CatalogReferencerObserver{scanLogger, tracker, reportBuilder},
		PostCataloguerFiltersOut: []CataloguerFilterObserver{scanLogger, tracker},
		Wrappers:                 []chain.CloserFunc{tracker.NoMoreEvents},
	}
	if !options.SkipRejects {
		config.PostAnalyserRejects = append(config.PostAnalyserRejects, new(analyserFailsFastObserver))
	}
	config.PostAnalyserRejects = append(config.PostAnalyserRejects, reportBuilder)

	launcher, err := multithreadedScanRuntime(ctx, options, config)
	return launcher, reportBuilder, err
}

func multithreadedScanRuntime(ctxNonCancelable context.Context, options Options, config *scanConfiguration) (analyserLauncher, error) {
	ctx, cancelFunc := context.WithCancel(ctxNonCancelable)

	launcher := scanAndBackupCommonLauncher(config, options, &chain.MultithreadedLink[[]BackingUpMediaRequest, []BackingUpMediaRequest]{
		NumberOfRoutines: 1,
		ConsumerBuilder:  chain.PassThrough[[]BackingUpMediaRequest](),
		Next: &chain.CloseWrapperLink[[]BackingUpMediaRequest]{
			CloserFuncs: slices.Concat(config.Wrappers, []chain.CloserFunc{chain.CloserFunc(cancelFunc)}),
			Next:        chain.EndOfTheChain[[]BackingUpMediaRequest](finalizer(config.PostCatalogFiltersIn)...),
		},
	})

	err := launcher.Starts(ctx, chain.NewErrorCollector(func(err error) {
		cancelFunc()
	}))
	return launcher, err
}
