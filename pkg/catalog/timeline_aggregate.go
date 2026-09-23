package catalog

import (
	"fmt"
	"slices"
	"time"

	"github.com/pkg/errors"
)

type TimelineAggregate struct {
	timeline *Timeline
	albums   []*Album
}

// NewTimelineAggregate builds a TimelineAggregate eagerly: the timeline is computed
// upfront so any inconsistency in the persisted albums (like a duplicated AlbumId) is
// surfaced immediately. Every mutating method on the aggregate keeps this invariant by
// rebuilding the timeline from the updated in-memory album list.
func NewTimelineAggregate(albums []*Album) (*TimelineAggregate, error) {
	timeline, err := NewTimeline(albums)
	if err != nil {
		return nil, err
	}
	return &TimelineAggregate{
		albums:   albums,
		timeline: timeline,
	}, nil
}

// CreateNewAlbum validates the request, appends the new album to the aggregate, rebuilds the
// timeline, and returns both the created Album and the MediaTransferRecords for medias that
// overlapping albums should hand over to it. Returns AlbumFolderNameAlreadyTakenErr if the
// computed folder name collides with an existing album.
func (t *TimelineAggregate) CreateNewAlbum(request CreateAlbumRequest) (Album, MediaTransferRecords, error) {
	if err := request.IsValid(); err != nil {
		return Album{}, nil, errors.Wrapf(err, "CreateNewAlbum(%s) failed", request)
	}

	folderName := generateFolderName(request.Name, request.Start)
	if request.ForcedFolderName != "" && request.ForcedFolderName != "/" {
		folderName = NewFolderName(request.ForcedFolderName)
	}

	albumId := AlbumId{
		Owner:      request.Owner,
		FolderName: folderName,
	}

	nameIsAlreadyTaken := slices.ContainsFunc(t.albums, func(album *Album) bool {
		return album.AlbumId.IsEqual(albumId)
	})
	if nameIsAlreadyTaken {
		return Album{}, nil, errors.Wrapf(AlbumFolderNameAlreadyTakenErr, "%s album id already exists", albumId)
	}

	album := Album{
		AlbumId: albumId,
		Name:    request.Name,
		Start:   request.Start,
		End:     request.End,
	}
	t.albums = append(t.albums, &album)

	var err error
	t.timeline, err = NewTimeline(t.albums)
	if err != nil {
		return Album{}, nil, err
	}

	records := make(MediaTransferRecords)
	for _, seg := range t.timeline.FindForAlbum(album.AlbumId) {
		if len(seg.Albums) > 1 {
			selector := MediaSelector{
				FromAlbums: extractAlbumIds(seg.Albums[1:]),
				Start:      seg.Start,
				End:        seg.End,
			}
			if selectors, found := records[album.AlbumId]; found {
				records[album.AlbumId] = append(selectors, selector)
			} else {
				records[album.AlbumId] = []MediaSelector{selector}
			}
		}
	}

	if len(records) == 0 {
		return album, nil, nil
	}
	return album, records, nil
}

func extractAlbumIds(albums []Album) []AlbumId {
	if len(albums) == 0 {
		return nil
	}

	ids := make([]AlbumId, len(albums), len(albums))
	for i, album := range albums {
		ids[i] = album.AlbumId
	}

	return ids
}

type RenameAlbumRequest struct {
	CurrentId        AlbumId
	NewName          string
	RenameFolder     bool   // RenameFolder set to TRUE will create a new album with a FolderName generated from the NewName
	ForcedFolderName string // ForcedFolderName set to non-empty will create a new album with the requested FolderName (RenameFolder is ignored)
}

func (r RenameAlbumRequest) String() string {
	return fmt.Sprintf("%s -> %s", r.CurrentId.String(), r.NewName)
}

// AlbumNameUpdated carries the outcome of a folder-name-changing rename as computed by the
// TimelineAggregate: the existing album to remove, the new album to insert, and the media
// transfer that moves the existing album's medias into the new one.
type AlbumNameUpdated struct {
	ExistingAlbum Album
	RenamedAlbum  Album
	MediaTransfer MediaTransferRecords
}

