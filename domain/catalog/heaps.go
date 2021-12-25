package catalog

import (
	"container/heap"
	"github.com/thomasduchatelle/dphoto/domain/catalogmodel"
)

type albumHeap struct {
	heap       []*catalogmodel.Album
	comparator func(a, b *catalogmodel.Album) int64
}

func newAlbumHeap(comparator func(a, b *catalogmodel.Album) int64) *albumHeap {
	return &albumHeap{
		comparator: comparator,
	}
}

func (a *albumHeap) Len() int           { return len(a.heap) }
func (a *albumHeap) Less(i, j int) bool { return a.comparator(a.heap[i], a.heap[j]) > 0 }
func (a *albumHeap) Swap(i, j int)      { a.heap[i], a.heap[j] = a.heap[j], a.heap[i] }
func (a *albumHeap) Push(x interface{}) { a.heap = append(a.heap, x.(*catalogmodel.Album)) }

func (a *albumHeap) Pop() interface{} {
	n := len(a.heap)
	x := a.heap[n-1]
	a.heap = a.heap[0 : n-1]
	return x
}

func (a *albumHeap) Head() (*catalogmodel.Album, bool) {
	if a.Len() == 0 {
		return nil, false
	}

	return a.heap[0], true
}

// HeapPush returns: TRUE if new element took the head
func (a *albumHeap) HeapPush(album *catalogmodel.Album) bool {
	heap.Push(a, album)
	return album.IsEqual(a.heap[0])
}

// Remove removes album fro heap and return TRUE if it was the head
func (a *albumHeap) Remove(albumToFind *catalogmodel.Album) bool {
	for index, album := range a.heap {
		if albumToFind.IsEqual(album) {
			heap.Remove(a, index)
			return index == 0
		}
	}

	return false
}

func (a *albumHeap) HasHead(album *catalogmodel.Album) (bool, *catalogmodel.Album) {
	head, notEmpty := a.Head()
	return notEmpty && album.IsEqual(head), head
}

// AsArray copies of the heap: slice where the first element is the head of the heap
func (a *albumHeap) AsArray() []*catalogmodel.Album {
	albums := make([]*catalogmodel.Album, a.Len())
	copy(albums, a.heap)
	return albums
}
