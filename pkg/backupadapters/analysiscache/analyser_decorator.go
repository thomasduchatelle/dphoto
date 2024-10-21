package analysiscache

import (
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

func (a *AnalyserCacheWrapper) Analyse(found backup.FoundMedia, analysedMediaObserver backup.AnalysedMediaObserver, rejectedMediaObserver backup.RejectedMediaObserver) {
	cache, missed, err := a.AnalyserCache.RestoreCache(found)
	if err != nil {
		a.fireMissed(found)
		rejectedMediaObserver.OnRejectedMedia(found, err)
		return
	}

	if !missed {
		a.fireHit(found)
		analysedMediaObserver.OnAnalysedMedia(cache)
		return
	}

	a.fireMissed(found)
	a.Delegate.Analyse(found, &interceptor{
		AnalysedMediaObserver: analysedMediaObserver,
		RejectedMediaObserver: rejectedMediaObserver,
		AnalyserCache:         a.AnalyserCache,
	}, rejectedMediaObserver)
}

func (a *AnalyserCacheWrapper) fireMissed(found backup.FoundMedia) {
	for _, observer := range a.Observers {
		observer.OnDecoratedAnalyser(found, false)
	}
}

func (a *AnalyserCacheWrapper) fireHit(found backup.FoundMedia) {
	for _, observer := range a.Observers {
		observer.OnDecoratedAnalyser(found, true)
	}
}

type interceptor struct {
	AnalysedMediaObserver backup.AnalysedMediaObserver
	RejectedMediaObserver backup.RejectedMediaObserver
	AnalyserCache         *AnalyserCache
}

func (a *interceptor) OnAnalysedMedia(media *backup.AnalysedMedia) {
	err := a.AnalyserCache.StoreCache(media)
	if err != nil {
		a.RejectedMediaObserver.OnRejectedMedia(media.FoundMedia, err)
		return
	}

	a.AnalysedMediaObserver.OnAnalysedMedia(media)
}
