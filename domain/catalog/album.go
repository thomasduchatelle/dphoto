// Package catalog provides tools to maintain an index of all medias that have been backed up.
package catalog

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	NotFoundError = errors.New("Album hasn't been found")
	NotEmptyError = errors.New("Album is not empty")
)

// FindAllAlbums find all albums owned by root user
func FindAllAlbums(owner string) ([]*Album, error) {
	return Repository.FindAllAlbums(owner)
}

// FindAllAlbumsWithStats returns the list of albums, with statistics for each
func FindAllAlbumsWithStats(owner string) ([]*AlbumStat, error) {
	albums, err := Repository.FindAllAlbums(owner)
	if err != nil {
		return nil, err
	}

	stats := make([]*AlbumStat, len(albums))
	for i, album := range albums {
		count, err := Repository.CountMedias(owner, album.FolderName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to count medias of %s", album.FolderName)
		}

		stats[i] = &AlbumStat{
			Album:      *album,
			TotalCount: count,
		}
	}

	return stats, err
}

// Create creates a new album
func Create(createRequest CreateAlbum) error {
	if createRequest.Owner == "" {
		return errors.Errorf("Album owner is mandatory")
	}
	if createRequest.Name == "" {
		return errors.Errorf("Album name is mandatory")
	}

	if createRequest.Start.IsZero() || createRequest.End.IsZero() {
		return errors.Errorf("Start and End times are mandatory")
	}

	if !createRequest.End.After(createRequest.Start) {
		return errors.Errorf("Album end must be strictly after its start")
	}

	createdAlbum := Album{
		Owner:      createRequest.Owner,
		Name:       createRequest.Name,
		FolderName: createRequest.ForcedFolderName,
		Start:      createRequest.Start,
		End:        createRequest.End,
	}

	if createdAlbum.FolderName == "" {
		createdAlbum.FolderName = generateAlbumFolder(createRequest.Name, createRequest.Start)
	} else {
		createdAlbum.FolderName = normaliseFolderName(createdAlbum.FolderName)
	}

	albums, err := Repository.FindAllAlbums(createRequest.Owner)
	if err != nil {
		return err
	}

	err = Repository.InsertAlbum(createdAlbum)
	if err != nil {
		return err
	}

	albums = append(albums, &createdAlbum)
	sort.Slice(albums, startsAscSort(albums))

	timeline, err := NewTimeline(albums)
	if err != nil {
		return err
	}

	filter := NewUpdateFilter(createdAlbum.Owner)
	for _, s := range timeline.FindForAlbum(createdAlbum.Owner, createdAlbum.FolderName) {
		filter.WithinRange(s.Start, s.End)
		for _, a := range s.Albums[1:] {
			filter.WithAlbum(a.FolderName)
		}
	}

	medias, count, err := Repository.UpdateMedias(filter, createdAlbum.FolderName)

	log.WithFields(log.Fields{
		"Album":           createdAlbum.FolderName,
		"MoveTransaction": medias,
		"MediaCount":      count,
	}).Infof("Album %s has been created, and %d medias has been virtually moved to it\n", createdAlbum.FolderName, count)

	return err
}

func generateAlbumFolder(name string, start time.Time) string {
	return normaliseFolderName(fmt.Sprintf("%s_%s", start.Format("2006-01"), name))
}

func normaliseFolderName(name string) string {
	nonAlphaNumeric := regexp.MustCompile("[^A-Za-z0-9-]+")
	return "/" + strings.Trim(nonAlphaNumeric.ReplaceAllString(name, "_"), "_")
}

// FindAlbum get an album by its business key (its folder name), or returns NotFoundError
func FindAlbum(owner string, folderName string) (*Album, error) {
	return Repository.FindAlbum(owner, normaliseFolderName(folderName))
}

// DeleteAlbum delete an album, medias it contains are dispatched to other albums.
func DeleteAlbum(owner string, folderNameToDelete string, emptyOnly bool) error {
	folderNameToDelete = normaliseFolderName(folderNameToDelete)
	if !emptyOnly {
		albums, err := Repository.FindAllAlbums(owner)
		if err != nil {
			return err
		}

		albums, removed := removeAlbumFrom(albums, folderNameToDelete)

		if removed != nil {
			_, err = assignMediasTo(owner, albums, removed, func(timeline *Timeline) (segments []PrioritySegment, missed []PrioritySegment, err error) {
				segments, missed = timeline.FindBetween(removed.Start, removed.End)
				reallyMissed, err := filterMissedSegmentWithMedias(owner, folderNameToDelete, missed)
				return segments, reallyMissed, nil
			})
			if err != nil {
				return errors.Wrapf(err, "Album cannot be deleted")
			}
		}
	}

	return Repository.DeleteEmptyAlbum(owner, folderNameToDelete)
}

