package catalog

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/catalogmodel"
	"testing"
	"time"
)

func TestHeap(t *testing.T) {
	a := assert.New(t)

	h := &albumHeap{
		comparator: endDescComparator,
	}
	heap.Init(h)

	for _, album := range AlbumCollection() {
		a := album
		heap.Push(h, a)
	}

	a.Equal("/2020-Q3", heap.Pop(h).(*catalogmodel.Album).FolderName)
	a.Equal("/Christmas_First_Week", heap.Pop(h).(*catalogmodel.Album).FolderName)
	a.Equal("/Christmas_Day", heap.Pop(h).(*catalogmodel.Album).FolderName)
	a.Equal("/2020-Q4", heap.Pop(h).(*catalogmodel.Album).FolderName)
	a.Equal("/New_Year", heap.Pop(h).(*catalogmodel.Album).FolderName)
	a.Equal("/Christmas_Holidays", heap.Pop(h).(*catalogmodel.Album).FolderName)
	a.Equal("/2021-Q1", heap.Pop(h).(*catalogmodel.Album).FolderName)
	a.Equal("/2021-May", heap.Pop(h).(*catalogmodel.Album).FolderName)
	a.Empty(h.heap)
}

func AlbumCollection() []*catalogmodel.Album {
	return []*catalogmodel.Album{
		{
			FolderName: "/2020-Q3",
			Start:      time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
			End:        time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			FolderName: "/2020-Q4",
			Start:      time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
			End:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			FolderName: "/Christmas_Holidays",
			Start:      time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC),
			End:        time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			FolderName: "/Christmas_First_Week",
			Start:      time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC),
			End:        time.Date(2020, 12, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			FolderName: "/Christmas_Day",
			Start:      time.Date(2020, 12, 24, 0, 0, 0, 0, time.UTC),
			End:        time.Date(2020, 12, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			FolderName: "/New_Year",
			Start:      time.Date(2020, 12, 31, 18, 0, 0, 0, time.UTC),
			End:        time.Date(2021, 1, 1, 18, 0, 0, 0, time.UTC),
		},
		{
			FolderName: "/2021-Q1",
			Start:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			End:        time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			FolderName: "/2021-May",
			Start:      time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC),
			End:        time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC),
		},
	}
}
