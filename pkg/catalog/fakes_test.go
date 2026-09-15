package catalog_test

import (
	"context"
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
// records passed to TransferMediasFromRecords and returns the pre-set TransferredMedias.
type TransferMediasInMemory struct {
	Records           []catalog.MediaTransferRecords
	TransferredMedias catalog.TransferredMedias
}

func NewTransferMediasInMemory() *TransferMediasInMemory {
	return &TransferMediasInMemory{
		TransferredMedias: catalog.NewTransferredMedias(),
	}
}

func (t *TransferMediasInMemory) TransferMediasFromRecords(_ context.Context, records catalog.MediaTransferRecords) (catalog.TransferredMedias, error) {
	t.Records = append(t.Records, records)
	return t.TransferredMedias, nil
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

// ReplaceAlbumCall represents one call to OnReplaceAlbum.
type ReplaceAlbumCall struct {
	Current         catalog.AlbumId
	CreationRequest catalog.CreateAlbumRequest
}

// ReplaceAlbumObserverInMemory implements catalog.ReplaceAlbumObserver: it captures every
// folder-changing rename notified to the observer.
type ReplaceAlbumObserverInMemory struct {
	Replaced []ReplaceAlbumCall
}

func (r *ReplaceAlbumObserverInMemory) OnReplaceAlbum(_ context.Context, current catalog.AlbumId, creationRequest catalog.CreateAlbumRequest) error {
	r.Replaced = append(r.Replaced, ReplaceAlbumCall{Current: current, CreationRequest: creationRequest})
	return nil
}

// AlbumRenamedCall represents one call to OnAlbumRenamed.
type AlbumRenamedCall struct {
	AlbumId catalog.AlbumId
	NewName string
}

// RenameAlbumObserverInMemory implements catalog.RenameAlbumObserver: it captures every
// in-place rename notified to the observer.
type RenameAlbumObserverInMemory struct {
	Renamed []AlbumRenamedCall
}

func (r *RenameAlbumObserverInMemory) OnAlbumRenamed(_ context.Context, albumId catalog.AlbumId, newName string) error {
	r.Renamed = append(r.Renamed, AlbumRenamedCall{AlbumId: albumId, NewName: newName})
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
