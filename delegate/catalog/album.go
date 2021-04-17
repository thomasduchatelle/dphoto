// Package provides tools to maintain a catalog of medias and albums owned by a user
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

// FindAllAlbum FindAlbum all albums owned by root user
func FindAllAlbum() ([]*Album, error) {
	return Repository.FindAllAlbums()
}

// Create creates a new album
func Create(createRequest CreateAlbum) error {
	if createRequest.Name == "" {
		return errors.Errorf("Album name is mandatory")
	}

	if createRequest.Start == nil || createRequest.End == nil {
		return errors.Errorf("Start and End times are mandatory")
	}

	if !createRequest.End.After(*createRequest.Start) {
		return errors.Errorf("Albem end must be strictly after its start")
	}

	createdAlbum := Album{
		Name:       createRequest.Name,
		FolderName: createRequest.ForcedFolderName,
		Start:      *createRequest.Start,
		End:        *createRequest.End,
	}

	if createdAlbum.FolderName == "" {
		createdAlbum.FolderName = generateAlbumFolder(createRequest.Name, *createRequest.Start)
	}

	albums, err := Repository.FindAllAlbums()
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

	filter := NewUpdateFilter()
	for _, s := range timeline.FindForAlbum(createdAlbum.FolderName) {
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
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	return strings.Trim(fmt.Sprintf("%s_%s", start.Format("2006-01"), re.ReplaceAllString(name, "_")), "_")
}

// FindAlbum get an album by its business key (its folder name), or returns NotFoundError
func FindAlbum(folderName string) (*Album, error) {
	return Repository.FindAlbum(folderName)
}

// DeleteAlbum delete an album, medias it contains are dispatched to other albums.
func DeleteAlbum(folderNameToDelete string, emptyOnly bool) error {
	if !emptyOnly {
		albums, err := Repository.FindAllAlbums()
		if err != nil {
			return err
		}

		albums, removed := removeAlbumFrom(albums, folderNameToDelete)

		if removed != nil {
			_, err = assignMediasTo(albums, removed, func(timeline *Timeline) []PrioritySegment {
				return timeline.FindBetween(removed.Start, removed.End)
			})
			if err != nil {
				return err
			}
		}
	}

	return Repository.DeleteEmptyAlbum(folderNameToDelete)
}

// RenameAlbum updates the displayed named of the album. Optionally changes the folder in which media will be stored
// and flag all its media to be moved to the new one.
func RenameAlbum(folderName, newName string, renameFolder bool) error {
	found, err := Repository.FindAlbum(folderName)
	if err != nil {
		return err // can be NotFoundError
	}

	if renameFolder {
		album := Album{
			Name:       newName,
			FolderName: generateAlbumFolder(newName, found.Start),
			Start:      found.Start,
			End:        found.End,
		}

		err = Repository.InsertAlbum(album)
		if err != nil {
			return err
		}

		transactionId, count, err := Repository.UpdateMedias(NewUpdateFilter().WithAlbum(folderName), album.FolderName)
		if err != nil {
			return err
		}

		log.WithFields(log.Fields{
			"AlbumFolderName":    album.FolderName,
			"PreviousFolderName": folderName,
			"MoveTransactionId":  transactionId,
			"AlbumMoved":         count,
		}).Infof("Album renamed and %d medias moved\n", count)
		return Repository.DeleteEmptyAlbum(folderName)
	}

	found.Name = newName
	return Repository.UpdateAlbum(*found)
}

// UpdateAlbum updates the dates of an album, medias will be re-assign between albums accordingly
func UpdateAlbum(folderName string, start, end time.Time) error {
	albums, err := Repository.FindAllAlbums()
	if err != nil {
		return err
	}

	albumsWithoutUpdated, updated := removeAlbumFrom(albums, folderName)
	if updated == nil {
		return NotFoundError
	}

	previousTimeRange := newTimeRangeFromAlbum(*updated)
	newTimeRange := timeRange{Start: start, End: end}
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

	count, err := assignMediasTo(updatedAlbums, updated, func(timeline *Timeline) (segments []PrioritySegment) {
		for _, timeRange := range previousTimeRange.Plus(newTimeRange) {
			segments = append(segments, timeline.FindBetween(timeRange.Start, timeRange.End)...)
		}

		return segments
	})
	if err != nil {
		log.WithFields(log.Fields{
			"FolderName": folderName,
		}).Infof("Album dates updated, %d medias moved\n", count)
	}
	return err
}

func assignMediasTo(albums []*Album, removedAlbum *Album, segmentsToReassignSupplier func(timeline *Timeline) []PrioritySegment) (int, error) {
	timeline, err := NewTimeline(albums)
	if err != nil {
		return 0, err
	}

	segmentsToReassign := segmentsToReassignSupplier(timeline)
	if len(segmentsToReassign) == 0 {
		return 0, nil
	}

	count := 0
	for _, s := range segmentsToReassign {
		filter := NewUpdateFilter().WithinRange(s.Start, s.End)
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
