package backup

import (
	"context"
	"sync"
)

func newBackupReportBuilder() *backupReportBuilder {
	return &backupReportBuilder{
		lock:          sync.Mutex{},
		countPerAlbum: make(map[string]IAlbumReport),
	}
}

type backupReportBuilder struct {
	lock          sync.Mutex
	skipped       MediaCounter
	countPerAlbum map[string]IAlbumReport
}

func (r *backupReportBuilder) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.skipped = r.skipped.Add(1, found.Size())

	return nil
}

func (r *backupReportBuilder) OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.skipped = r.skipped.Add(1, media.FoundMedia.Size())

	return nil
}

func (r *backupReportBuilder) OnBackingUpMediaRequestUploaded(ctx context.Context, request BackingUpMediaRequest) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	albumName := request.CatalogReference.AlbumFolderName()
	count, ok := r.countPerAlbum[albumName]
	if !ok {
		count = &albumReport{}
		r.countPerAlbum[albumName] = count
	}

	if request.CatalogReference.AlbumCreated() {
		count.(*albumReport).isNew = true
	}

	switch request.AnalysedMedia.Type {
	case MediaTypeImage:
		count.(*albumReport).image = count.(*albumReport).image.Add(1, request.AnalysedMedia.FoundMedia.Size())
	case MediaTypeVideo:
		count.(*albumReport).video = count.(*albumReport).video.Add(1, request.AnalysedMedia.FoundMedia.Size())
	default:
		count.(*albumReport).other = count.(*albumReport).other.Add(1, request.AnalysedMedia.FoundMedia.Size())
	}

	return nil
}

func (r *backupReportBuilder) Skipped() MediaCounter {
	return r.skipped
}

func (r *backupReportBuilder) CountPerAlbum() map[string]IAlbumReport {
	return r.countPerAlbum
}

type albumReport struct {
	isNew bool
	image MediaCounter
	video MediaCounter
	other MediaCounter
}

func (c *albumReport) IsNew() bool {
	return c.isNew
}

func (c *albumReport) Total() MediaCounter {
	return c.image.AddCounter(c.video).AddCounter(c.other)
}

func (c *albumReport) OfType(mediaType MediaType) MediaCounter {
	switch mediaType {
	case MediaTypeImage:
		return c.image
	case MediaTypeVideo:
		return c.video
	default:
		return c.other
	}
}
