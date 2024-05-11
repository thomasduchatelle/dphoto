// Package catalog provides tools to maintain an index of all medias that have been backed up.
package catalog

import (
	"context"
	"strings"
	"time"
)

// FindAllAlbums find all albums owned by root user
func FindAllAlbums(owner Owner) ([]*Album, error) {
	return repositoryPort.FindAlbumsByOwner(context.TODO(), owner)
}

// FindAlbums get several albums by their business keys
func FindAlbums(keys []AlbumId) ([]*Album, error) {
	return repositoryPort.FindAlbumByIds(context.TODO(), keys...)
}

// FindAlbum get an album by its business key (its folder name), or returns AlbumNotFoundError
func FindAlbum(id AlbumId) (*Album, error) {
	albums, err := repositoryPort.FindAlbumByIds(context.TODO(), id)
	if err != nil {
		return nil, err
	}
	if len(albums) == 0 {
		return nil, AlbumNotFoundError
	}
	return albums[0], nil
}

// priorityDescComparator is positive if a is more important than b
func priorityDescComparator(a, b *Album) int64 {
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

func endDescComparator(a, b *Album) int64 {
	return b.End.Unix() - a.End.Unix()
}

// startsAscSort sorts albums by start date ascending, then by priority descending (equivalent to end date ascending)
func startsAscComparator(a, b *Album) int64 {
	if a.Start.Equal(b.Start) {
		return priorityDescComparator(a, b)
	}
	return b.Start.Unix() - a.Start.Unix()
}

func albumDuration(album *Album) time.Duration {
	return album.End.Sub(album.Start)
}
