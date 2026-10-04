package catalog_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func TestBackfillCovers_BackfillForOwner(t *testing.T) {
	owner := ownermodel.Owner("ironman")
	avengersAlbum := &catalog.Album{AlbumId: avengersId, Name: "Avengers"}
	stealthAlbum := &catalog.Album{AlbumId: stealthId, Name: "Stealth"}

	albumRepository := NewAlbumRepositoryInMemory(avengersAlbum, stealthAlbum)
	coverRepository := NewCoverRepositoryInMemory(coversFor(stealthId,
		catalog.Cover{MediaId: "media-5", Filename: "photo-5.jpg", Origin: catalog.CoverOriginRandom},
	))
	mediaRepository := &MediaReadRepositoryInMemory{
		Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
			avengersId: {image1, image2},
			stealthId:  {image5},
		},
	}
	observer := &coverBackfillObserverFake{}

	backfill := &catalog.BackfillCovers{
		FindAlbumByOwnerPort: albumRepository,
		CoverService: &catalog.CoverService{
			CoverRepository:     coverRepository,
			MediaReadRepository: mediaRepository,
			Randomiser:          deterministicRandomiser,
		},
		Observers: []catalog.CoverBackfillObserver{observer},
	}

	got, err := backfill.BackfillForOwner(context.Background(), owner)
	if !assert.NoError(t, err) {
		return
	}

	expectedChanged := map[catalog.AlbumId][]catalog.Cover{
		avengersId: {
			{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
			{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
		},
	}

	assert.Equal(t, expectedChanged, got, "returned changed map")
	assert.Equal(t, map[catalog.AlbumId][]catalog.Cover{
		avengersId: {
			{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
			{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
		},
		stealthId: {
			{MediaId: "media-5", Filename: "photo-5.jpg", Origin: catalog.CoverOriginRandom},
		},
	}, coverRepository.Covers, "canonical covers persisted")
	assert.Equal(t, []map[catalog.AlbumId][]catalog.Cover{expectedChanged}, observer.Notifications, "observer notified with the resulting map")
}

type coverBackfillObserverFake struct {
	Notifications []map[catalog.AlbumId][]catalog.Cover
	Err           error
}

func (o *coverBackfillObserverFake) OnCoverBackfilled(_ context.Context, coversByAlbumId map[catalog.AlbumId][]catalog.Cover) error {
	o.Notifications = append(o.Notifications, coversByAlbumId)
	return o.Err
}
