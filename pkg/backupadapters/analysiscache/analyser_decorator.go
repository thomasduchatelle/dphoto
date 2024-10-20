package analysiscache

import (
	"github.com/thomasduchatelle/dphoto/pkg/backup"
)

func (d *AnalyserCache) Decorate(analyser backup.RunnerAnalyser) backup.RunnerAnalyser {
	return &AnalyserCacheWrapper{
		Delegate:      analyser,
		AnalyserCache: d,
	}
}

type AnalyserCacheWrapper struct {
	Delegate      backup.RunnerAnalyser
	AnalyserCache *AnalyserCache
}

func (a *AnalyserCacheWrapper) Analyse(found backup.FoundMedia, analysedMediaObserver backup.AnalysedMediaObserver, rejectedMediaObserver backup.RejectedMediaObserver) {
	cache, missed, err := a.AnalyserCache.RestoreCache(found)
	if err != nil {
		rejectedMediaObserver.OnRejectedMedia(found, err)
		return
	}

	if !missed {
		analysedMediaObserver.OnAnalysedMedia(cache)
		return
	}

	a.Delegate.Analyse(found, &AnalyserCacheObserver{
		AnalysedMediaObserver: analysedMediaObserver,
		RejectedMediaObserver: rejectedMediaObserver,
		AnalyserCache:         a.AnalyserCache,
	}, rejectedMediaObserver)
}

type AnalyserCacheObserver struct {
	AnalysedMediaObserver backup.AnalysedMediaObserver
	RejectedMediaObserver backup.RejectedMediaObserver
	AnalyserCache         *AnalyserCache
}

func (a *AnalyserCacheObserver) OnAnalysedMedia(media *backup.AnalysedMedia) {
	err := a.AnalyserCache.StoreCache(media)
	if err != nil {
		a.RejectedMediaObserver.OnRejectedMedia(media.FoundMedia, err)
		return
	}

	a.AnalysedMediaObserver.OnAnalysedMedia(media)
}