// RenameAlbum removes the existing album from the aggregate and inserts a new one carrying
// the new folder name / display name. It returns everything the application service needs to
// persist the rename and move the medias:
//   - AlbumNotFoundErr if currentId is unknown to the aggregate.
//   - AlbumFolderNameAlreadyTakenErr if the target folder name is used by another album.
//
// The rename is expected to change the folder name; the caller is responsible for handling the
// name-only case (a simple UpdateAlbumName on the same row) before reaching this method.
func (t *TimelineAggregate) RenameAlbum(request RenameAlbumRequest) (*AlbumNameUpdated, error) {
	if request.NewName == "" {
		return nil, AlbumNameMandatoryErr
	}

	existing, err := t.findAlbumById(request.CurrentId)
	if err != nil {
		return nil, err
	}

	folderName := generateFolderName(request.NewName, existing.Start)
	if request.ForcedFolderName != "" && request.ForcedFolderName != "/" {
		folderName = NewFolderName(request.ForcedFolderName)
	}
	newAlbumId := AlbumId{Owner: request.CurrentId.Owner, FolderName: folderName}

	nameIsAlreadyTakenByAnother := slices.ContainsFunc(t.albums, func(album *Album) bool {
		return album.AlbumId.IsEqual(newAlbumId) && !album.AlbumId.IsEqual(request.CurrentId)
	})
	if nameIsAlreadyTakenByAnother {
		return nil, errors.Wrapf(AlbumFolderNameAlreadyTakenErr, "%s album id already exists", newAlbumId)
	}

	renamed := Album{
		AlbumId: newAlbumId,
		Name:    request.NewName,
		Start:   existing.Start,
		End:     existing.End,
	}

	for i, alb := range t.albums {
		if alb.AlbumId.IsEqual(existing.AlbumId) {
			t.albums[i] = &renamed
			break
		}
	}

	return &AlbumNameUpdated{
		ExistingAlbum: existing,
		RenamedAlbum:  renamed,
		MediaTransfer: MediaTransferRecords{
			newAlbumId: []MediaSelector{{
				FromAlbums: []AlbumId{existing.AlbumId},
				Start:      existing.Start,
				End:        existing.End,
			}},
		},
	}, nil
}

// FindAlbum returns a copy of the album carrying the given id, or AlbumNotFoundErr if the
// aggregate does not contain it. It is exposed so use cases that do not mutate the aggregate
// can still read from the same canonical source (e.g. an in-place rename that just needs the
// album's dates and current name to build an AlbumRenamed event).
func (t *TimelineAggregate) FindAlbum(albumId AlbumId) (Album, error) {
	return t.findAlbumById(albumId)
}

func (t *TimelineAggregate) findAlbumById(albumId AlbumId) (Album, error) {
	index := slices.IndexFunc(t.albums, func(album *Album) bool {
		return album.AlbumId.IsEqual(albumId)
	})
	if index == -1 {
		return Album{}, errors.Wrapf(AlbumNotFoundErr, "album %s not found", albumId)
	}
	existing := *t.albums[index]
	return existing, nil
}

func (t *TimelineAggregate) RemoveAlbum(deletedAlbumId AlbumId) (MediaTransferRecords, []MediaSelector, error) {
	records := make(MediaTransferRecords)
	orphaned := make([]MediaSelector, 0)
	var err error
	var deletedAlbum *Album

	t.albums, deletedAlbum = t.removeAlbumFrom(t.albums, deletedAlbumId.FolderName)
	if deletedAlbum == nil {
		return nil, nil, AlbumNotFoundErr
	}

	t.timeline, err = NewTimeline(t.albums)
	if err != nil {
		return nil, nil, err
	}

	segments := t.timeline.FindSegmentsBetween(deletedAlbum.Start, deletedAlbum.End)
	for _, seg := range segments {
		selector := MediaSelector{
			FromAlbums: []AlbumId{deletedAlbumId},
			Start:      seg.Start,
			End:        seg.End,
		}

		if len(seg.Albums) == 0 {
			orphaned = append(orphaned, selector)

		} else if priorityDescComparator(deletedAlbum, &seg.Albums[0]) > 0 {
			selectors, _ := records[seg.Albums[0].AlbumId]
			records[seg.Albums[0].AlbumId] = append(selectors, selector)
		}
	}

	return records, orphaned, nil
}

// removeAlbumFrom removes the album with the given folderName from the list of albums
func (t *TimelineAggregate) removeAlbumFrom(albums []*Album, folderName FolderName) ([]*Album, *Album) {
	for index, album := range albums {
		if album.FolderName == folderName {
			return append(albums[:index], albums[index+1:]...), album
		}
	}

	return albums, nil
}

// AlbumDatesUpdated carries the outcome of an album date amendment as computed by the
// TimelineAggregate: the DatesUpdate (with the previous start/end preserved so the caller
// can detect a no-op via DatesUpdate.DatesNotChanged()), the media transfer records needed
// to reallocate medias to their new album, and any selectors that would be left orphan.
type AlbumDatesUpdated struct {
	DatesUpdate   DatesUpdate
	MediaTransfer MediaTransferRecords
	Orphaned      []MediaSelector
}

