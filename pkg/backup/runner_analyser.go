package backup

import (
	"context"
	"github.com/pkg/errors"
)

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

type RejectedMediaObserver interface {
	// OnRejectedMedia is called when the media is invalid and cannot be used ; the error is returned only if there is a technical issue.
	OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error
}

type Analyser interface {
	Analyse(ctx context.Context, found FoundMedia, analysedMediaObserver AnalysedMediaObserver, rejectsObserver RejectedMediaObserver) error
}

type AnalyserDecorator interface {
	Decorate(analyseFunc Analyser, observers ...AnalyserDecoratorObserver) Analyser
}

type AnalyserDecoratorObserver interface {
	OnDecoratedAnalyser(ctx context.Context, found FoundMedia, cacheHit bool) error
}

// ... unsure ...

type AnalyserObserver interface {
	AnalysedMediaObserver
	RejectedMediaObserver
}

func NewAnalyserMediaHandler(options Options, observers ...AnalyserObserver) AnalyserObserver {
	delegates := NewCompositeAnalyserObserver(observers...)
	if options.SkipRejects {
		return &analyserNoDateTimeFilter{observer: delegates}
	}

	return &analyserNoDateTimeFilter{observer: &analyserFailsFastObserver{observer: delegates}}
}

type analyserNoDateTimeFilter struct {
	observer AnalyserObserver
}

func (a *analyserNoDateTimeFilter) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	return a.observer.OnRejectedMedia(ctx, found, cause) // TODO is that required to have a pass-through?
}

func (a *analyserNoDateTimeFilter) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	if media.Details.DateTime.IsZero() {
		return a.observer.OnRejectedMedia(ctx, media.FoundMedia, ErrAnalyserNoDateTime)
	}

	return a.observer.OnAnalysedMedia(ctx, media)
}

type analyserFailsFastObserver struct {
	observer AnalyserObserver
}

func (a *analyserFailsFastObserver) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	return a.observer.OnAnalysedMedia(ctx, media) // TODO is that required to have a pass-through?
}

func (a *analyserFailsFastObserver) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	return errors.Wrapf(cause, "invalid media '%s'", found)
}
