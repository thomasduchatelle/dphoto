package album

import (
	"container/heap"
	"sort"
	"strings"
	"time"
)

type Timeline struct {
	timeline []segment
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

func (c *builderCursor) closeCurrent(end time.Time, timeline *Timeline) {
	if c.current != nil {
		closing := *c.current
		c.current = nil

		closing.to = end
		timeline.timeline = append(timeline.timeline, closing)
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

	if c.current != nil {
		c.current.albums = append(c.current.albums, album)
	}

	return newPriority
}

func NewTimeline(albums []Album) (*Timeline, error) {
	timeline := new(Timeline)

	if len(albums) == 0 {
		return timeline, nil
	}

	cursor := &builderCursor{
		current:         nil,
		nextToCloseHeap: newAlbumHeap(firstToEndComparator),
		priorityHeap:    newAlbumHeap(priorityComparator),
	}

	for index, a := range albums {
		album := a

		for toClose, ok := cursor.nextToCloseHeap.Head(); ok && !toClose.End.After(album.Start); toClose, ok = cursor.nextToCloseHeap.Head() {
			removeAlbumFromCursor(timeline, cursor)
		}

		newPriority := cursor.appendAlbum(&album)

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
	index := sort.Search(len(t.timeline), func(i int) bool {
		return t.timeline[i].to.After(date)
	})


	if index < len(t.timeline) {
		var albums []*Album
		for _, a := range t.timeline[index].albums {
			if !a.Start.After(date) && a.End.After(date) {
				albums = append(albums, a)
			}
		}
		return albums
	}

	return nil
}

// return nil if not found
func (t *Timeline) FindAt(date time.Time) *Album {
	albums := t.FindAllAt(date)
	if len(albums) > 0 {
		return albums[0]
	}

	return nil
}

func priorityComparator(a, b *Album) int64 {
	durationDiff := albumDuration(b).Seconds() - albumDuration(a).Seconds()
	if durationDiff != 0 {
		return int64(durationDiff)
	}

	startDiff := b.Start.Unix() - a.Start.Unix()
	if startDiff != 0 {
		return startDiff
	}

	endDiff := b.End.Unix() - a.End.Unix()
	if endDiff != 0 {
		return endDiff
	}

	return int64(strings.Compare(a.Name, b.Name))
}

func albumDuration(album *Album) time.Duration {
	return album.End.Sub(album.Start)
}
