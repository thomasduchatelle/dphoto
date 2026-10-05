package catalog_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func TestDeleteAlbum_DeleteAlbum(t *testing.T) {
	const owner = "ironman"
	mar24 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	may24 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	jan24 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	nov24 := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
	dec24 := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)

	toDeleteAlbumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers-1")}
	toDeleteAlbum := &catalog.Album{AlbumId: toDeleteAlbumId, Name: "Avenger 1", Start: mar24, End: may24}
	allYearAlbumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/lifetime")}
	allYearAlbum := &catalog.Album{AlbumId: allYearAlbumId, Name: "lifetime", Start: jan24, End: jan25}
	isolatedAlbumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/isolated")}
	isolatedAlbum := &catalog.Album{AlbumId: isolatedAlbumId, Name: "Isolated", Start: nov24, End: dec24}

	testError := errors.Errorf("TEST error throwing")

	mediaAt := func(id string, albumId catalog.AlbumId, when time.Time) *catalog.MediaMeta {
		return &catalog.MediaMeta{
			Id:       catalog.MediaId(id),
			Filename: id + ".jpg",
			Type:     catalog.MediaTypeImage,
			Details:  catalog.MediaDetails{DateTime: when},
		}
	}

	movingMedia := mediaAt("movie-1", toDeleteAlbumId, mar24.AddDate(0, 0, 10))
	otherMovingMedia := mediaAt("movie-2", toDeleteAlbumId, mar24.AddDate(0, 0, 20))

	t.Run("it should transfer medias to the surrounding album and inherit the deleted album's CHERRY_PICKED cover", func(t *testing.T) {
		albumRepository := NewAlbumRepositoryInMemory(allYearAlbum, toDeleteAlbum)
		mediaStore := &MediaReadRepositoryInMemory{
			Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
				toDeleteAlbumId: {movingMedia, otherMovingMedia},
			},
		}
		coverRepository := NewCoverRepositoryInMemory(coversFor(toDeleteAlbumId,
			catalog.Cover{MediaId: movingMedia.Id, Filename: movingMedia.Filename, Origin: catalog.CoverOriginCherryPicked},
			catalog.Cover{MediaId: otherMovingMedia.Id, Filename: otherMovingMedia.Filename, Origin: catalog.CoverOriginRandom},
		))
		observer := &AlbumDeletedObserverInMemory{}

		deleteAlbum := catalog.NewDeleteAlbum(
			albumRepository,
			mediaStore,
			&catalog.TransferMediasFromRepository{TransferMediasRepository: mediaStore},
			&catalog.CoverService{
				CoverRepository:     coverRepository,
				MediaReadRepository: mediaStore,
				Randomiser:          deterministicRandomiser,
			},
			observer,
		)

		err := deleteAlbum.DeleteAlbum(context.Background(), toDeleteAlbumId)
		if !assert.NoError(t, err) {
			return
		}

		assert.NotContains(t, albumRepository.Albums, toDeleteAlbumId, "deleted album is gone from the repository")
		assert.ElementsMatch(t, []*catalog.MediaMeta{movingMedia, otherMovingMedia}, mediaStore.Medias[allYearAlbumId], "medias have been transferred to the surrounding album")
		assert.Empty(t, mediaStore.Medias[toDeleteAlbumId], "no medias left on the deleted album")

		assert.NotContains(t, coverRepository.Covers, toDeleteAlbumId, "deleted album's canonical cover record is removed")
		assert.Equal(t, []catalog.Cover{
			{MediaId: movingMedia.Id, Filename: movingMedia.Filename, Origin: catalog.CoverOriginCherryPicked},
			{MediaId: otherMovingMedia.Id, Filename: otherMovingMedia.Filename, Origin: catalog.CoverOriginRandom},
		}, coverRepository.Covers[allYearAlbumId], "destination inherits CHERRY_PICKED cover and fills remaining slots")

		if assert.Len(t, observer.Events, 1) {
			event := observer.Events[0]
			assert.Equal(t, toDeleteAlbumId, event.DeletedAlbumId)
			assert.ElementsMatch(t, []catalog.MediaId{movingMedia.Id, otherMovingMedia.Id}, event.TransferredMedias.Transfers[allYearAlbumId])
			assert.Equal(t, []catalog.AlbumId{toDeleteAlbumId}, event.TransferredMedias.FromAlbums)
			assert.Contains(t, event.Covers, toDeleteAlbumId, "deleted album is in the Covers map")
			assert.Nil(t, event.Covers[toDeleteAlbumId], "deleted album covers are nil (sentinel for 'delete')")
			assert.NotNil(t, event.Covers[allYearAlbumId], "destination has new covers")
		}
	})

	t.Run("it should delete the canonical cover record and fire an empty-transfer event when no surrounding album covers it", func(t *testing.T) {
		albumRepository := NewAlbumRepositoryInMemory(isolatedAlbum, toDeleteAlbum)
		mediaStore := &MediaReadRepositoryInMemory{Medias: map[catalog.AlbumId][]*catalog.MediaMeta{}}
		coverRepository := NewCoverRepositoryInMemory(coversFor(toDeleteAlbumId,
			catalog.Cover{MediaId: "ghost", Filename: "ghost.jpg", Origin: catalog.CoverOriginCherryPicked},
		))
		observer := &AlbumDeletedObserverInMemory{}

		deleteAlbum := catalog.NewDeleteAlbum(
			albumRepository,
			mediaStore,
			&catalog.TransferMediasFromRepository{TransferMediasRepository: mediaStore},
			&catalog.CoverService{
				CoverRepository:     coverRepository,
				MediaReadRepository: mediaStore,
				Randomiser:          deterministicRandomiser,
			},
			observer,
		)

		err := deleteAlbum.DeleteAlbum(context.Background(), toDeleteAlbumId)
		if !assert.NoError(t, err) {
			return
		}

		assert.NotContains(t, albumRepository.Albums, toDeleteAlbumId, "deleted album is gone")
		assert.Contains(t, albumRepository.Albums, isolatedAlbumId, "isolated album is left untouched")
		assert.NotContains(t, coverRepository.Covers, toDeleteAlbumId, "deleted album's canonical cover record is removed")

		if assert.Len(t, observer.Events, 1) {
			event := observer.Events[0]
			assert.Equal(t, toDeleteAlbumId, event.DeletedAlbumId)
			assert.True(t, event.TransferredMedias.IsEmpty(), "no media was transferred")
			assert.Equal(t, map[catalog.AlbumId][]catalog.Cover{toDeleteAlbumId: nil}, event.Covers, "deleted album is the only entry in Covers, mapped to nil")
		}
	})

	t.Run("it should return OrphanedMediasErr without touching any album when the deletion would orphan medias", func(t *testing.T) {
		q1AlbumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/q1")}
		q1Album := &catalog.Album{AlbumId: q1AlbumId, Name: "q1", Start: jan24, End: apr24}
		albumRepository := NewAlbumRepositoryInMemory(q1Album, toDeleteAlbum)
		mediaStore := &MediaReadRepositoryInMemory{
			Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
				toDeleteAlbumId: {mediaAt("orphan-1", toDeleteAlbumId, apr24.AddDate(0, 0, 10))},
			},
		}
		coverRepository := NewCoverRepositoryInMemory()
		observer := &AlbumDeletedObserverInMemory{}

		deleteAlbum := catalog.NewDeleteAlbum(
			albumRepository,
			mediaStore,
			&catalog.TransferMediasFromRepository{TransferMediasRepository: mediaStore},
			&catalog.CoverService{
				CoverRepository:     coverRepository,
				MediaReadRepository: mediaStore,
				Randomiser:          deterministicRandomiser,
			},
			observer,
		)

		err := deleteAlbum.DeleteAlbum(context.Background(), toDeleteAlbumId)
		assert.ErrorIs(t, err, catalog.OrphanedMediasErr)
		assert.Contains(t, albumRepository.Albums, toDeleteAlbumId, "no album was removed")
		assert.NotEmpty(t, mediaStore.Medias[toDeleteAlbumId], "no media was moved")
		assert.Empty(t, observer.Events, "no event was fired")
	})

	t.Run("it should return AlbumNotFoundErr when the album does not exist", func(t *testing.T) {
		albumRepository := NewAlbumRepositoryInMemory(allYearAlbum)
		mediaStore := &MediaReadRepositoryInMemory{}
		coverRepository := NewCoverRepositoryInMemory()
		observer := &AlbumDeletedObserverInMemory{}

		deleteAlbum := catalog.NewDeleteAlbum(
			albumRepository,
			mediaStore,
			&catalog.TransferMediasFromRepository{TransferMediasRepository: mediaStore},
			&catalog.CoverService{
				CoverRepository:     coverRepository,
				MediaReadRepository: mediaStore,
				Randomiser:          deterministicRandomiser,
			},
			observer,
		)

		err := deleteAlbum.DeleteAlbum(context.Background(), toDeleteAlbumId)
		assert.ErrorIs(t, err, catalog.AlbumNotFoundErr)
		assert.Empty(t, observer.Events)
	})

	t.Run("it should skip firing the event when the repository fails to delete the album after the transfer happened", func(t *testing.T) {
		albumRepository := NewAlbumRepositoryInMemory(allYearAlbum, toDeleteAlbum)
		mediaStore := &MediaReadRepositoryInMemory{
			Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
				toDeleteAlbumId: {movingMedia},
			},
		}
		coverRepository := NewCoverRepositoryInMemory()
		observer := &AlbumDeletedObserverInMemory{}

		deleteAlbum := catalog.NewDeleteAlbum(
			&failingDeleteAlbumInterceptor{AlbumRepositoryInMemory: albumRepository, err: testError},
			mediaStore,
			&catalog.TransferMediasFromRepository{TransferMediasRepository: mediaStore},
			&catalog.CoverService{
				CoverRepository:     coverRepository,
				MediaReadRepository: mediaStore,
				Randomiser:          deterministicRandomiser,
			},
			observer,
		)

		err := deleteAlbum.DeleteAlbum(context.Background(), toDeleteAlbumId)
		assert.ErrorIs(t, err, testError)
		assert.Contains(t, albumRepository.Albums, toDeleteAlbumId, "album should still be in the fake repository since the failing port is a different implementation")
		assert.ElementsMatch(t, []*catalog.MediaMeta{movingMedia}, mediaStore.Medias[allYearAlbumId], "media transfer happened before the deletion attempt")
		assert.Empty(t, observer.Events, "no event was fired")
	})
}

type failingDeleteAlbumInterceptor struct {
	*AlbumRepositoryInMemory
	err error
}

func (i *failingDeleteAlbumInterceptor) DeleteAlbum(_ context.Context, _ catalog.AlbumId) error {
	return i.err
}
