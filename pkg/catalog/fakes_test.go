package catalog_test

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// AlbumRepositoryInMemory is a single in-memory Fake that implements every album-store port
// used across the pkg/catalog tests. It is backed by a single map of albums keyed by AlbumId.
//
// When a test needs a non-zero result from CountMediasBySelectors, set MediasBySelector to a
// function returning the count for a given (owner, selector). By default, every selector reports
// zero medias.
type AlbumRepositoryInMemory struct {
	Albums           map[catalog.AlbumId]*catalog.Album
	MediasBySelector func(owner ownermodel.Owner, selector catalog.MediaSelector) int
}

func NewAlbumRepositoryInMemory(albums ...*catalog.Album) *AlbumRepositoryInMemory {
	store := &AlbumRepositoryInMemory{
		Albums: make(map[catalog.AlbumId]*catalog.Album),
	}
	for _, album := range albums {
		store.Albums[album.AlbumId] = album
	}
	return store
}

func (r *AlbumRepositoryInMemory) FindAlbumsByOwner(_ context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
	var albums []*catalog.Album
	for _, album := range r.Albums {
		if album.Owner == owner {
			albums = append(albums, album)
		}
	}
	return albums, nil
}

func (r *AlbumRepositoryInMemory) FindAlbumByIds(_ context.Context, ids ...catalog.AlbumId) ([]*catalog.Album, error) {
	var albums []*catalog.Album
	for _, album := range r.Albums {
		if slices.Contains(ids, album.AlbumId) {
			albums = append(albums, album)
		}
	}
	return albums, nil
}

func (r *AlbumRepositoryInMemory) FindAlbumById(_ context.Context, id catalog.AlbumId) (*catalog.Album, error) {
	if album, ok := r.Albums[id]; ok {
		return album, nil
	}
	return nil, catalog.AlbumNotFoundErr
}

func (r *AlbumRepositoryInMemory) FindMedias(_ context.Context, _ *catalog.FindMediaRequest) ([]*catalog.MediaMeta, error) {
	panic("FindMedias not implemented on AlbumRepositoryInMemory")
}

func (r *AlbumRepositoryInMemory) FindMediaCurrentAlbum(_ context.Context, _ ownermodel.Owner, _ catalog.MediaId) (*catalog.AlbumId, error) {
	panic("FindMediaCurrentAlbum not implemented on AlbumRepositoryInMemory")
}

func (r *AlbumRepositoryInMemory) CountMedia(_ context.Context, _ ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	panic("CountMedia not implemented on AlbumRepositoryInMemory")
}

func (r *AlbumRepositoryInMemory) InsertAlbum(_ context.Context, album catalog.Album) error {
	saved := album
	r.Albums[album.AlbumId] = &saved
	return nil
}

func (r *AlbumRepositoryInMemory) DeleteAlbum(_ context.Context, albumId catalog.AlbumId) error {
	delete(r.Albums, albumId)
	return nil
}

func (r *AlbumRepositoryInMemory) UpdateAlbumName(_ context.Context, albumId catalog.AlbumId, newName string) error {
	album, ok := r.Albums[albumId]
	if !ok {
		return catalog.AlbumNotFoundErr
	}
	album.Name = newName
	return nil
}

func (r *AlbumRepositoryInMemory) AmendDates(_ context.Context, albumId catalog.AlbumId, start, end time.Time) error {
	album, ok := r.Albums[albumId]
	if !ok {
		return catalog.AlbumNotFoundErr
	}
	album.Start = start
	album.End = end
	return nil
}

func (r *AlbumRepositoryInMemory) CountMediasBySelectors(_ context.Context, owner ownermodel.Owner, selectors []catalog.MediaSelector) (int, error) {
	if r.MediasBySelector == nil {
		return 0, nil
	}
	count := 0
	for _, selector := range selectors {
		count += r.MediasBySelector(owner, selector)
	}
	return count, nil
}

// TransferMediasInMemory implements catalog.TransferMediasRepositoryPort: it captures every
// records passed to TransferMediasFromRecords and returns the pre-set Transferred map.
type TransferMediasInMemory struct {
	Records     []catalog.MediaTransferRecords
	Transferred map[catalog.AlbumId][]catalog.MediaId
}

func NewTransferMediasInMemory() *TransferMediasInMemory {
	return &TransferMediasInMemory{
		Transferred: make(map[catalog.AlbumId][]catalog.MediaId),
	}
}

func (t *TransferMediasInMemory) TransferMediasFromRecords(_ context.Context, records catalog.MediaTransferRecords) (map[catalog.AlbumId][]catalog.MediaId, error) {
	t.Records = append(t.Records, records)
	return t.Transferred, nil
}

// MediaTransferInMemory implements catalog.MediaTransfer: it captures every records passed to
// Transfer.
type MediaTransferInMemory struct {
	Records []catalog.MediaTransferRecords
}

func (m *MediaTransferInMemory) Transfer(_ context.Context, records catalog.MediaTransferRecords) error {
	m.Records = append(m.Records, records)
	return nil
}

// TimelineMutationObserverInMemory implements catalog.TimelineMutationObserver: it captures
// every transfer notified to the observer.
type TimelineMutationObserverInMemory struct {
	Notifications []catalog.TransferredMedias
}

