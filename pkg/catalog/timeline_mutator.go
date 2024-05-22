package catalog

import (
	"github.com/pkg/errors"
	"sort"
	"time"
)

// TimelineMutator is used to measure the impact of a change on the timeline
type TimelineMutator struct{}

func NewTimelineMutator() *TimelineMutator {
	return new(TimelineMutator)
}

func (t TimelineMutator) RemoveAlbum(currentAlbums []*Album, deletedAlbumId AlbumId) (MediaTransferRecords, []MediaSelector, error) {
	records := make(MediaTransferRecords)
	orphaned := make([]MediaSelector, 0)

	albumsWithoutOneToDelete, deletedAlbum := t.removeAlbumFrom(currentAlbums, deletedAlbumId.FolderName)
	if deletedAlbum == nil {
		return nil, nil, AlbumNotFoundError
	}
	sort.Slice(albumsWithoutOneToDelete, startsAscSort(albumsWithoutOneToDelete))

	timeline, err := NewTimeline(albumsWithoutOneToDelete)
	if err != nil {
		return nil, nil, err
	}

	segments := timeline.FindSegmentsBetween(deletedAlbum.Start, deletedAlbum.End)
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
func (t TimelineMutator) removeAlbumFrom(albums []*Album, folderName FolderName) ([]*Album, *Album) {
	for index, album := range albums {
		if album.FolderName == folderName {
			return append(albums[:index], albums[index+1:]...), album
		}
	}

	return albums, nil
}

func (t TimelineMutator) AmendDates(timeline []*Album, amendedAlbum Album) (MediaTransferRecords, []MediaSelector, error) {
	originalTimeline, err := NewTimeline(timeline)
	if err != nil {
		return nil, nil, err
	}

	updatedList, previous, err := t.copyListWithAmendedAlbum(timeline, amendedAlbum)
	amendedTimeline, err := NewTimeline(updatedList)
	if err != nil {
		return nil, nil, err
	}

	start := minTime(amendedAlbum.Start, previous.Start)
	end := maxTime(amendedAlbum.End, previous.End)

	cursor := struct {
		time             time.Time
		originalSegments []PrioritySegment
		amendedSegments  []PrioritySegment
		records          MediaTransferRecords
		orphaned         []MediaSelector
	}{
		time:             start,
		originalSegments: originalTimeline.FindSegmentsBetween(start, end),
		amendedSegments:  amendedTimeline.FindSegmentsBetween(start, end),
		records:          make(MediaTransferRecords),
	}

	for len(cursor.originalSegments) > 0 && len(cursor.amendedSegments) > 0 {
		nextTime := minTime(cursor.originalSegments[0].End, cursor.amendedSegments[0].End)
		wasLeading := t.isLeadByAlbum(amendedAlbum.AlbumId, cursor.originalSegments[0])
		takeTheLead := t.isLeadByAlbum(amendedAlbum.AlbumId, cursor.amendedSegments[0])

		if wasLeading && !takeTheLead {
			selector := MediaSelector{
				FromAlbums: []AlbumId{amendedAlbum.AlbumId},
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

			target := amendedAlbum.AlbumId
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

func (t TimelineMutator) isLeadByAlbum(albumId AlbumId, seg PrioritySegment) bool {
	return len(seg.Albums) > 0 && seg.Albums[0].AlbumId.IsEqual(albumId)
}

func (t TimelineMutator) copyListWithAmendedAlbum(timeline []*Album, amendedAlbum Album) ([]*Album, Album, error) {
	var previous *Album
	amendedTimeline := make([]*Album, len(timeline), len(timeline))

	for i, album := range timeline {
		if album.AlbumId.IsEqual(amendedAlbum.AlbumId) {
			previous = album
			amendedTimeline[i] = &amendedAlbum
		} else {
			amendedTimeline[i] = album
		}
	}

	if previous == nil {
		return nil, Album{}, errors.Errorf("album %s not found in timeline %+v", amendedAlbum.AlbumId.FolderName, timeline)
	}

	return amendedTimeline, *previous, nil
}

func startsAscSort(albums []*Album) func(i int, j int) bool {
	return func(i, j int) bool {
		return startsAscComparator(albums[i], albums[j]) > 0
	}
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
