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

func (t *TimelineAggregate) CreateNewAlbum(ctx context.Context, request CreateAlbumRequest, observers ...CreateAlbumObserver) (Album, error) {
	if err := request.IsValid(); err != nil {
		return Album{}, errors.Wrapf(err, "CreateNewAlbum(%s) failed", request)
	}

	album, err := t.convert(request)
	if err != nil {
		return Album{}, err
	}

	return *album, nil
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

func extractAlbumIds(albums []Album) []AlbumId {
	if len(albums) == 0 {
		return nil
	}

	ids := make([]AlbumId, len(albums), len(albums))
	for i, album := range albums {
		ids[i] = album.AlbumId
	}

	return ids
}

func (t *TimelineAggregate) RemoveAlbum(deletedAlbumId AlbumId) (MediaTransferRecords, []MediaSelector, error) {
	records := make(MediaTransferRecords)
	orphaned := make([]MediaSelector, 0)
	var err error
	var deletedAlbum *Album

	t.albums, deletedAlbum = t.removeAlbumFrom(t.albums, deletedAlbumId.FolderName)
	if deletedAlbum == nil {
		return nil, nil, AlbumNotFoundError
	}

	t.timeline, err = NewTimeline(t.albums)
	if err != nil {
		return nil, nil, err
	}

	segments := t.timeline.FindSegmentsBetween(deletedAlbum.Start, deletedAlbum.End)
	for _, seg := range segments {
		selector := MediaSelector{
			FromAlbums: []AlbumId{deletedAlbumId},
			Start:      seg.Start,
			End:        seg.End,
		}

		if len(seg.Albums) == 0 {
			orphaned = append(orphaned, selector)

		} else if priorityDescComparator(deletedAlbum, &seg.Albums[0]) > 0 {
			selectors, _ := records[seg.Albums[0].AlbumId]
			records[seg.Albums[0].AlbumId] = append(selectors, selector)
		}
	}

	return records, orphaned, nil
}

// removeAlbumFrom removes the album with the given folderName from the list of albums
func (t *TimelineAggregate) removeAlbumFrom(albums []*Album, folderName FolderName) ([]*Album, *Album) {
	for index, album := range albums {
		if album.FolderName == folderName {
			return append(albums[:index], albums[index+1:]...), album
		}
	}

	return albums, nil
}

func (t *TimelineAggregate) ValidateAmendDates(albumId AlbumId, start, end time.Time) (*DatesUpdate, error) {
	index := slices.IndexFunc(t.albums, func(album *Album) bool {
		return album.AlbumId.IsEqual(albumId)
	})

	if index == -1 {
		return nil, errors.Wrapf(AlbumNotFoundError, "album %s not found", albumId)
	}

	amended := DatesUpdate{
		UpdatedAlbum:  *t.albums[index],
		PreviousStart: t.albums[index].Start,
		PreviousEnd:   t.albums[index].End,
	}
	amended.UpdatedAlbum.Start = start
	amended.UpdatedAlbum.End = end

	return &amended, nil
}

func (t *TimelineAggregate) AmendDates(amendedAlbum DatesUpdate) (MediaTransferRecords, []MediaSelector, error) {
	index := slices.IndexFunc(t.albums, func(album *Album) bool {
		return album.AlbumId.IsEqual(amendedAlbum.UpdatedAlbum.AlbumId)
	})
	if index == -1 {
		return nil, nil, errors.Wrapf(AlbumNotFoundError, "album %s not found", amendedAlbum.UpdatedAlbum.AlbumId)
	}

	var err error

	previousAlbum := *t.albums[index]
	if !previousAlbum.Start.Equal(amendedAlbum.PreviousStart) || !previousAlbum.End.Equal(amendedAlbum.PreviousEnd) {
		// restore TimelineAggregate as it was before the amendment
		previousAlbum.Start = amendedAlbum.PreviousStart
		previousAlbum.End = amendedAlbum.PreviousEnd
		t.albums[index] = &previousAlbum

		t.timeline, err = NewTimeline(t.albums)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to create timeline during AmendDates(%s)", amendedAlbum)
		}
	}

	originalTimeline, err := t.getTimeline()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create timeline during AmendDates(%s)", amendedAlbum)
	}

	t.albums[index] = &amendedAlbum.UpdatedAlbum
	t.timeline, err = NewTimeline(t.albums)
	if err != nil {
		return nil, nil, err
	}

	start := minTime(amendedAlbum.UpdatedAlbum.Start, previousAlbum.Start)
	end := maxTime(amendedAlbum.UpdatedAlbum.End, previousAlbum.End)

	cursor := struct {
		time             time.Time
		originalSegments []PrioritySegment
		amendedSegments  []PrioritySegment
		records          MediaTransferRecords
		orphaned         []MediaSelector
	}{
		time:             start,
		originalSegments: originalTimeline.FindSegmentsBetween(start, end),
		amendedSegments:  t.timeline.FindSegmentsBetween(start, end),
		records:          make(MediaTransferRecords),
	}

	for len(cursor.originalSegments) > 0 && len(cursor.amendedSegments) > 0 {
		nextTime := minTime(cursor.originalSegments[0].End, cursor.amendedSegments[0].End)
		wasLeading := t.isLeadByAlbum(amendedAlbum.UpdatedAlbum.AlbumId, cursor.originalSegments[0])
		takeTheLead := t.isLeadByAlbum(amendedAlbum.UpdatedAlbum.AlbumId, cursor.amendedSegments[0])

		if wasLeading && !takeTheLead {
			selector := MediaSelector{
				FromAlbums: []AlbumId{amendedAlbum.UpdatedAlbum.AlbumId},
				Start:      cursor.time,
				End:        nextTime,
			}

			if len(cursor.amendedSegments[0].Albums) == 0 {
				cursor.orphaned = append(cursor.orphaned, selector)

			} else {
				target := cursor.amendedSegments[0].Albums[0].AlbumId
				if selectors, found := cursor.records[target]; found {
					cursor.records[target] = append(selectors, selector)
				} else {
					cursor.records[target] = []MediaSelector{selector}
				}
			}
		} else if !wasLeading && takeTheLead && len(cursor.amendedSegments[0].Albums) > 1 {
			selector := MediaSelector{
				FromAlbums: extractAlbumIds(cursor.amendedSegments[0].Albums[1:]), // TODO this is required because there is no index only by dates (no album)
				Start:      cursor.time,
				End:        nextTime,
			}

			target := amendedAlbum.UpdatedAlbum.AlbumId
			if selectors, found := cursor.records[target]; found {
				cursor.records[target] = append(selectors, selector)
			} else {
				cursor.records[target] = []MediaSelector{selector}
			}
		}

		if cursor.originalSegments[0].End.Equal(cursor.amendedSegments[0].End) {
			cursor.time = cursor.originalSegments[0].End
			cursor.originalSegments = cursor.originalSegments[1:]
			cursor.amendedSegments = cursor.amendedSegments[1:]
		} else if cursor.originalSegments[0].End.Before(cursor.amendedSegments[0].End) {
			cursor.time = cursor.originalSegments[0].End
			cursor.originalSegments = cursor.originalSegments[1:]
		} else {
			cursor.time = cursor.amendedSegments[0].End
			cursor.amendedSegments = cursor.amendedSegments[1:]
		}
	}

	return cursor.records, cursor.orphaned, nil
}

func (t *TimelineAggregate) isLeadByAlbum(albumId AlbumId, seg PrioritySegment) bool {
	return len(seg.Albums) > 0 && seg.Albums[0].AlbumId.IsEqual(albumId)
}

func (t *TimelineAggregate) FindAt(date time.Time) (*Album, bool, error) {
	_, err := t.getTimeline()
	if err != nil {
		return nil, false, errors.Wrapf(err, "failed to create timeline during FindAt(%s)", date)
	}

	albumId, found := t.timeline.FindAt(date)
	return albumId, found, nil
}

// getTimeline returns the timeline, creating it if it doesn't exist yet (lazy initialisation)
func (t *TimelineAggregate) getTimeline() (*Timeline, error) {
	var err error
	if t.timeline == nil {
		t.timeline, err = NewTimeline(t.albums)
	}

	return t.timeline, err
}
