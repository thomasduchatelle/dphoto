package backup

type IProgressObserver interface {
	AnalysedMediaObserver
	RejectedMediaObserver
	AnalyserDecoratorObserver
	CatalogReferencerObserver
	CataloguerFilterObserver
}

func ProgressObserverMustImplementsObserverInterfaces() IProgressObserver {
	return &ProgressObserver{}
}