func filterMissedSegmentWithMedias(owner string, folderName string, missed []PrioritySegment) ([]PrioritySegment, error) {
	var reallyMissed []PrioritySegment
	for _, m := range missed {
		page, err := Repository.FindMedias(owner, normaliseFolderName(folderName), FindMediaFilter{
			PageRequest: PageRequest{Size: 1},
			TimeRange:   TimeRange{Start: m.Start, End: m.End},
		})
		if err != nil {
			return nil, err
		}

		if len(page.Content) > 0 {
			reallyMissed = append(reallyMissed, m)
		}
	}

	return reallyMissed, nil
}

// RenameAlbum updates the displayed named of the album. Optionally changes the folder in which media will be stored
// and flag all its media to be moved to the new one.
func RenameAlbum(owner string, folderName, newName string, renameFolder bool) error {
	folderName = normaliseFolderName(folderName)

	found, err := Repository.FindAlbum(owner, folderName)
	if err != nil {
		return err // can be NotFoundError
	}

	if renameFolder {
		album := Album{
			Owner:      owner,
			Name:       newName,
			FolderName: generateAlbumFolder(newName, found.Start),
			Start:      found.Start,
			End:        found.End,
		}

		err = Repository.InsertAlbum(album)
		if err != nil {
			return err
		}

		transactionId, count, err := Repository.UpdateMedias(NewUpdateFilter(owner).WithAlbum(folderName), album.FolderName)
		if err != nil {
			return err
		}

		log.WithFields(log.Fields{
			"AlbumFolderName":    album.FolderName,
			"PreviousFolderName": folderName,
			"MoveTransactionId":  transactionId,
			"AlbumMoved":         count,
		}).Infof("Album renamed and %d medias moved\n", count)
		return Repository.DeleteEmptyAlbum(owner, folderName)
	}

	found.Name = newName
	return Repository.UpdateAlbum(*found)
}

// UpdateAlbum updates the dates of an album, medias will be re-assign between albums accordingly
func UpdateAlbum(owner string, folderName string, start, end time.Time) error {
	folderName = normaliseFolderName(folderName)

	albums, err := Repository.FindAllAlbums(owner)
	if err != nil {
		return err
	}

	albumsWithoutUpdated, updated := removeAlbumFrom(albums, folderName)
	if updated == nil {
		return NotFoundError
	}

	previousTimeRange := NewTimeRangeFromAlbum(*updated)
	newTimeRange := TimeRange{Start: start, End: end}
	if previousTimeRange.Equals(newTimeRange) {
		log.WithFields(log.Fields{
			"AlbumFolderName": folderName,
			"Start":           start,
			"End":             end,
		}).Infoln("Album date unchanged, nothing to do.")
		return nil
	}

	updated.Start = start
	updated.End = end
	updatedAlbums := append(albumsWithoutUpdated, updated)
	sort.Slice(updatedAlbums, startsAscSort(updatedAlbums))

	count, err := assignMediasTo(owner, updatedAlbums, updated, func(timeline *Timeline) (segments []PrioritySegment, missed []PrioritySegment, err error) {
		for _, timeRange := range previousTimeRange.Plus(newTimeRange) {
			subSegments, subMissed := timeline.FindBetween(timeRange.Start, timeRange.End)
			segments = append(segments, subSegments...)
			missed = append(missed, subMissed...)
		}

		reallyMissed, err := filterMissedSegmentWithMedias(owner, folderName, missed)
		return segments, reallyMissed, nil
	})
	if err != nil {
		return errors.Wrapf(err, "Album dates couldn't be updated.")
	}

	err = Repository.UpdateAlbum(*updated)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"FolderName": folderName,
	}).Infof("Album dates updated, %d medias moved\n", count)
	return nil
}

func assignMediasTo(owner string, albums []*Album, removedAlbum *Album, segmentsToReassignSupplier func(timeline *Timeline) (segments []PrioritySegment, missed []PrioritySegment, err error)) (int, error) {
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
		filter := NewUpdateFilter(owner).WithinRange(s.Start, s.End)
		if removedAlbum != nil && !removedAlbum.IsEqual(&s.Albums[0]) {
			filter.WithAlbum(removedAlbum.FolderName)
		}
		for _, a := range s.Albums[1:] {
			filter.WithAlbum(a.FolderName)
		}

		_, mediaCount, err := Repository.UpdateMedias(filter, s.Albums[0].FolderName)
		if err != nil {
			return 0, err
		}

		count += mediaCount
	}

	return count, nil
}

// remove and keep order
func removeAlbumFrom(albums []*Album, folderName string) ([]*Album, *Album) {
	index := -1
	var removed *Album
	for i, a := range albums {
		if a.FolderName == folderName {
			index = i
			removedCopy := a
			removed = removedCopy
		}
	}

	if index >= 0 {
		return append(albums[:index], albums[index+1:]...), removed
	}

	return albums, removed
}

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

func startsAscComparator(a, b *Album) int64 {
	return b.Start.Unix() - a.Start.Unix()
}

func startsAscSort(albums []*Album) func(i int, j int) bool {
	return func(i, j int) bool {
		return startsAscComparator(albums[i], albums[j]) > 0
	}
}

func albumDuration(album *Album) time.Duration {
	return album.End.Sub(album.Start)
}
