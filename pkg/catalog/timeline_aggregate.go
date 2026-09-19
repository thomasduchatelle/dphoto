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

// NewLazyTimelineAggregate creates a new TimelineAggregate without timeline pre-computation. The timeline will be computed at the first AddNew call.
func NewLazyTimelineAggregate(albums []*Album) *TimelineAggregate {
	return &TimelineAggregate{
		albums: albums,
	}
}

func NewInitialisedTimelineAggregate(albums []*Album) (*TimelineAggregate, error) {
	timeline, err := NewTimeline(albums)
	return &TimelineAggregate{
		albums:   albums,
		timeline: timeline,
	}, err
}

func (t *TimelineAggregate) CreateNewAlbum(request CreateAlbumRequest) (Album, error) {
	if err := request.IsValid(); err != nil {
		return Album{}, errors.Wrapf(err, "CreateNewAlbum(%s) failed", request)
	}

	album, err := t.convert(request)
	if err != nil {
		return Album{}, err
	}

	return *album, nil
}

func (t *TimelineAggregate) convert(request CreateAlbumRequest) (*Album, error) {
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
		return nil, errors.Wrapf(AlbumFolderNameAlreadyTakenErr, "%s album id already exists", albumId)
	}

	return &Album{
		AlbumId: albumId,
		Name:    request.Name,
		Start:   request.Start,
		End:     request.End,
	}, nil
}

func (t *TimelineAggregate) AddNew(addedAlbum Album) (MediaTransferRecords, error) {
	t.albums = append(t.albums, &addedAlbum)

	var err error
	t.timeline, err = NewTimeline(t.albums)
	if err != nil {
		return nil, err
	}

	records := make(MediaTransferRecords)
	for _, seg := range t.timeline.FindForAlbum(addedAlbum.AlbumId) {
		if len(seg.Albums) > 1 {
			selector := MediaSelector{
				FromAlbums: extractAlbumIds(seg.Albums[1:]),
				Start:      seg.Start,
				End:        seg.End,
			}
			if selectors, found := records[addedAlbum.AlbumId]; found {
				records[addedAlbum.AlbumId] = append(selectors, selector)
			} else {
				records[addedAlbum.AlbumId] = []MediaSelector{selector}
			}
		}
	}

	if len(records) == 0 {
		return nil, nil
	}

	return records, nil
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

// AlbumRenamed carries the outcome of a folder-name-changing rename: the existing album to
// remove, the new album to insert, and the media transfer that moves the existing album's
// medias into the new one.
type AlbumRenamed struct {
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
func (t *TimelineAggregate) RenameAlbum(request RenameAlbumRequest) (*AlbumRenamed, error) {
	if t.timeline == nil {
		return nil, errors.Errorf("TimelineAggregate.RenameAlbum must be called from NewInitialisedTimelineAggregate (t.timeline is nil)")
	}
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

	return &AlbumRenamed{
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

func (t *TimelineAggregate) ValidateAmendDates(albumId AlbumId, start, end time.Time) (*DatesUpdate, error) {
	index := slices.IndexFunc(t.albums, func(album *Album) bool {
		return album.AlbumId.IsEqual(albumId)
	})

	if index == -1 {
		return nil, errors.Wrapf(AlbumNotFoundErr, "album %s not found", albumId)
	}

	amended := DatesUpdate{
		UpdatedAlbum:  *t.albums[index],
		PreviousStart: t.albums[index].Start,
		PreviousEnd:   t.albums[index].End,
	}
	amended.UpdatedAlbum.Start = start
	amended.UpdatedAlbum.End = end

	return &amended, nil
}

func (t *TimelineAggregate) AmendDates(amendedAlbum DatesUpdate) (MediaTransferRecords, []MediaSelector, error) {
	index := slices.IndexFunc(t.albums, func(album *Album) bool {
		return album.AlbumId.IsEqual(amendedAlbum.UpdatedAlbum.AlbumId)
	})
	if index == -1 {
		return nil, nil, errors.Wrapf(AlbumNotFoundErr, "album %s not found", amendedAlbum.UpdatedAlbum.AlbumId)
	}

	var err error

	previousAlbum := *t.albums[index]
	if !previousAlbum.Start.Equal(amendedAlbum.PreviousStart) || !previousAlbum.End.Equal(amendedAlbum.PreviousEnd) {
		// restore TimelineAggregate as it was before the amendment
		previousAlbum.Start = amendedAlbum.PreviousStart
		previousAlbum.End = amendedAlbum.PreviousEnd
		t.albums[index] = &previousAlbum

		t.timeline, err = NewTimeline(t.albums)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to create timeline during AmendDates(%s)", amendedAlbum)
		}
	}

	originalTimeline, err := t.getTimeline()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create timeline during AmendDates(%s)", amendedAlbum)
	}

	t.albums[index] = &amendedAlbum.UpdatedAlbum
	t.timeline, err = NewTimeline(t.albums)
	if err != nil {
		return nil, nil, err
	}

	start := minTime(amendedAlbum.UpdatedAlbum.Start, previousAlbum.Start)
	end := maxTime(amendedAlbum.UpdatedAlbum.End, previousAlbum.End)

	cursor := struct {
		time             time.Time
		originalSegments []PrioritySegment
		amendedSegments  []PrioritySegment
		records          MediaTransferRecords
		orphaned         []MediaSelector
	}{
		time:             start,
		originalSegments: originalTimeline.FindSegmentsBetween(start, end),
		amendedSegments:  t.timeline.FindSegmentsBetween(start, end),
		records:          make(MediaTransferRecords),
	}

	for len(cursor.originalSegments) > 0 && len(cursor.amendedSegments) > 0 {
		nextTime := minTime(cursor.originalSegments[0].End, cursor.amendedSegments[0].End)
		wasLeading := t.isLeadByAlbum(amendedAlbum.UpdatedAlbum.AlbumId, cursor.originalSegments[0])
		takeTheLead := t.isLeadByAlbum(amendedAlbum.UpdatedAlbum.AlbumId, cursor.amendedSegments[0])

		if wasLeading && !takeTheLead {
			selector := MediaSelector{
				FromAlbums: []AlbumId{amendedAlbum.UpdatedAlbum.AlbumId},
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

			target := amendedAlbum.UpdatedAlbum.AlbumId
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

	return cursor.records, cursor.orphaned, nil
}

func (t *TimelineAggregate) isLeadByAlbum(albumId AlbumId, seg PrioritySegment) bool {
	return len(seg.Albums) > 0 && seg.Albums[0].AlbumId.IsEqual(albumId)
}

func (t *TimelineAggregate) FindAt(date time.Time) (*Album, bool, error) {
	_, err := t.getTimeline()
	if err != nil {
		return nil, false, errors.Wrapf(err, "failed to create timeline during FindAt(%s)", date)
	}

	albumId, found := t.timeline.FindAt(date)
	return albumId, found, nil
}

// getTimeline returns the timeline, creating it if it doesn't exist yet (lazy initialisation)
func (t *TimelineAggregate) getTimeline() (*Timeline, error) {
	var err error
	if t.timeline == nil {
		t.timeline, err = NewTimeline(t.albums)
	}

	return t.timeline, err
}
