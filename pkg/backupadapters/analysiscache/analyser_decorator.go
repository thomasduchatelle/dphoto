package analysiscache

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
)

func (d *AnalyserCache) Decorate(analyser backup.Analyser, observers ...backup.AnalyserDecoratorObserver) backup.Analyser {
	return &AnalyserCacheWrapper{
		Delegate:      analyser,
		AnalyserCache: d,
		Observers:     observers,
	}
}

type AnalyserCacheWrapper struct {
	Delegate      backup.Analyser
	AnalyserCache *AnalyserCache
	Observers     []backup.AnalyserDecoratorObserver
}

func (a *AnalyserCacheWrapper) Analyse(ctx context.Context, found backup.FoundMedia, analysedMediaObserver backup.AnalysedMediaObserver) error {
	cache, missed, err := a.AnalyserCache.RestoreCache(found)
	if err != nil {
		return err
	}

	if !missed {
		err = a.fireHit(ctx, found)
		if err != nil {
			return err
		}
		return analysedMediaObserver.OnAnalysedMedia(ctx, cache)
	}

	err = a.fireMissed(ctx, found)
	if err != nil {
		return err
	}
	return a.Delegate.Analyse(ctx, found, &interceptor{
		AnalysedMediaObserver: analysedMediaObserver,
		AnalyserCache:         a.AnalyserCache,
	})
}

func (a *AnalyserCacheWrapper) fireMissed(ctx context.Context, found backup.FoundMedia) error {
	for _, observer := range a.Observers {
		err := observer.OnDecoratedAnalyser(ctx, found, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AnalyserCacheWrapper) fireHit(ctx context.Context, found backup.FoundMedia) error {
	for _, observer := range a.Observers {
		err := observer.OnDecoratedAnalyser(ctx, found, true)
		if err != nil {
			return err
		}
	}
	return nil
}

type interceptor struct {
	AnalysedMediaObserver backup.AnalysedMediaObserver
	AnalyserCache         *AnalyserCache
}

func (a *interceptor) OnAnalysedMedia(ctx context.Context, media *backup.AnalysedMedia) error {
	err := a.AnalyserCache.StoreCache(media)
	if err != nil {
		return err
	}

	return a.AnalysedMediaObserver.OnAnalysedMedia(ctx, media)
}
