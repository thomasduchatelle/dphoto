package backup

var (
	defaultAnalyser = &CoreAnalyser{
		options: DetailsReaderOptions{},
	}
)

func getDefaultAnalyser() Analyser {
	return defaultAnalyser
}

func ClearDetailsReader() {
	defaultAnalyser.detailsReaders = nil
}

func RegisterDetailsReader(reader DetailsReaderAdapter) {
	defaultAnalyser.detailsReaders = append(defaultAnalyser.detailsReaders, reader)
}

type Analyser interface {
	Analyse(found FoundMedia, analysedMediaObserver AnalysedMediaObserver, rejectedMediaObserver RejectedMediaObserver)
}

type AnalysedMediaObserver interface {
	OnAnalysedMedia(media *AnalysedMedia)
}

type RejectedMediaObserver interface {
	OnRejectedMedia(found FoundMedia, err error)
}

type AnalyserDecorator interface {
	Decorate(analyseFunc Analyser, observers ...AnalyserDecoratorObserver) Analyser
}

type AnalyserDecoratorObserver interface {
	OnDecoratedAnalyser(found FoundMedia, cacheHit bool)
}

type RunnerAnalyserFunc func(found FoundMedia, analysedMediaObserver AnalysedMediaObserver, rejectedMediaObserver RejectedMediaObserver)

func (r RunnerAnalyserFunc) Analyse(found FoundMedia, analysedMediaObserver AnalysedMediaObserver, rejectedMediaObserver RejectedMediaObserver) {
	r(found, analysedMediaObserver, rejectedMediaObserver)
}
