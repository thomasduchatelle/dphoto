package backup

import (
	"context"
	"github.com/pkg/errors"
)

// NewProgressObserver publishes on a channel the progress of the backup
func NewProgressObserver(sizeHint int) *ProgressObserver {
	return &ProgressObserver{
		EventChannel: make(chan *progressEvent, sizeHint*5),
	}
}

type ProgressObserver struct {
	EventChannel chan *progressEvent
}

func (p *ProgressObserver) OnUniqueBackupMediaRequest(ctx context.Context, request []BackingUpMediaRequest) error {
	for _, media := range request {
		p.EventChannel <- &progressEvent{Type: trackCatalogued, Count: 1, Size: media.AnalysedMedia.FoundMedia.Size()}
	}

	return nil
}

func (p *ProgressObserver) OnDuplicatedBackupMediaRequest(ctx context.Context, request []BackingUpMediaRequest) error {
	for _, media := range request {
		p.EventChannel <- &progressEvent{Type: trackDuplicatedInVolume, Count: 1, Size: media.AnalysedMedia.FoundMedia.Size()}
	}

	return nil
}

func (p *ProgressObserver) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	p.EventChannel <- &progressEvent{Type: trackAnalysisFailed, Count: 1, Size: found.Size()}
	return nil
}

func (p *ProgressObserver) OnDecoratedAnalyser(ctx context.Context, found FoundMedia, cacheHit bool) error {
	if cacheHit {
		p.EventChannel <- &progressEvent{Type: trackAnalysedFromCache, Count: 1, Size: found.Size()}
	}

	return nil
}

func (p *ProgressObserver) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	p.EventChannel <- &progressEvent{Type: ProgressEventAnalysed, Count: 1, Size: media.FoundMedia.Size()}
	return nil
}

func (p *ProgressObserver) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	count := MediaCounterZero
	for _, request := range requests {
		if request.CatalogReference.AlbumCreated() {
			p.EventChannel <- &progressEvent{Type: trackAlbumCreated, Count: 1, Album: request.CatalogReference.AlbumFolderName()}
		}

		count = count.Add(1, request.AnalysedMedia.FoundMedia.Size())
	}

	p.EventChannel <- &progressEvent{Type: ProgressEventCatalogued, Count: count.Count, Size: count.Size}

	return nil
}

func (p *ProgressObserver) OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error {
	switch {
	case errors.Is(cause, ErrCatalogerFilterMustBeInAlbum):
		p.EventChannel <- &progressEvent{Type: trackWrongAlbum, Count: 1, Size: media.FoundMedia.Size()}
		return nil
	case errors.Is(cause, ErrCatalogerFilterMustNotAlreadyExists):
		p.EventChannel <- &progressEvent{Type: trackAlreadyExistsInCatalog, Count: 1, Size: media.FoundMedia.Size()}
		return nil

	case errors.Is(cause, ErrCatalogerFilterMustNotAlreadyExists):
		p.EventChannel <- &progressEvent{Type: trackAlreadyExistsInCatalog, Count: 1, Size: media.FoundMedia.Size()}
		return nil

	default:
		return errors.Wrapf(cause, "filter error is not supported. Media: %s", media.FoundMedia)
	}
}