// AmendDates validates the amendment and, if the dates actually change, mutates the
// aggregate to reflect the new dates and computes the media transfer records plus any
// orphaned selectors. When the dates are unchanged, the aggregate is left untouched and the
// returned AlbumDatesUpdated has DatesUpdate.DatesNotChanged() == true and empty transfer
// records: it is up to the caller (the use case) to short-circuit on that signal.
//
// Returns AlbumNotFoundErr if albumId is unknown to the aggregate.
func (t *TimelineAggregate) AmendDates(albumId AlbumId, start, end time.Time) (*AlbumDatesUpdated, error) {
	index := slices.IndexFunc(t.albums, func(album *Album) bool {
		return album.AlbumId.IsEqual(albumId)
	})
	if index == -1 {
		return nil, errors.Wrapf(AlbumNotFoundErr, "album %s not found", albumId)
	}

	previousAlbum := *t.albums[index]
	updatedAlbum := previousAlbum
	updatedAlbum.Start = start
	updatedAlbum.End = end

	datesUpdate := DatesUpdate{
		UpdatedAlbum:  updatedAlbum,
		PreviousStart: previousAlbum.Start,
		PreviousEnd:   previousAlbum.End,
	}

	if datesUpdate.DatesNotChanged() {
		return &AlbumDatesUpdated{DatesUpdate: datesUpdate}, nil
	}

	originalTimeline := t.timeline

	t.albums[index] = &updatedAlbum
	var err error
	t.timeline, err = NewTimeline(t.albums)
	if err != nil {
		return nil, err
	}

	rangeStart := minTime(updatedAlbum.Start, previousAlbum.Start)
	rangeEnd := maxTime(updatedAlbum.End, previousAlbum.End)

	cursor := struct {
		time             time.Time
		originalSegments []PrioritySegment
		amendedSegments  []PrioritySegment
		records          MediaTransferRecords
		orphaned         []MediaSelector
	}{
		time:             rangeStart,
		originalSegments: originalTimeline.FindSegmentsBetween(rangeStart, rangeEnd),
		amendedSegments:  t.timeline.FindSegmentsBetween(rangeStart, rangeEnd),
		records:          make(MediaTransferRecords),
	}

	for len(cursor.originalSegments) > 0 && len(cursor.amendedSegments) > 0 {
		nextTime := minTime(cursor.originalSegments[0].End, cursor.amendedSegments[0].End)
		wasLeading := t.isLeadByAlbum(albumId, cursor.originalSegments[0])
		takeTheLead := t.isLeadByAlbum(albumId, cursor.amendedSegments[0])

		if wasLeading && !takeTheLead {
			selector := MediaSelector{
				FromAlbums: []AlbumId{albumId},
				Start:      cursor.time,
				End:        nextTime,
			}

			if len(cursor.amendedSegments[0].Albums) == 0 {
				cursor.orphaned = append(cursor.orphaned, selector)

			} else {
				target := cursor.amendedSegments[0].Albums[0].AlbumId
				if selectors, found := cursor.records[target]; found {
					cursor.records[target] = append(selectors, selector)
				} else {
					cursor.records[target] = []MediaSelector{selector}
				}
			}
		} else if !wasLeading && takeTheLead && len(cursor.amendedSegments[0].Albums) > 1 {
			selector := MediaSelector{
				FromAlbums: extractAlbumIds(cursor.amendedSegments[0].Albums[1:]), // TODO this is required because there is no index only by dates (no album)
				Start:      cursor.time,
				End:        nextTime,
			}

			target := albumId
			if selectors, found := cursor.records[target]; found {
				cursor.records[target] = append(selectors, selector)
			} else {
				cursor.records[target] = []MediaSelector{selector}
			}
		}

		if cursor.originalSegments[0].End.Equal(cursor.amendedSegments[0].End) {
			cursor.time = cursor.originalSegments[0].End
			cursor.originalSegments = cursor.originalSegments[1:]
			cursor.amendedSegments = cursor.amendedSegments[1:]
		} else if cursor.originalSegments[0].End.Before(cursor.amendedSegments[0].End) {
			cursor.time = cursor.originalSegments[0].End
			cursor.originalSegments = cursor.originalSegments[1:]
		} else {
			cursor.time = cursor.amendedSegments[0].End
			cursor.amendedSegments = cursor.amendedSegments[1:]
		}
	}

	return &AlbumDatesUpdated{
		DatesUpdate:   datesUpdate,
		MediaTransfer: cursor.records,
		Orphaned:      cursor.orphaned,
	}, nil
}

func (t *TimelineAggregate) isLeadByAlbum(albumId AlbumId, seg PrioritySegment) bool {
	return len(seg.Albums) > 0 && seg.Albums[0].AlbumId.IsEqual(albumId)
}

func (t *TimelineAggregate) FindAt(date time.Time) (*Album, bool) {
	album, found := t.timeline.FindAt(date)
	return album, found
}
