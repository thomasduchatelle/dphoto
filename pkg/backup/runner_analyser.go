package backup

import "context"

var (
	defaultAnalyser = &CoreAnalyser{
		options: DetailsReaderOptions{},
	}
)

func getDefaultAnalyser() Analyser {
	return defaultAnalyser
}

// ClearDetailsReader remove all details reader from the default Analyser
func ClearDetailsReader() {
	defaultAnalyser.detailsReaders = nil
}

// RegisterDetailsReader adds a details reader implementation to the default Analyser
func RegisterDetailsReader(reader DetailsReaderAdapter) {
	defaultAnalyser.detailsReaders = append(defaultAnalyser.detailsReaders, reader)
}

type AnalysedMediaObserver interface {
	OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error
}

type Analyser interface {
	Analyse(ctx context.Context, found FoundMedia, analysedMediaObserver AnalysedMediaObserver) error
}

type AnalyserDecorator interface {
	Decorate(analyseFunc Analyser, observers ...AnalyserDecoratorObserver) Analyser
}

type AnalyserDecoratorObserver interface {
	OnDecoratedAnalyser(ctx context.Context, found FoundMedia, cacheHit bool) error
}

// ... Async ...

type RejectedMediaObserver interface {
	OnRejectedMedia(ctx context.Context, found FoundMedia, cause error)
}

type AnalyserAsync interface {
	Analyse(ctx context.Context, found FoundMedia, analysedMediaObserver AnalysedMediaObserver, rejectedMediaObserver RejectedMediaObserver)
}

type AnalyserAsyncWrapper struct {
	Analyser Analyser
}

func (a *AnalyserAsyncWrapper) Analyse(ctx context.Context, found FoundMedia, analysedMediaObserver AnalysedMediaObserver, rejectedMediaObserver RejectedMediaObserver) {
	err := a.Analyser.Analyse(ctx, found, analysedMediaObserver)
	if err != nil {
		rejectedMediaObserver.OnRejectedMedia(ctx, found, err)
	}
}
