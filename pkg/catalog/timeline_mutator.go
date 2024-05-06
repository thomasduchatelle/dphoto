package catalog

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

// TimelineMutator is used to measure the impact of a change on the timeline
type TimelineMutator struct{}

// MediaTransferRecords is a description of all medias that needs to be moved accordingly to the Timeline change
type MediaTransferRecords map[AlbumId][]MediaSelector

func (r MediaTransferRecords) String() string {
	if len(r) == 0 {
		return "<no media transfer>"
	}

	var transfer []string
	for albumId, selector := range r {
		transfer = append(transfer, fmt.Sprintf("%s [%s]", albumId, selector))
	}
	return strings.Join(transfer, ", ")
}

type TimelineMutationObserver interface {
	Observe(ctx context.Context, transfers TransferredMedias) error
}

type MediaSelector struct {
	//ExclusiveAlbum *AlbumId  // ExclusiveAlbum is the Album in which medias are NOT (optional)
	FromAlbums []AlbumId // FromAlbums is a list of potential origins of medias ; is mandatory on CreateAlbum case because media are not indexed by date, only per album.
	Start      time.Time // Start is the first date of matching medias, included
	End        time.Time // End is the last date of matching media, excluded at the second
}

func (m MediaSelector) String() string {
	return fmt.Sprintf("%s -> %s", m.Start.Format(time.DateTime), m.End.Format(time.DateTime))
}

func NewTimelineMutator() *TimelineMutator {
	return new(TimelineMutator)
}

func (t TimelineMutator) AddNew(currentAlbums []*Album, addedAlbum Album) (MediaTransferRecords, error) {
	albums := append(currentAlbums, &addedAlbum)

	timeline, err := NewTimeline(albums)
	if err != nil {
		return nil, err
	}

	records := make(MediaTransferRecords)
	for _, seg := range timeline.FindForAlbum(addedAlbum.AlbumId) {
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

func (t TimelineMutator) RemoveAlbum(currentAlbums []*Album, deletedAlbumId AlbumId) (MediaTransferRecords, []MediaSelector, error) {
	records := make(MediaTransferRecords)
	orphaned := make([]MediaSelector, 0)

	albumsWithoutOneToDelete, deletedAlbum := removeAlbumFrom(currentAlbums, deletedAlbumId.FolderName)
	if deletedAlbum == nil {
		return nil, nil, NotFoundError
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
