package catalog

import (
	"github.com/pkg/errors"
	"slices"
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
	start        time.Time
	priorityHeap *albumHeap
}

type PrioritySegment struct {
	Start  time.Time
	End    time.Time
	Albums []Album // sorted by priority
}

func (c *builderCursor) closeCurrent(end time.Time, timeline *Timeline) {
	if !c.start.IsZero() {
		var albums []*Album
		for _, album := range c.priorityHeap.heap {
			if album.End.After(c.start) && album.Start.Before(end) {
				albums = append(albums, album)
			}
		}
		slices.SortFunc(albums, func(a, b *Album) int {
			return -int(priorityDescComparator(a, b))
		})
		timeline.segments = append(timeline.segments, segment{
			from:   c.start,
			to:     end,
			albums: albums,
		})

		c.start = time.Time{}
	}
}

func (c *builderCursor) startSegment(start time.Time) {
	c.start = start
}

// add album to both heaps, return TRUE if album is the new priority
func (c *builderCursor) appendAlbum(album *Album) bool {
	return c.priorityHeap.HeapPush(album)
}

// NewTimeline creates a Timeline object used to compute overlaps between Album. List of albums must be sorted by Start date ASC (End sorting does not matter).
func NewTimeline(albums []*Album) (*Timeline, error) {
	slices.SortFunc(albums, func(a, b *Album) int {
		return -int(startsAscComparator(a, b))
	})

	timeline := &Timeline{
		albums: albums,
	}

	if len(albums) == 0 {
		return timeline, nil
	}

	cursor := &builderCursor{
		priorityHeap: newAlbumHeap(priorityDescComparator),
	}

	for index, next := range albums {
		if index > 0 && albums[index-1].Start.After(next.Start) {
			return nil, errors.Errorf("Albums must be sorted by Start date ASC, %s [index %d] is before %s", next.String(), index, albums[index-1].String())
		}
		if !next.End.After(next.Start) {
			return nil, errors.Errorf("Invalid album, end date must be after start date: %s", next.String())
		}

		for lead, hasLead := cursor.priorityHeap.Head(); hasLead && lead.End.Before(next.Start); lead, hasLead = cursor.priorityHeap.Head() {
			electNewLeader(cursor, lead.End, timeline)
		}

		takesTheLead := cursor.appendAlbum(next)
		if takesTheLead {
			electNewLeader(cursor, next.Start, timeline)
		}
	}

	for lead, hasLead := cursor.priorityHeap.Head(); hasLead; lead, hasLead = cursor.priorityHeap.Head() {
		electNewLeader(cursor, lead.End, timeline)
	}

	return timeline, nil
}

func electNewLeader(cursor *builderCursor, atTheTime time.Time, timeline *Timeline) {
	cursor.closeCurrent(atTheTime, timeline)

	for rottenHead, notEmpty := cursor.priorityHeap.Head(); notEmpty && !rottenHead.End.After(atTheTime); rottenHead, notEmpty = cursor.priorityHeap.Head() {
		cursor.priorityHeap.RemoveHead()
	}

	if _, hasNewLead := cursor.priorityHeap.Head(); hasNewLead {
		cursor.startSegment(atTheTime)
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

func (t *Timeline) FindForAlbum(albumId AlbumId) (segments []PrioritySegment) {
	for _, seg := range t.segments {
		if seg.albums[0].AlbumId.IsEqual(albumId) {
			segments = append(segments, PrioritySegment{
				Start:  seg.from,
				End:    seg.to,
				Albums: toSortedArray(seg.albums, priorityDescComparator),
			})
		}
	}

	return segments
}

// FindBetween is deprecated, use FindSegmentsBetween instead
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

// FindSegmentsBetween returns a list of segments between start and end date. Segments will cover the whole period, but might not have any album.
func (t *Timeline) FindSegmentsBetween(start, end time.Time) (segments []PrioritySegment) {
	startIndex := sort.Search(len(t.segments), func(i int) bool {
		return t.segments[i].to.After(start)
	})

	if startIndex >= len(t.segments) {
		return []PrioritySegment{
			{Start: start, End: end},
		}
	}

	endIndex := sort.Search(len(t.segments)-startIndex, func(i int) bool {
		return !t.segments[startIndex+i].from.Before(end)
	})

	previousEnd := start
	for _, seg := range t.segments[startIndex : startIndex+endIndex] {
		if previousEnd.Before(seg.from) {
			segments = append(segments, PrioritySegment{
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
		segments = append(segments, PrioritySegment{
			Start: start,
			End:   end,
		})
	} else if segments[len(segments)-1].End.Before(end) {
		segments = append(segments, PrioritySegment{
			Start: segments[len(segments)-1].End,
			End:   end,
		})
	}

	return segments
}

// FindSegmentsBetweenAndFilter returns a list of segments between start and end date, only segments lead by the given albumId will be returned.
func (t *Timeline) FindSegmentsBetweenAndFilter(start, end time.Time, albumId AlbumId) (segments []PrioritySegment) {
	for _, seg := range t.FindSegmentsBetween(start, end) {
		if seg.Albums[0].AlbumId.IsEqual(albumId) {
			segments = append(segments, seg)
		}
	}
	return
}

// AppendAlbum generates a new timeline from memory
func (t *Timeline) AppendAlbum(album *Album) (*Timeline, error) {
	albums := make([]*Album, len(t.albums)+1)
	albums[0] = album
	copy(albums[1:], t.albums)

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
