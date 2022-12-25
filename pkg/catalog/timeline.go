package catalog

import (
	"container/heap"
	"github.com/pkg/errors"
	"sort"
	"time"
)

// Timeline can be used to find to which album a media will belongs.
type Timeline struct {
	segments []segment
	albums   []*Album // albums is only used to re-generate a new Timeline
}

type segment struct {
	from   time.Time
	to     time.Time
	albums []*Album
}

type builderCursor struct {
	current         *segment
	nextToCloseHeap *albumHeap
	priorityHeap    *albumHeap
}

type PrioritySegment struct {
	Start  time.Time
	End    time.Time
	Albums []Album // sorted by priority
}

func (c *builderCursor) closeCurrent(end time.Time, timeline *Timeline) {
	if c.current != nil {
		closing := *c.current
		c.current = nil

		closing.to = end
		timeline.segments = append(timeline.segments, closing)
	}
}

func (c *builderCursor) removeNextToClose() (*Album, bool) {
	if c.nextToCloseHeap.Len() == 0 {
		return nil, false
	}

	toClose := heap.Pop(c.nextToCloseHeap).(*Album)
	hasPriority := c.priorityHeap.Remove(toClose)

	return toClose, hasPriority
}

func (c *builderCursor) startSegment(start time.Time) {
	c.current = &segment{
		from:   start,
		albums: c.priorityHeap.AsArray(),
	}
}

// add album to both heaps, return TRUE if album is the new priority
func (c *builderCursor) appendAlbum(album *Album) bool {
	c.nextToCloseHeap.HeapPush(album)
	newPriority := c.priorityHeap.HeapPush(album)

	if !newPriority && c.current != nil {
		c.current.albums = append(c.current.albums, album)
	}

	return newPriority
}

// NewTimeline creates a Timeline object used to compute overlaps between Album. List of albums must be sorted by Start date ASC (End sorting do not matter).
func NewTimeline(albums []*Album) (*Timeline, error) {
	timeline := &Timeline{
		albums: albums,
	}

	if len(albums) == 0 {
		return timeline, nil
	}

	cursor := &builderCursor{
		current:         nil,
		nextToCloseHeap: newAlbumHeap(endDescComparator),
		priorityHeap:    newAlbumHeap(priorityDescComparator),
	}

	for index, a := range albums {
		album := a
		if index > 0 && albums[index-1].Start.After(album.Start) {
			return nil, errors.Errorf("Albums must be sorted by Start date ASC, %s [index %d] is before %s", album.String(), index, albums[index-1].String())
		}
		if !album.End.After(album.Start) {
			return nil, errors.Errorf("Invalid album, end date must be after start date: %s", album.String())
		}

		for toClose, ok := cursor.nextToCloseHeap.Head(); ok && !toClose.End.After(album.Start); toClose, ok = cursor.nextToCloseHeap.Head() {
			removeAlbumFromCursor(timeline, cursor)
		}

		newPriority := cursor.appendAlbum(album)

		if newPriority {
			cursor.closeCurrent(album.Start, timeline)
		}

		hasNext := index+1 < len(albums)
		if cursor.current == nil && (!hasNext || !albums[index+1].Start.Equal(a.Start)) {
			cursor.startSegment(album.Start)
		}
	}

	for cursor.nextToCloseHeap.Len() > 0 {
		removeAlbumFromCursor(timeline, cursor)
	}

	return timeline, nil
}

func removeAlbumFromCursor(timeline *Timeline, cursor *builderCursor) {
	toClose, hadPriority := cursor.removeNextToClose()

	if cursor.current != nil && hadPriority {
		cursor.closeCurrent(toClose.End, timeline)
	}

	// skipped if next to close has same end date
	if nextToClose, hasNext := cursor.nextToCloseHeap.Head(); cursor.current == nil && hasNext && !nextToClose.End.Equal(toClose.End) {
		cursor.startSegment(toClose.End)
	}
}

func (t *Timeline) FindAllAt(date time.Time) []*Album {
	index := sort.Search(len(t.segments), func(i int) bool {
		return t.segments[i].to.After(date)
	})

	if index < len(t.segments) {
		var albums []*Album
		for _, a := range t.segments[index].albums {
			if !a.Start.After(date) && a.End.After(date) {
				albums = append(albums, a)
			}
		}
		return albums
	}

	return nil
}

// FindAt returns nil if not found
func (t *Timeline) FindAt(date time.Time) (*Album, bool) {
	albums := t.FindAllAt(date)
	if len(albums) > 0 {
		return albums[0], true
	}

	return nil, false
}

func (t *Timeline) FindForAlbum(owner, folderName string) (segments []PrioritySegment) {
	for _, seg := range t.segments {
		if seg.albums[0].Owner == owner && seg.albums[0].FolderName == folderName {
			segments = append(segments, PrioritySegment{
				Start:  seg.from,
				End:    seg.to,
				Albums: toSortedArray(seg.albums, priorityDescComparator),
			})
		}
	}

	return segments
}

func (t *Timeline) FindBetween(start, end time.Time) (segments []PrioritySegment, missed []PrioritySegment) {
	startIndex := sort.Search(len(t.segments), func(i int) bool {
		return t.segments[i].to.After(start)
	})

	if startIndex >= len(t.segments) {
		return
	}

	endIndex := sort.Search(len(t.segments)-startIndex, func(i int) bool {
		return !t.segments[startIndex+i].from.Before(end)
	})

	previousEnd := start
	for _, seg := range t.segments[startIndex : startIndex+endIndex] {
		if previousEnd.Before(seg.from) {
			missed = append(missed, PrioritySegment{
				Start: previousEnd,
				End:   seg.from,
			})
		}
		previousEnd = seg.to
		segments = append(segments, PrioritySegment{
			Start:  maxTime(seg.from, start),
			End:    minTime(seg.to, end),
			Albums: toSortedArray(seg.albums, priorityDescComparator),
		})
	}

	if len(segments) == 0 {
		missed = append(missed, PrioritySegment{
			Start: start,
			End:   end,
		})
	} else if segments[len(segments)-1].End.Before(end) {
		missed = append(missed, PrioritySegment{
			Start: segments[len(segments)-1].End,
			End:   end,
		})
	}

	return segments, missed
}

// AppendAlbum generates a new timeline from memory
func (t *Timeline) AppendAlbum(album *Album) (*Timeline, error) {
	albums := make([]*Album, len(t.albums)+1)
	albums[0] = album
	copy(albums[1:], t.albums)

	sort.Slice(albums, func(i, j int) bool {
		return startsAscComparator(albums[i], albums[j]) > 0
	})

	return NewTimeline(albums)
}

func toSortedArray(albums []*Album, comparator func(a *Album, b *Album) int64) []Album {
	sortedAlbums := make([]Album, len(albums))
	for i, a := range albums {
		sortedAlbums[i] = *a
	}

	sort.Slice(sortedAlbums, func(i, j int) bool {
		return comparator(&sortedAlbums[i], &sortedAlbums[j]) > 0
	})

	return sortedAlbums
}
