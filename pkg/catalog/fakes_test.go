package catalog_test

import (
	"context"
	"slices"
	"time"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// AlbumRepositoryInMemory is a single in-memory backing store implementing every album-store
// port used across the pkg/catalog tests. Assertions are made through its read methods (or by
// inspecting the exported Albums map/CountByOwnerAndSelector state) rather than by recording
// method calls.
type AlbumRepositoryInMemory struct {
	Albums              map[catalog.AlbumId]*catalog.Album
	MediaCountsBySelect map[ownermodel.Owner]int // total medias by owner returned to CountMediasBySelectors regardless of selectors
}

func NewAlbumRepositoryInMemory(albums ...*catalog.Album) *AlbumRepositoryInMemory {
	r := &AlbumRepositoryInMemory{
		Albums:              make(map[catalog.AlbumId]*catalog.Album),
		MediaCountsBySelect: make(map[ownermodel.Owner]int),
	}
	for _, album := range albums {
		r.Albums[album.AlbumId] = album
	}
	return r
}

func (r *AlbumRepositoryInMemory) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
	var albums []*catalog.Album
	for _, album := range r.Albums {
		if album.Owner == owner {
			albums = append(albums, album)
		}
	}
	return albums, nil
}

func (r *AlbumRepositoryInMemory) FindAlbumByIds(ctx context.Context, ids ...catalog.AlbumId) ([]*catalog.Album, error) {
	var albums []*catalog.Album
	for _, album := range r.Albums {
		if slices.Contains(ids, album.AlbumId) {
			albums = append(albums, album)
		}
	}
	return albums, nil
}

func (r *AlbumRepositoryInMemory) FindAlbumById(ctx context.Context, id catalog.AlbumId) (*catalog.Album, error) {
	if album, ok := r.Albums[id]; ok {
		return album, nil
	}
	return nil, catalog.AlbumNotFoundErr
}

func (r *AlbumRepositoryInMemory) FindMedias(ctx context.Context, request *catalog.FindMediaRequest) ([]*catalog.MediaMeta, error) {
	panic("FindMedias not implemented in AlbumRepositoryInMemory")
}

func (r *AlbumRepositoryInMemory) FindMediaCurrentAlbum(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId) (*catalog.AlbumId, error) {
	panic("FindMediaCurrentAlbum not implemented in AlbumRepositoryInMemory")
}

func (r *AlbumRepositoryInMemory) CountMedia(ctx context.Context, ids ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	panic("CountMedia not implemented in AlbumRepositoryInMemory")
}

func (r *AlbumRepositoryInMemory) InsertAlbum(ctx context.Context, album catalog.Album) error {
	copyAlbum := album
	r.Albums[album.AlbumId] = &copyAlbum
	return nil
}

func (r *AlbumRepositoryInMemory) DeleteAlbum(ctx context.Context, albumId catalog.AlbumId) error {
	delete(r.Albums, albumId)
	return nil
}

func (r *AlbumRepositoryInMemory) UpdateAlbumName(ctx context.Context, albumId catalog.AlbumId, newName string) error {
	album, ok := r.Albums[albumId]
	if !ok {
		return catalog.AlbumNotFoundErr
	}
	album.Name = newName
	return nil
}

func (r *AlbumRepositoryInMemory) AmendDates(ctx context.Context, albumId catalog.AlbumId, start, end time.Time) error {
	album, ok := r.Albums[albumId]
	if !ok {
		return catalog.AlbumNotFoundErr
	}
	album.Start = start
	album.End = end
	return nil
}

func (r *AlbumRepositoryInMemory) CountMediasBySelectors(ctx context.Context, owner ownermodel.Owner, selectors []catalog.MediaSelector) (int, error) {
	return r.MediaCountsBySelect[owner], nil
}

// MediaTransferInMemory records the media transfer requests made and returns a canned
// TransferredMedias value so downstream observers can be exercised. It implements both
// catalog.MediaTransfer and catalog.TransferMediasRepositoryPort.
type MediaTransferInMemory struct {
	TransferRecords   []catalog.MediaTransferRecords
	TransferredMedias catalog.TransferredMedias
}

func (m *MediaTransferInMemory) Transfer(ctx context.Context, records catalog.MediaTransferRecords) error {
	m.TransferRecords = append(m.TransferRecords, records)
	return nil
}

func (m *MediaTransferInMemory) TransferMediasFromRecords(ctx context.Context, records catalog.MediaTransferRecords) (catalog.TransferredMedias, error) {
	m.TransferRecords = append(m.TransferRecords, records)
	return m.TransferredMedias, nil
}

// TimelineMutationObserverInMemory records every notification sent to it.
type TimelineMutationObserverInMemory struct {
	Notifications []catalog.TransferredMedias
}

func (t *TimelineMutationObserverInMemory) OnTransferredMedias(ctx context.Context, transfers catalog.TransferredMedias) error {
	t.Notifications = append(t.Notifications, transfers)
	return nil
}

// CreateAlbumObserverInMemory records every album passed to ObserveCreateAlbum.
type CreateAlbumObserverInMemory struct {
	CreatedAlbums []catalog.Album
}

func (c *CreateAlbumObserverInMemory) ObserveCreateAlbum(ctx context.Context, createdAlbum catalog.Album) error {
	c.CreatedAlbums = append(c.CreatedAlbums, createdAlbum)
	return nil
}

// DeleteAlbumObserverInMemory records every OnDeleteAlbum invocation.
type DeleteAlbumObserverInMemory struct {
	Deleted   []catalog.AlbumId
	Transfers []catalog.MediaTransferRecords
}

func (d *DeleteAlbumObserverInMemory) OnDeleteAlbum(ctx context.Context, deletedAlbum catalog.AlbumId, transfers catalog.MediaTransferRecords) error {
	d.Deleted = append(d.Deleted, deletedAlbum)
	d.Transfers = append(d.Transfers, transfers)
	return nil
}

// RenameAlbumObserverInMemory records every OnRenameAlbum invocation.
type RenameAlbumObserverInMemory struct {
	RenamedFrom      []catalog.AlbumId
	CreationRequests []catalog.CreateAlbumRequest
}

func (r *RenameAlbumObserverInMemory) OnRenameAlbum(ctx context.Context, current catalog.AlbumId, creationRequest catalog.CreateAlbumRequest) error {
	r.RenamedFrom = append(r.RenamedFrom, current)
	r.CreationRequests = append(r.CreationRequests, creationRequest)
	return nil
}

// AlbumDatesAmendedObserverInMemory records every dates amended notification, both with and
// without a timeline.
type AlbumDatesAmendedObserverInMemory struct {
	DateAmendedAlbums []catalog.DatesUpdate
}

func (a *AlbumDatesAmendedObserverInMemory) OnAlbumDatesAmendedWithTimeline(ctx context.Context, timeline *catalog.TimelineAggregate, amendedAlbum catalog.DatesUpdate) error {
	a.DateAmendedAlbums = append(a.DateAmendedAlbums, amendedAlbum)
	return nil
}

func (a *AlbumDatesAmendedObserverInMemory) OnAlbumDatesAmended(ctx context.Context, amendedAlbum catalog.DatesUpdate) error {
	a.DateAmendedAlbums = append(a.DateAmendedAlbums, amendedAlbum)
	return nil
}
