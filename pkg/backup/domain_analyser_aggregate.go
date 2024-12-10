package backup

import (
	"context"
	"github.com/pkg/errors"
)

var (
	ErrAnalyserNoDateTime   = errors.New("media must have a date time included in the metadata")
	ErrAnalyserNotSupported = errors.New("media format is not supported")
)

type AnalysedMediaObserver interface {
	OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error
}

type RejectedMediaObserver interface {
	// OnRejectedMedia is called when the media is invalid and cannot be used ; the error is returned only if there is a technical issue.
	OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error
}

type analyserAggregate struct {
	analyser               Analyser
	analysedMediaObservers AnalysedMediaObservers
	rejectedMediaObservers RejectedMediaObservers
}

func (a *analyserAggregate) OnFoundMedia(ctx context.Context, media FoundMedia) error {
	analyseMedia, err := a.analyser.Analyse(ctx, media)
	if err != nil {
		return a.rejectedMediaObservers.OnRejectedMedia(ctx, media, err)
	}

	if analyseMedia.Details.DateTime.IsZero() {
		return a.rejectedMediaObservers.OnRejectedMedia(ctx, media, ErrAnalyserNoDateTime)
	}

	return a.analysedMediaObservers.OnAnalysedMedia(ctx, analyseMedia)
}

type AnalysedMediaObserverFunc func(ctx context.Context, media *AnalysedMedia) error

func (a AnalysedMediaObserverFunc) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	return a(ctx, media)
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
