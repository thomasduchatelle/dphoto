package backup

import (
	"context"
	"github.com/pkg/errors"
)

// scanningController is the engine running the scanning code: it could be synchronous or multithreaded.
type scanningController interface {
	scanMonitoringIntegrator
	Launcher(analyser Analyser, chain *analyserObserverChain, tracker scanCompleteObserver) scanningLauncher
	WrapAnalysedMediasBatchObserverIntoAnalysedMediaObserver(ctx context.Context, batchSize int, cataloguerAdapter analysedMediasBatchObserver) AnalysedMediaObserver
}

type scanningLauncher analyserLauncher

type scanningOptions struct {
	Options
	analyser   Analyser   // analyser is mandatory
	cataloguer Cataloguer // cataloguer is mandatory
}

func newScanningChain(ctx context.Context, controller scanningController, options scanningOptions) (scanningLauncher, error) {
	if options.analyser == nil || options.cataloguer == nil {
		return nil, errors.New("scanningOptions.analyser and scanningOptions.cataloguer are mandatory.")
	}

	var postAnalyserRejects []RejectedMediaObserver
	if !options.SkipRejects {
		postAnalyserRejects = append(postAnalyserRejects, new(analyserFailsFastObserver))
	}

	chain := &analyserObserverChain{
		AnalysedMediaObservers: controller.AppendPostAnalyserSuccess(&analyserNoDateTimeFilter{
			analyserObserverChain{
				AnalysedMediaObservers: []AnalysedMediaObserver{
					controller.WrapAnalysedMediasBatchObserverIntoAnalysedMediaObserver(ctx, options.BatchSize, &analyserToCatalogReferencer{
						CatalogReferencer: options.cataloguer,
						CatalogReferencerObservers: controller.AppendPreCataloguerFilter(&applyFiltersOnCataloguer{
							CatalogReferencerObservers: controller.AppendPostCatalogFiltersIn(),
							CataloguerFilterObservers:  controller.AppendPostCatalogFiltersOut(),
							CataloguerFilters:          postCatalogFiltersList(options.Options),
						}),
					}),
				},
				RejectedMediaObservers: controller.AppendPostAnalyserFilterRejects(),
			},
		}),
		RejectedMediaObservers: controller.AppendPostAnalyserRejects(postAnalyserRejects...),
	}

	return controller.Launcher(options.analyser, chain, controller.ScanCompleteObserver()), nil
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
