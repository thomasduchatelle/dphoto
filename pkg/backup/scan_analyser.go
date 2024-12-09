package backup

import (
	"context"
	"github.com/pkg/errors"
)

func newDefaultAnalyser(readers ...DetailsReader) *CoreAnalyser {
	return &CoreAnalyser{
		detailsReaders: readers,
	}
}

type AnalysedMediaObserver interface {
	OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error
}

type AnalysedMediaObserverFunc func(ctx context.Context, media *AnalysedMedia) error

func (a AnalysedMediaObserverFunc) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	return a(ctx, media)
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

type AnalysedMediaObservers []AnalysedMediaObserver

func (a AnalysedMediaObservers) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	for _, observer := range a {
		if err := observer.OnAnalysedMedia(ctx, media); err != nil {
			return err
		}
	}

	return nil
}

type RejectedMediaObservers []RejectedMediaObserver

func (a RejectedMediaObservers) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	for _, observer := range a {
		if err := observer.OnRejectedMedia(ctx, found, cause); err != nil {
			return err
		}
	}

	return nil
}

type analyserFailsFastObserver struct {
}

func (a *analyserFailsFastObserver) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	return errors.Wrapf(cause, "invalid media '%s'", found)
}
