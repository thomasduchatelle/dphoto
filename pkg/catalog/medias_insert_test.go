package catalog_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func TestInsertMedias_Insert(t *testing.T) {
	owner := ownermodel.Owner("ironman")
	avengersFolder := avengersId.FolderName
	mediaRequest := func(id catalog.MediaId, filename string) catalog.CreateMediaRequest {
		return catalog.CreateMediaRequest{
			Id:         id,
			FolderName: avengersFolder,
			Filename:   filename,
			Type:       catalog.MediaTypeImage,
		}
	}

	mediaRepository := &MediaReadRepositoryInMemory{}
	coverRepository := NewCoverRepositoryInMemory()
	observer := &insertMediasObserverFake{}

	insert := catalog.NewInsertMedias(
		mediaRepository,
		&catalog.CoverService{
			CoverRepository:     coverRepository,
			MediaReadRepository: mediaRepository,
			Randomiser:          deterministicRandomiser,
		},
		observer,
	)

	err := insert.Insert(context.Background(), owner, []catalog.CreateMediaRequest{
		mediaRequest("media-1", "photo-1.jpg"),
		mediaRequest("media-2", "photo-2.jpg"),
	})
	if !assert.NoError(t, err) {
		return
	}

	expectedCovers := []catalog.Cover{
		{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
		{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
	}
	assert.Equal(t, map[catalog.AlbumId][]catalog.Cover{avengersId: expectedCovers}, coverRepository.Covers, "canonical covers persisted")
	assert.Equal(t, []catalog.MediasInserted{
		{
			Inserted: map[catalog.AlbumId][]catalog.MediaId{avengersId: {"media-1", "media-2"}},
			Covers:   map[catalog.AlbumId][]catalog.Cover{avengersId: expectedCovers},
		},
	}, observer.Events, "observer notified with the resulting event")
}

type insertMediasObserverFake struct {
	Events []catalog.MediasInserted
}

func (o *insertMediasObserverFake) OnMediasInserted(_ context.Context, event catalog.MediasInserted) error {
	o.Events = append(o.Events, event)
	return nil
}
