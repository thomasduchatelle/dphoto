package backup

type IProgressObserver interface {
	AnalysedMediaObserver
	RejectedMediaObserver
	AnalyserDecoratorObserver
}

func ProgressObserverMustImplementsObserverInterfaces() IProgressObserver {
	return &ProgressObserver{}
}
