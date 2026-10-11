package catalog_test

import (
	"context"
	"slices"
	"time"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

var AlbumNotEmptyErr = errors.New("album cannot be deleted while it still contains medias")

type albumWithMedias struct {
	Album  catalog.Album
	Medias []*catalog.MediaMeta
}

// CatalogInMemory is the single in-memory Fake backing every test in pkg/catalog. It
// merges the albums and the medias into a single consistent store: a media only lives
// inside the album that holds it, so moving a media between albums is a single, atomic
// mutation and deleting an album that still contains medias is explicitly rejected
// (AlbumNotEmptyErr).
//
// It implements every port the production code consumes:
// catalog.TimelineRepository, catalog.MediaReadRepository, catalog.TransferMediasRepositoryPort,
// catalog.CountMediasBySelectorsPort, catalog.InsertMediasRepositoryPort, catalog.FindAlbumByOwnerPort,
// catalog.RepositoryAdapter.
type CatalogInMemory struct {
	albums map[catalog.AlbumId]*albumWithMedias
}

// AlbumSeed is the single seed structure consumed by NewCatalogInMemory: an album with
// (optionally) the medias it already contains. Setup is pre-execution: medias must be in
// their current album, never in the album they will end up in after the use case runs.
type AlbumSeed struct {
	Album  *catalog.Album
	Medias []*catalog.MediaMeta
}

func withAlbum(album *catalog.Album, medias ...*catalog.MediaMeta) AlbumSeed {
	return AlbumSeed{Album: album, Medias: medias}
}

// withMedias is a shortcut for albums whose Name / Start / End don't matter to the test;
// it fabricates a stub Album from the AlbumId so the only visible information is "there
// is an album with these medias in it".
func withMedias(albumId catalog.AlbumId, medias ...*catalog.MediaMeta) AlbumSeed {
	return AlbumSeed{
		Album:  &catalog.Album{AlbumId: albumId, Name: string(albumId.FolderName)},
		Medias: medias,
	}
}

func NewCatalogInMemory(seeds ...AlbumSeed) *CatalogInMemory {
	store := &CatalogInMemory{albums: make(map[catalog.AlbumId]*albumWithMedias)}
	for _, seed := range seeds {
		clone := *seed.Album
		store.albums[seed.Album.AlbumId] = &albumWithMedias{
			Album:  clone,
			Medias: append([]*catalog.MediaMeta(nil), seed.Medias...),
		}
	}
	return store
}

func (c *CatalogInMemory) AlbumIds() []catalog.AlbumId {
	ids := make([]catalog.AlbumId, 0, len(c.albums))
	for id := range c.albums {
		ids = append(ids, id)
	}
	return ids
}

func (c *CatalogInMemory) MediasByAlbum() map[catalog.AlbumId][]*catalog.MediaMeta {
	result := make(map[catalog.AlbumId][]*catalog.MediaMeta, len(c.albums))
	for id, entry := range c.albums {
		result[id] = append([]*catalog.MediaMeta(nil), entry.Medias...)
	}
	return result
}

func (c *CatalogInMemory) AlbumsByIds() map[catalog.AlbumId]catalog.Album {
	result := make(map[catalog.AlbumId]catalog.Album, len(c.albums))
	for id, entry := range c.albums {
		result[id] = entry.Album
	}
	return result
}

func (c *CatalogInMemory) FindAlbumsByOwner(_ context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
	var albums []*catalog.Album
	for _, entry := range c.albums {
		if entry.Album.Owner == owner {
			album := entry.Album
			albums = append(albums, &album)
		}
	}
	return albums, nil
}

func (c *CatalogInMemory) LoadTimeline(ctx context.Context, owner ownermodel.Owner) (*catalog.TimelineAggregate, error) {
	albums, err := c.FindAlbumsByOwner(ctx, owner)
	if err != nil {
		return nil, err
	}
	return catalog.NewTimelineAggregate(albums)
}

func (c *CatalogInMemory) FindAlbumByIds(_ context.Context, ids ...catalog.AlbumId) ([]*catalog.Album, error) {
	var albums []*catalog.Album
	for id, entry := range c.albums {
		if slices.Contains(ids, id) {
			album := entry.Album
			albums = append(albums, &album)
		}
	}
	return albums, nil
}

func (c *CatalogInMemory) CountMedia(_ context.Context, ids ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	result := make(map[catalog.AlbumId]int, len(ids))
	for _, id := range ids {
		if entry, ok := c.albums[id]; ok {
			result[id] = len(entry.Medias)
		}
	}
	return result, nil
}

func (c *CatalogInMemory) InsertAlbum(_ context.Context, album catalog.Album) error {
	if _, exists := c.albums[album.AlbumId]; exists {
		return errors.Errorf("album %s already exists", album.AlbumId)
	}
	saved := album
	c.albums[album.AlbumId] = &albumWithMedias{Album: saved}
	return nil
}

func (c *CatalogInMemory) DeleteAlbum(_ context.Context, albumId catalog.AlbumId) error {
	entry, ok := c.albums[albumId]
	if !ok {
		return catalog.AlbumNotFoundErr
	}
	if len(entry.Medias) > 0 {
		return errors.Wrapf(AlbumNotEmptyErr, "album %s still has %d medias", albumId, len(entry.Medias))
	}
	delete(c.albums, albumId)
	return nil
}

func (c *CatalogInMemory) UpdateAlbumName(_ context.Context, albumId catalog.AlbumId, newName string) (catalog.Album, catalog.Album, error) {
	entry, ok := c.albums[albumId]
	if !ok {
		return catalog.Album{}, catalog.Album{}, catalog.AlbumNotFoundErr
	}
	previous := entry.Album
	entry.Album.Name = newName
	return previous, entry.Album, nil
}

func (c *CatalogInMemory) AmendDates(_ context.Context, albumId catalog.AlbumId, start, end time.Time) error {
	entry, ok := c.albums[albumId]
	if !ok {
		return catalog.AlbumNotFoundErr
	}
	entry.Album.Start = start
	entry.Album.End = end
	return nil
}

func (c *CatalogInMemory) FindMedias(_ context.Context, request *catalog.FindMediaRequest) ([]*catalog.MediaMeta, error) {
	var medias []*catalog.MediaMeta
	for id, entry := range c.albums {
		if id.Owner != request.Owner {
			continue
		}
		if _, ok := request.AlbumFolderNames[id.FolderName]; !ok {
			continue
		}
		medias = append(medias, entry.Medias...)
	}
	return medias, nil
}

func (c *CatalogInMemory) FindMediaCurrentAlbum(_ context.Context, owner ownermodel.Owner, mediaId catalog.MediaId) (*catalog.AlbumId, error) {
	for id, entry := range c.albums {
		if id.Owner != owner {
			continue
		}
		for _, media := range entry.Medias {
			if media.Id == mediaId {
				albumId := id
				return &albumId, nil
			}
		}
	}
	return nil, catalog.MediaNotFoundError
}

func (c *CatalogInMemory) InsertMedias(_ context.Context, owner ownermodel.Owner, medias []catalog.CreateMediaRequest) error {
	for _, media := range medias {
		albumId := catalog.AlbumId{Owner: owner, FolderName: media.FolderName}
		entry, ok := c.albums[albumId]
		if !ok {
			return errors.Wrapf(catalog.AlbumNotFoundErr, "InsertMedias: album %s does not exist", albumId)
		}
		entry.Medias = append(entry.Medias, &catalog.MediaMeta{
			Id:        media.Id,
			Signature: media.Signature,
			Filename:  media.Filename,
			Type:      media.Type,
			Details:   media.Details,
		})
	}
	return nil
}

func (c *CatalogInMemory) TransferMediasFromRecords(_ context.Context, records catalog.MediaTransferRecords) (map[catalog.AlbumId][]catalog.MediaId, error) {
	transfers := make(map[catalog.AlbumId][]catalog.MediaId)
	for destination, selectors := range records {
		destEntry, destExists := c.albums[destination]
		for _, selector := range selectors {
			for _, source := range selector.FromAlbums {
				sourceEntry, ok := c.albums[source]
				if !ok {
					continue
				}
				remaining := sourceEntry.Medias[:0]
				for _, media := range sourceEntry.Medias {
					if selectorMatches(selector, media) {
						transfers[destination] = append(transfers[destination], media.Id)
						if destExists {
							destEntry.Medias = append(destEntry.Medias, media)
						}
					} else {
						remaining = append(remaining, media)
					}
				}
				sourceEntry.Medias = remaining
			}
		}
	}
	return transfers, nil
}

func (c *CatalogInMemory) CountMediasBySelectors(_ context.Context, owner ownermodel.Owner, selectors []catalog.MediaSelector) (int, error) {
	count := 0
	for _, selector := range selectors {
		for _, source := range selector.FromAlbums {
			if source.Owner != owner {
				continue
			}
			entry, ok := c.albums[source]
			if !ok {
				continue
			}
			for _, media := range entry.Medias {
				if selectorMatches(selector, media) {
					count++
				}
			}
		}
	}
	return count, nil
}

func selectorMatches(selector catalog.MediaSelector, media *catalog.MediaMeta) bool {
	date := media.Details.DateTime
	return !date.Before(selector.Start) && date.Before(selector.End)
}

// catalogTransferMediasFailing delegates every TransferMediasRepositoryPort call to the
// embedded CatalogInMemory except TransferMediasFromRecords, which returns Err. It lets
// use-case tests prove that an in-flight transfer failure leaves the catalog in a
// recoverable state: source medias are not moved, destination album is not touched, and
// downstream steps (delete, dates persistence, event-fire, ...) never execute.
type catalogTransferMediasFailing struct {
	*CatalogInMemory
	Err error
}

func (c *catalogTransferMediasFailing) TransferMediasFromRecords(_ context.Context, _ catalog.MediaTransferRecords) (map[catalog.AlbumId][]catalog.MediaId, error) {
	return nil, c.Err
}

// CoverRepositorySeed is used to pre-populate a CoverRepositoryInMemory.
type CoverRepositorySeed struct {
	AlbumId catalog.AlbumId
	Covers  []catalog.Cover
}

// CoverRepositoryInMemory implements catalog.CoverRepository backed by a map keyed by
// AlbumId. Cases can seed the fake with existing covers via the constructor.
type CoverRepositoryInMemory struct {
	Covers map[catalog.AlbumId][]catalog.Cover
}

func NewCoverRepositoryInMemory(seeds ...CoverRepositorySeed) *CoverRepositoryInMemory {
	store := &CoverRepositoryInMemory{
		Covers: make(map[catalog.AlbumId][]catalog.Cover),
	}
	for _, seed := range seeds {
		store.Covers[seed.AlbumId] = append([]catalog.Cover(nil), seed.Covers...)
	}
	return store
}

func (c *CoverRepositoryInMemory) FindCoversByAlbum(_ context.Context, albumId catalog.AlbumId) ([]catalog.Cover, error) {
	covers, ok := c.Covers[albumId]
	if !ok {
		return nil, nil
	}
	return append([]catalog.Cover(nil), covers...), nil
}

func (c *CoverRepositoryInMemory) FindCoversByAlbums(_ context.Context, albumIds ...catalog.AlbumId) (map[catalog.AlbumId][]catalog.Cover, error) {
	result := make(map[catalog.AlbumId][]catalog.Cover, len(albumIds))
	for _, albumId := range albumIds {
		if covers, ok := c.Covers[albumId]; ok {
			result[albumId] = append([]catalog.Cover(nil), covers...)
		}
	}
	return result, nil
}

func (c *CoverRepositoryInMemory) SaveCovers(_ context.Context, albumId catalog.AlbumId, covers []catalog.Cover) error {
	if len(covers) > catalog.MaxCoversPerAlbum {
		return catalog.TooManyCoversErr
	}
	if len(covers) == 0 {
		delete(c.Covers, albumId)
		return nil
	}
	c.Covers[albumId] = append([]catalog.Cover(nil), covers...)
	return nil
}

func (c *CoverRepositoryInMemory) MoveCovers(_ context.Context, from, to catalog.AlbumId) error {
	if from.IsEqual(to) {
		return nil
	}
	covers, ok := c.Covers[from]
	if !ok {
		return nil
	}
	c.Covers[to] = append([]catalog.Cover(nil), covers...)
	delete(c.Covers, from)
	return nil
}

// AlbumCreatedObserverInMemory implements catalog.AlbumCreatedObserver: it captures every
// AlbumCreated event notified to the observer.
type AlbumCreatedObserverInMemory struct {
	Events []catalog.AlbumCreated
}

func (c *AlbumCreatedObserverInMemory) OnAlbumCreated(_ context.Context, event catalog.AlbumCreated) error {
	c.Events = append(c.Events, event)
	return nil
}

// AlbumDeletedObserverInMemory implements catalog.AlbumDeletedObserver: it captures every
// AlbumDeleted event notified to the observer.
type AlbumDeletedObserverInMemory struct {
	Events []catalog.AlbumDeleted
}

func (d *AlbumDeletedObserverInMemory) OnAlbumDeleted(_ context.Context, event catalog.AlbumDeleted) error {
	d.Events = append(d.Events, event)
	return nil
}

// AlbumRenamedObserverInMemory implements catalog.AlbumRenamedObserver: it captures every
// AlbumRenamed event notified to the observer.
type AlbumRenamedObserverInMemory struct {
	Events []catalog.AlbumRenamed
}

func (r *AlbumRenamedObserverInMemory) OnAlbumRenamed(_ context.Context, event catalog.AlbumRenamed) error {
	r.Events = append(r.Events, event)
	return nil
}

// AlbumDatesAmendedObserverInMemory implements catalog.AlbumDatesAmendedObserver: it
// captures every AlbumDatesAmended event notified to the observer.
type AlbumDatesAmendedObserverInMemory struct {
	Events []catalog.AlbumDatesAmended
}

func (a *AlbumDatesAmendedObserverInMemory) OnAlbumDatesAmended(_ context.Context, event catalog.AlbumDatesAmended) error {
	a.Events = append(a.Events, event)
	return nil
}

// photoBuilder helps the tests describe a media by its date in a single line, with
// deterministic Id and Filename. Example: photoAt("photo5jan26", jan5_26) yields a media
// id "photo5jan26" and filename "photo5jan26.jpg" dated jan5_26.
func photoAt(id string, when time.Time) *catalog.MediaMeta {
	return &catalog.MediaMeta{
		Id:       catalog.MediaId(id),
		Filename: id + ".jpg",
		Type:     catalog.MediaTypeImage,
		Details:  catalog.MediaDetails{DateTime: when},
	}
}

func pickedCover(media *catalog.MediaMeta) catalog.Cover {
	return catalog.Cover{MediaId: media.Id, Filename: media.Filename, Origin: catalog.CoverOriginCherryPicked}
}

func randomCover(media *catalog.MediaMeta) catalog.Cover {
	return catalog.Cover{MediaId: media.Id, Filename: media.Filename, Origin: catalog.CoverOriginRandom}
}


