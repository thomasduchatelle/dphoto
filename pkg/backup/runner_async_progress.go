package backup

import (
	"context"
	"github.com/pkg/errors"
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

func (p *ProgressObserver) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	p.EventChannel <- &ProgressEvent{Type: ProgressEventAnalysed, Count: 1, Size: found.Size()}
	p.EventChannel <- &ProgressEvent{Type: ProgressEventRejected, Count: 1, Size: found.Size()} // Hacky - Should be SKIPPED not ProgressEventDuplicate
	return nil
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

func (p *ProgressObserver) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	count := MediaCounterZero
	for _, request := range requests {
		if request.CatalogReference.AlbumCreated() {
			p.EventChannel <- &ProgressEvent{Type: ProgressEventAlbumCreated, Count: 1, Album: request.CatalogReference.AlbumFolderName()}
		}

		count = count.Add(1, request.AnalysedMedia.FoundMedia.Size())
	}

	p.EventChannel <- &ProgressEvent{Type: ProgressEventCatalogued, Count: count.Count, Size: count.Size}

	return nil
}

func (p *ProgressObserver) OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error {
	switch {
	case errors.Is(cause, ErrCatalogerFilterMustBeInAlbum):
		p.EventChannel <- &ProgressEvent{Type: ProgressEventWrongAlbum, Count: 1, Size: media.FoundMedia.Size()}
		return nil
	case errors.Is(cause, ErrCatalogerFilterMustNotAlreadyExists):
		p.EventChannel <- &ProgressEvent{Type: ProgressEventAlreadyExists, Count: 1, Size: media.FoundMedia.Size()}
		return nil

	default:
		return errors.Wrapf(cause, "filter error is not supported. Media: %s", media.FoundMedia)
	}
}
