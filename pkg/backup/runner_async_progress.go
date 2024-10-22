package backup

import (
	"context"
)

// NewProgressObserver publishes on a channel the progress of the backup
func NewProgressObserver(sizeHint int) *ProgressObserver {
	return &ProgressObserver{
		EventChannel: make(chan *ProgressEvent, sizeHint*5),
	}
}

type ProgressObserver struct {
	EventChannel chan *ProgressEvent
}

func (p *ProgressObserver) OnRejectedMedia(ctx context.Context, found FoundMedia, err error) {
	p.EventChannel <- &ProgressEvent{Type: ProgressEventAnalysed, Count: 1, Size: found.Size()}
	p.EventChannel <- &ProgressEvent{Type: ProgressEventRejected, Count: 1, Size: found.Size()} // Hacky - Should be SKIPPED not ProgressEventDuplicate
}

func (p *ProgressObserver) OnDecoratedAnalyser(ctx context.Context, found FoundMedia, cacheHit bool) error {
	if cacheHit {
		p.EventChannel <- &ProgressEvent{Type: ProgressEventAnalysedFromCache, Count: 1, Size: found.Size()}
	}

	return nil
}

func (p *ProgressObserver) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	p.EventChannel <- &ProgressEvent{Type: ProgressEventAnalysed, Count: 1, Size: media.FoundMedia.Size()}
	return nil
}
