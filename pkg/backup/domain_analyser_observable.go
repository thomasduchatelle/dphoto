package backup

import (
	"context"
	"github.com/pkg/errors"
)

type Analyser interface {
	Analyse(ctx context.Context, found FoundMedia, analysedMediaObserver AnalysedMediaObserver, rejectsObserver RejectedMediaObserver) error
}

func newDefaultAnalyser(readers ...DetailsReader) *observableAnalyser {
	return &observableAnalyser{
		analyser: &mediaReader{
			detailsReaders: readers,
		},
	}
}

type observableAnalyser struct {
	analyser *mediaReader
}

func (a *observableAnalyser) Analyse(ctx context.Context, found FoundMedia, analysedMediaObserver AnalysedMediaObserver, rejectsObserver RejectedMediaObserver) error {
	media, err := a.analyser.analyseMedia(found)
	if err != nil {
		return rejectsObserver.OnRejectedMedia(ctx, found, err)
	}
	return analysedMediaObserver.OnAnalysedMedia(ctx, media)
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

type AnalyserDecorator interface {
	Decorate(analyse Analyser, observers ...AnalyserDecoratorObserver) Analyser
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
