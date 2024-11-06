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

// scanListeners list the listeners that will be notified during the scan process.
type scanListeners struct {
	scanCompleteObserver      scanCompleteObserver
	PostAnalyserSuccess       []AnalysedMediaObserver
	PostAnalyserFilterRejects []RejectedMediaObserver
	PostAnalyserRejects       []RejectedMediaObserver
	PreCataloguerFilter       []CatalogReferencerObserver
	PostCatalogFiltersIn      []CatalogReferencerObserver
	PostCatalogFiltersOut     []CataloguerFilterObserver
}

func (m *scanListeners) ScanCompleteObserver() scanCompleteObserver {
	return m.scanCompleteObserver
}

func (m *scanListeners) AppendPostCatalogFiltersOut(observers ...CataloguerFilterObserver) []CataloguerFilterObserver {
	list := make([]CataloguerFilterObserver, len(m.PostCatalogFiltersOut), len(m.PostCatalogFiltersOut)+len(observers))
	_ = copy(list, m.PostCatalogFiltersOut)
	list = append(list, observers...)

	return list
}

func (m *scanListeners) AppendPostCatalogFiltersIn(observers ...CatalogReferencerObserver) []CatalogReferencerObserver {
	list := make([]CatalogReferencerObserver, len(m.PostCatalogFiltersIn), len(m.PostCatalogFiltersIn)+len(observers))
	_ = copy(list, m.PostCatalogFiltersIn)
	list = append(list, observers...)

	return list
}

func (m *scanListeners) AppendPreCataloguerFilter(observers ...CatalogReferencerObserver) []CatalogReferencerObserver {
	list := make([]CatalogReferencerObserver, len(m.PreCataloguerFilter), len(m.PreCataloguerFilter)+len(observers))
	_ = copy(list, m.PreCataloguerFilter)
	list = append(list, observers...)

	return list
}

func (m *scanListeners) AppendPostAnalyserSuccess(observers ...AnalysedMediaObserver) []AnalysedMediaObserver {
	list := make([]AnalysedMediaObserver, len(m.PostAnalyserSuccess), len(m.PostAnalyserSuccess)+len(observers))
	_ = copy(list, m.PostAnalyserSuccess)
	list = append(list, observers...)

	return list
}

func (m *scanListeners) AppendPostAnalyserRejects(observers ...RejectedMediaObserver) []RejectedMediaObserver {
	list := make([]RejectedMediaObserver, len(m.PostAnalyserRejects), len(m.PostAnalyserRejects)+len(observers)+1)
	_ = copy(list, m.PostAnalyserRejects)
	list = append(list, observers...)

	return list
}

func (m *scanListeners) AppendPostAnalyserFilterRejects(observers ...RejectedMediaObserver) []RejectedMediaObserver {
	list := make([]RejectedMediaObserver, len(m.PostAnalyserFilterRejects), len(m.PostAnalyserFilterRejects)+len(observers))
	_ = copy(list, m.PostAnalyserFilterRejects)
	list = append(list, observers...)

	return list
}
