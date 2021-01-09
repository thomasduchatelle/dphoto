package album

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
	NotFoundError = errors.New("Album not found")
)

func FindAll() ([]Album, error) {
	return Repository.FindAll()
}

func Create(createrequest CreateAlbum) error {
	if createrequest.Name == "" {
		return errors.Errorf("Album name is mandatory")
	}

	if createrequest.Start == nil || createrequest.End == nil {
		return errors.Errorf("Start and End times are mandatory")
	}

	if !createrequest.End.After(*createrequest.Start) {
		return errors.Errorf("Albem end must be strictly after its start")
	}

	createdAlbum := Album{
		Name:       createrequest.Name,
		FolderName: createrequest.ForcedFolderName,
		Start:      *createrequest.Start,
		End:        *createrequest.End,
	}

	if createdAlbum.FolderName == "" {
		createdAlbum.FolderName = generateAlbumFolder(createrequest.Name, *createrequest.Start)
	}

	albums, err := Repository.FindAll()
	if err != nil {
		return err
	}

	err = Repository.Insert(createdAlbum)
	if err != nil {
		return err
	}

	albums = append(albums, createdAlbum)
	sort.Slice(albums, startsAscSort(albums))

	timeline, err := NewTimeline(albums)
	if err != nil {
		return err
	}

	filter := NewFilter()
	for _, s := range timeline.FindForAlbum(createdAlbum.FolderName) {
		filter.WithinRange(s.Start, s.End)
		for _, a := range s.Albums[1:] {
			filter.WithAlbum(a.FolderName)
		}
	}

	return Repository.UpdateMedias(filter, NewMoveMediaUpdate(createdAlbum))
}

func generateAlbumFolder(name string, start time.Time) string {
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	return strings.Trim(fmt.Sprintf("%s_%s", start.Format("2006-01"), re.ReplaceAllString(name, "_")), "_")
}

func Find(folderName string) (*Album, error) {
	return Repository.Find(folderName)
}

func Delete(folderNameToDelete string, emptyOnly bool) error {
	if !emptyOnly {
		albums, err := Repository.FindAll()
		if err != nil {
			return err
		}

		albums, removed := removeAlbumFrom(albums, folderNameToDelete)

		if removed != nil {
			err = assignMediasTo(albums, removed, func(timeline *Timeline) []PrioritySegment {
				return timeline.FindBetween(removed.Start, removed.End)
			})
			if err != nil {
				return err
			}
		}
	}

	return Repository.DeleteEmpty(folderNameToDelete)
}

// Rename an album, and flag all medias to be moved...
// folderName: optional, force to use a specific name
func Rename(folderName, newName string, renameFolder bool) error {
	found, err := Repository.Find(folderName)
	if err != nil {
		return err
	}
	if found == nil {
		return NotFoundError
	}

	if renameFolder {
		album := Album{
			Name:       newName,
			FolderName: generateAlbumFolder(newName, found.Start),
			Start:      found.Start,
			End:        found.End,
		}

		err = Repository.Insert(album)
		if err != nil {
			return err
		}

		err = Repository.UpdateMedias(NewFilter().WithAlbum(folderName), NewMoveMediaUpdate(album))
		if err != nil {
			return err
		}

		return Repository.DeleteEmpty(folderName)
	}

	found.Name = newName
	return Repository.Update(*found)
}

func Update(folderName string, start, end time.Time) error {
	albums, err := Repository.FindAll()
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
			"FolderName": folderName,
			"Start":      start,
			"End":        end,
		}).Infoln("Album date unchanged, nothing to do.")
		return nil
	}

	updated.Start = start
	updated.End = end
	updatedAlbums := append(albumsWithoutUpdated, *updated)
	sort.Slice(updatedAlbums, startsAscSort(updatedAlbums))

	return assignMediasTo(updatedAlbums, updated, func(timeline *Timeline) (segments []PrioritySegment) {
		for _, timeRange := range previousTimeRange.Plus(newTimeRange) {
			segments = append(segments, timeline.FindBetween(timeRange.Start, timeRange.End)...)
		}

		return segments
	})
}

func assignMediasTo(albums []Album, removedAlbum *Album, segmentsToReassignSupplier func(timeline *Timeline) []PrioritySegment) error {
	timeline, err := NewTimeline(albums)
	if err != nil {
		return err
	}

	segmentsToReassign := segmentsToReassignSupplier(timeline)
	if len(segmentsToReassign) == 0 {
		return nil
	}

	for _, s := range segmentsToReassign {
		filter := NewFilter().WithinRange(s.Start, s.End)
		if removedAlbum != nil && !removedAlbum.IsEqual(&s.Albums[0]) {
			filter.WithAlbum(removedAlbum.FolderName)
		}
		for _, a := range s.Albums[1:] {
			filter.WithAlbum(a.FolderName)
		}

		err = Repository.UpdateMedias(filter, NewMoveMediaUpdate(s.Albums[0]))
		if err != nil {
			return err
		}
	}

	return nil
}

// remove and keep order
func removeAlbumFrom(albums []Album, folderName string) ([]Album, *Album) {
	index := -1
	var removed *Album
	for i, a := range albums {
		if a.FolderName == folderName {
			index = i
			removedCopy := a
			removed = &removedCopy
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

func startsAscSort(albums []Album) func(i int, j int) bool {
	return func(i, j int) bool {
		return startsAscComparator(&albums[i], &albums[j]) > 0
	}
}

func albumDuration(album *Album) time.Duration {
	return album.End.Sub(album.Start)
}
