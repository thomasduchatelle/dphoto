package catalog

import (
	"context"
	"github.com/pkg/errors"
	"slices"
	"time"
)

type TimelineAggregate struct {
	timeline *Timeline
	albums   []*Album
}

// NewLazyTimelineAggregate creates a new TimelineAggregate without timeline pre-computation. The timeline will be computed at the first AddNew call.
func NewLazyTimelineAggregate(albums []*Album) *TimelineAggregate {
	return &TimelineAggregate{
		albums: albums,
	}
}

func NewInitialisedTimelineAggregate(albums []*Album) (*TimelineAggregate, error) {
	timeline, err := NewTimeline(albums)
	return &TimelineAggregate{
		albums:   albums,
		timeline: timeline,
	}, err
}

func (t *TimelineAggregate) CreateNewAlbum(ctx context.Context, request CreateAlbumRequest, observers ...CreateAlbumObserver) (*AlbumId, error) {
	if err := request.IsValid(); err != nil {
		return nil, errors.Wrapf(err, "CreateNewAlbum(%s) failed", request)
	}

	album, err := t.convert(request)
	if err != nil {
		return nil, err
	}

	for index, observer := range observers {
		if err = observer.ObserveCreateAlbum(ctx, *album); err != nil {
			return nil, errors.Wrapf(err, "CreateNewAlbum(%s) failed at observer %d/%d", request, index, len(observers))
		}
	}

	return &album.AlbumId, nil
}

func (t *TimelineAggregate) convert(request CreateAlbumRequest) (*Album, error) {
	folderName := generateFolderName(request.Name, request.Start)
	if request.ForcedFolderName != "" {
		folderName = NewFolderName(request.ForcedFolderName)
	}

	albumId := AlbumId{
		Owner:      request.Owner,
		FolderName: folderName,
	}

	nameIsAlreadyTaken := slices.ContainsFunc(t.albums, func(album *Album) bool {
		return album.AlbumId.IsEqual(albumId)
	})
	if nameIsAlreadyTaken {
		return nil, errors.Wrapf(AlbumFolderNameAlreadyTakenErr, "%s album id already exists", albumId)
	}

	return &Album{
		AlbumId: albumId,
		Name:    request.Name,
		Start:   request.Start,
		End:     request.End,
	}, nil
}

func (t *TimelineAggregate) AddNew(addedAlbum Album) (MediaTransferRecords, error) {
	t.albums = append(t.albums, &addedAlbum)

	var err error
	t.timeline, err = NewTimeline(t.albums)
	if err != nil {
		return nil, err
	}

	records := make(MediaTransferRecords)
	for _, seg := range t.timeline.FindForAlbum(addedAlbum.AlbumId) {
		if len(seg.Albums) > 1 {
			selector := MediaSelector{
				FromAlbums: extractAlbumIds(seg.Albums[1:]),
				Start:      seg.Start,
				End:        seg.End,
			}
			if selectors, found := records[addedAlbum.AlbumId]; found {
				records[addedAlbum.AlbumId] = append(selectors, selector)
			} else {
				records[addedAlbum.AlbumId] = []MediaSelector{selector}
			}
		}
	}

	if len(records) == 0 {
		return nil, nil
	}

	return records, nil
}

func (t *TimelineAggregate) FindAt(date time.Time) (*Album, bool, error) {
	if t.timeline == nil {
		var err error
		t.timeline, err = NewTimeline(t.albums)
		if err != nil {
			return nil, false, errors.Wrapf(err, "failed to create timeline during FindAt(%s)", date)
		}
	}

	albumId, found := t.timeline.FindAt(date)
	return albumId, found, nil
}
