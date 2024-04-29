package catalog

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

// TimelineMutator is used to measure the impact of a change on the timeline
type TimelineMutator struct {
	Current []*Album
}

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
	FromAlbums []AlbumId // FromAlbums is a list
	Start      time.Time // Start is the first date of matching medias, included
	End        time.Time // End is the last date of matching media, excluded at the second
}

func (m MediaSelector) String() string {
	return fmt.Sprintf("%s -> %s", m.Start.Format(time.DateTime), m.End.Format(time.DateTime))
}

func NewTimelineMutator(current []*Album) *TimelineMutator {
	return &TimelineMutator{
		Current: current,
	}
}

func (t *TimelineMutator) AddNew(album Album) (MediaTransferRecords, error) {
	albums := append(t.Current, &album)
	sort.Slice(albums, startsAscSort(albums))

	timeline, err := NewTimeline(albums)
	if err != nil {
		return nil, err
	}

	records := make(MediaTransferRecords)
	for _, seg := range timeline.FindForAlbum(album.AlbumId) {
		if len(seg.Albums) > 1 {
			selector := MediaSelector{
				FromAlbums: extractAlbumIds(seg.Albums[1:]),
				Start:      seg.Start,
				End:        seg.End,
			}
			if selectors, found := records[album.AlbumId]; found {
				records[album.AlbumId] = append(selectors, selector)
			} else {
				records[album.AlbumId] = []MediaSelector{selector}
			}
		}
	}

	if len(records) == 0 {
		return nil, nil
	}

	return records, nil
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
