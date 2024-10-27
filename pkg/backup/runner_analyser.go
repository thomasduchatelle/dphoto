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

func newAnalyserObserverChain(options Options, observers ...interface{}) *analyserObserverChain {
	delegates := newCompositeAnalyserObserver(observers...)

	if options.SkipRejects {
		return &analyserObserverChain{
			AnalysedMediaObservers: []AnalysedMediaObserver{
				&analyserNoDateTimeFilter{
					analyserObserverChain{
						AnalysedMediaObservers: []AnalysedMediaObserver{delegates},
						RejectedMediaObservers: []RejectedMediaObserver{delegates}},
				},
			},
			RejectedMediaObservers: []RejectedMediaObserver{delegates},
		}
	}

	return &analyserObserverChain{
		AnalysedMediaObservers: []AnalysedMediaObserver{
			&analyserNoDateTimeFilter{
				analyserObserverChain{
					AnalysedMediaObservers: []AnalysedMediaObserver{delegates},
					RejectedMediaObservers: []RejectedMediaObserver{new(analyserFailsFastObserver)},
				},
			},
		},
		RejectedMediaObservers: []RejectedMediaObserver{delegates, new(analyserFailsFastObserver)},
	}
}

type analyserObserverChain struct {
	AnalysedMediaObservers []AnalysedMediaObserver
	RejectedMediaObservers []RejectedMediaObserver
}

func (a *analyserObserverChain) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	for _, observer := range a.AnalysedMediaObservers {
		if err := observer.OnAnalysedMedia(ctx, media); err != nil {
			return err
		}
	}

	return nil
}

func (a *analyserObserverChain) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	for _, observer := range a.RejectedMediaObservers {
		if err := observer.OnRejectedMedia(ctx, found, cause); err != nil {
			return err
		}
	}

	return nil
}

type analyserNoDateTimeFilter struct {
	analyserObserverChain
}

func (a *analyserNoDateTimeFilter) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	if media.Details.DateTime.IsZero() {
		return a.analyserObserverChain.OnRejectedMedia(ctx, media.FoundMedia, ErrAnalyserNoDateTime)
	}

	return a.analyserObserverChain.OnAnalysedMedia(ctx, media)
}

type analyserFailsFastObserver struct {
}

func (a *analyserFailsFastObserver) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	return errors.Wrapf(cause, "invalid media '%s'", found)
}
