package backup

type scanMonitoringIntegrator interface {
	ScanCompleteObserver() scanCompleteObserver
	AppendPostAnalyserSuccess(observers ...AnalysedMediaObserver) []AnalysedMediaObserver
	AppendPostAnalyserRejects(observers ...RejectedMediaObserver) []RejectedMediaObserver
	AppendPostAnalyserFilterRejects(observers ...RejectedMediaObserver) []RejectedMediaObserver
	AppendPreCataloguerFilter(observers ...CatalogReferencerObserver) []CatalogReferencerObserver
	AppendPostCatalogFiltersIn(observers ...CatalogReferencerObserver) []CatalogReferencerObserver
	AppendPostCatalogFiltersOut(observers ...CataloguerFilterObserver) []CataloguerFilterObserver
}
