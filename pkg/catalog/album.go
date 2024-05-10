// Package catalog provides tools to maintain an index of all medias that have been backed up.
package catalog

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sort"
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

func filterMissedSegmentWithMedias(albumId AlbumId, missed []PrioritySegment) ([]PrioritySegment, error) {
	var reallyMissed []PrioritySegment
	for _, m := range missed {
		request := NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName).WithinRange(m.Start, m.End)
		medias, err := repositoryPort.FindMediaIds(context.TODO(), request)

		if err != nil {
			return nil, err
		}

		if len(medias) > 0 {
			reallyMissed = append(reallyMissed, m)
		}
	}

	return reallyMissed, nil
}

// UpdateAlbum updates the dates of an album, medias will be re-assign between albums accordingly
func UpdateAlbum(albumId AlbumId, start, end time.Time) error {
	albums, err := repositoryPort.FindAlbumsByOwner(context.TODO(), albumId.Owner)
	if err != nil {
		return err
	}

	albumsWithoutUpdated, updated := removeAlbumFrom(albums, albumId.FolderName)
	if updated == nil {
		return AlbumNotFoundError
	}

	previousTimeRange := NewTimeRangeFromAlbum(*updated)
	newTimeRange := TimeRange{Start: start, End: end}
	if previousTimeRange.Equals(newTimeRange) {
		log.WithFields(log.Fields{
			"AlbumFolderName": albumId.FolderName,
			"Start":           start,
			"End":             end,
		}).Infoln("Album date unchanged, nothing to do.")
		return nil
	}

	updated.Start = start
	updated.End = end
	updatedAlbums := append(albumsWithoutUpdated, updated)
	sort.Slice(updatedAlbums, startsAscSort(updatedAlbums))

	count, err := assignMediasTo(albumId.Owner, updatedAlbums, updated, func(timeline *Timeline) (segments []PrioritySegment, missed []PrioritySegment, err error) {
		for _, timeRange := range previousTimeRange.Plus(newTimeRange) {
			subSegments, subMissed := timeline.FindBetween(timeRange.Start, timeRange.End)
			segments = append(segments, subSegments...)
			missed = append(missed, subMissed...)
		}

		reallyMissed, err := filterMissedSegmentWithMedias(albumId, missed)
		return segments, reallyMissed, nil
	})
	if err != nil {
		return errors.Wrapf(err, "Album dates couldn't be updated.")
	}

	err = repositoryPort.UpdateAlbum(context.TODO(), *updated)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"FolderName": albumId.FolderName,
	}).Infof("Album dates updated, %d medias moved\n", count)
	return nil
}

func assignMediasTo(owner Owner, albums []*Album, removedAlbum *Album, segmentsToReassignSupplier func(timeline *Timeline) (segments []PrioritySegment, missed []PrioritySegment, err error)) (int, error) {
	timeline, err := NewTimeline(albums)
	if err != nil {
		return 0, err
	}

	segmentsToReassign, missedSegments, err := segmentsToReassignSupplier(timeline)
	if len(missedSegments) > 0 {
		segRanges := make([]string, len(missedSegments))
		for i, seg := range missedSegments {
			segRanges[i] = fmt.Sprintf("%s -> %s", seg.Start.Format("2006-01-02"), seg.End.Format("2006-01-02"))
		}
		return 0, errors.Errorf("some dates are not covered, create albums to cover them before retrying (%s)", strings.Join(segRanges, " ; "))
	}

	if len(segmentsToReassign) == 0 {
		return 0, nil
	}

	count := 0
	for _, s := range segmentsToReassign {
		filter := NewFindMediaRequest(owner).WithinRange(s.Start, s.End)
		if removedAlbum != nil && !removedAlbum.IsEqual(&s.Albums[0]) {
			filter.WithAlbum(removedAlbum.FolderName)
		}
		for _, a := range s.Albums[1:] {
			filter.WithAlbum(a.FolderName)
		}

		mediaCount, err := transferMedias(filter, s.Albums[0].FolderName)
		if err != nil {
			return 0, err
		}

		count += mediaCount
	}

	return count, nil
}

// remove and keep order
func removeAlbumFrom(albums []*Album, folderName FolderName) ([]*Album, *Album) {
	for index, album := range albums {
		if album.FolderName == folderName {
			return append(albums[:index], albums[index+1:]...), album
		}
	}

	return albums, nil
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