func (o *TimelineMutationObserverInMemory) OnTransferredMedias(_ context.Context, transfers catalog.TransferredMedias) error {
	o.Notifications = append(o.Notifications, transfers)
	return nil
}

// CreateAlbumObserverInMemory implements catalog.CreateAlbumObserver: it captures every album
// created through the observer.
type CreateAlbumObserverInMemory struct {
	CreatedAlbums []catalog.Album
}

func (c *CreateAlbumObserverInMemory) ObserveCreateAlbum(_ context.Context, createdAlbum catalog.Album) error {
	c.CreatedAlbums = append(c.CreatedAlbums, createdAlbum)
	return nil
}

// DeleteAlbumObserverInMemory implements catalog.DeleteAlbumObserver: it captures every album
// deleted through the observer.
type DeleteAlbumObserverInMemory struct {
	Deleted   []catalog.AlbumId
	Transfers []catalog.MediaTransferRecords
}

func (d *DeleteAlbumObserverInMemory) OnDeleteAlbum(_ context.Context, deletedAlbum catalog.AlbumId, transfers catalog.MediaTransferRecords) error {
	d.Deleted = append(d.Deleted, deletedAlbum)
	d.Transfers = append(d.Transfers, transfers)
	return nil
}

// RenameAlbumCall represents one call to OnRenameAlbum.
type RenameAlbumCall struct {
	Current         catalog.AlbumId
	CreationRequest catalog.CreateAlbumRequest
}

// RenameAlbumObserverInMemory implements catalog.RenameAlbumObserver: it captures every rename
// notified to the observer.
type RenameAlbumObserverInMemory struct {
	Renamed []RenameAlbumCall
}

func (r *RenameAlbumObserverInMemory) OnRenameAlbum(_ context.Context, current catalog.AlbumId, creationRequest catalog.CreateAlbumRequest) error {
	r.Renamed = append(r.Renamed, RenameAlbumCall{Current: current, CreationRequest: creationRequest})
	return nil
}

// AlbumDatesAmendedObserverInMemory implements both catalog.AlbumDatesAmendedObserverWithTimeline
// and catalog.AlbumDatesAmendedObserver: it captures every amended-dates event.
type AlbumDatesAmendedObserverInMemory struct {
	DateAmendedAlbums []catalog.DatesUpdate
}

func (a *AlbumDatesAmendedObserverInMemory) OnAlbumDatesAmendedWithTimeline(_ context.Context, _ *catalog.TimelineAggregate, amendedAlbum catalog.DatesUpdate) error {
	a.DateAmendedAlbums = append(a.DateAmendedAlbums, amendedAlbum)
	return nil
}

func (a *AlbumDatesAmendedObserverInMemory) OnAlbumDatesAmended(_ context.Context, amendedAlbum catalog.DatesUpdate) error {
	a.DateAmendedAlbums = append(a.DateAmendedAlbums, amendedAlbum)
	return nil
}

// TransferMediasServiceFake implements catalog.TransferMediasService: for each destination
// album in the records, it fabricates one MediaId per (source album, day in the selector
// range) so tests can rely on a deterministic, non-empty result without wiring a real
// repository. It also captures every Records passed to TransferMedias.
//
// The generated MediaId format is:
//
//	fake-media-<owner>-<sourceFolder>-<yyyy-mm-dd>
//
// where <sourceFolder> is the source album's FolderName with leading '/' removed. FromAlbums
// on the returned TransferredMedias is derived from the input records (origins that are not
// themselves destinations), matching the production behaviour.
type TransferMediasServiceFake struct {
	Records []catalog.MediaTransferRecords
}

func (f *TransferMediasServiceFake) TransferMedias(_ context.Context, records catalog.MediaTransferRecords) (catalog.TransferredMedias, error) {
	f.Records = append(f.Records, records)

	transfers := make(map[catalog.AlbumId][]catalog.MediaId)
	for destination, selectors := range records {
		var ids []catalog.MediaId
		for _, selector := range selectors {
			for _, source := range selector.FromAlbums {
				for day := truncateToDay(selector.Start); day.Before(selector.End); day = day.AddDate(0, 0, 1) {
					ids = append(ids, fakeMediaId(source, day))
				}
			}
		}
		if len(ids) > 0 {
			transfers[destination] = ids
		}
	}

	result := catalog.TransferredMedias{Transfers: transfers}
	if result.IsEmpty() {
		return result, nil
	}

	for _, selectors := range records {
		for _, selector := range selectors {
			for _, source := range selector.FromAlbums {
				if _, isDestination := result.Transfers[source]; isDestination {
					continue
				}
				if slices.Contains(result.FromAlbums, source) {
					continue
				}
				result.FromAlbums = append(result.FromAlbums, source)
			}
		}
	}

	return result, nil
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func fakeMediaId(source catalog.AlbumId, day time.Time) catalog.MediaId {
	folder := string(source.FolderName)
	if len(folder) > 0 && folder[0] == '/' {
		folder = folder[1:]
	}
	return catalog.MediaId(fmt.Sprintf("fake-media-%s-%s-%s", source.Owner, folder, day.Format("2006-01-02")))
}
