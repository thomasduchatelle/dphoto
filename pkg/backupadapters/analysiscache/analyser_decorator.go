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

func (a *AnalyserCacheWrapper) Analyse(ctx context.Context, found backup.FoundMedia) (*backup.AnalysedMedia, error) {
	cache, missed, err := a.AnalyserCache.RestoreCache(found)
	if err != nil {
		return nil, err
	}

	if !missed {
		err = a.fireHit(ctx, found)
		if err != nil {
			return nil, err
		}
		return cache, nil
	}

	analyse, err := a.Delegate.Analyse(ctx, found)
	if err != nil {
		return nil, err
	}

	err = a.AnalyserCache.StoreCache(analyse)
	return analyse, err
}

func (a *AnalyserCacheWrapper) fireHit(ctx context.Context, found backup.FoundMedia) error {
	for _, observer := range a.Observers {
		err := observer.OnSkipDelegateAnalyser(ctx, found)
		if err != nil {
			return err
		}
	}
	return nil
}
