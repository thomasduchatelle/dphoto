package catalog

import (
	"container/heap"
	"slices"
)

type albumHeap struct {
	heap       []*Album
	comparator func(a, b *Album) int64
}

func newAlbumHeap(comparator func(a, b *Album) int64) *albumHeap {
	return &albumHeap{
		comparator: comparator,
	}
}

func (a *albumHeap) Len() int           { return len(a.heap) }
func (a *albumHeap) Less(i, j int) bool { return a.comparator(a.heap[i], a.heap[j]) > 0 }
func (a *albumHeap) Swap(i, j int)      { a.heap[i], a.heap[j] = a.heap[j], a.heap[i] }
func (a *albumHeap) Push(x interface{}) { a.heap = append(a.heap, x.(*Album)) }

// Pop removes the last element of the heap
func (a *albumHeap) Pop() interface{} {
	n := len(a.heap)
	x := a.heap[n-1]
	a.heap = a.heap[0 : n-1]
	return x
}

// Head returns the head of the heap
func (a *albumHeap) Head() (*Album, bool) {
	if a.Len() == 0 {
		return nil, false
	}

	return a.heap[0], true
}

// HeapPush adds album to the heap and return TRUE if it is the head
func (a *albumHeap) HeapPush(album *Album) bool {
	heap.Push(a, album)
	return album.IsEqual(a.heap[0])
}

// Remove removes album from the heap and return TRUE if it was the head
func (a *albumHeap) Remove(albumToFind *Album) bool {
	for index, album := range a.heap {
		if albumToFind.IsEqual(album) {
			heap.Remove(a, index)
			return index == 0
		}
	}

	return false
}

// HasHead returns TRUE if the head of the heap is the same as album
func (a *albumHeap) HasHead(album *Album) (bool, *Album) {
	head, notEmpty := a.Head()
	return notEmpty && album.IsEqual(head), head
}

// AsArray returns the heap as an array
func (a *albumHeap) AsArray() []*Album {
	// TODO Do better using the heap !
	albums := make([]*Album, a.Len())
	copy(albums, a.heap)
	slices.SortFunc(albums, func(a, b *Album) int {
		return -int(priorityDescComparator(a, b))
	})
	return albums
}

// RemoveHead removes the head of the heap
func (a *albumHeap) RemoveHead() *Album {
	if a.Len() == 0 {
		return nil
	}

	return heap.Remove(a, 0).(*Album)
}
