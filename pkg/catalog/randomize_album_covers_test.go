package catalog_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func TestRandomizeAlbumCovers_Randomize(t *testing.T) {
	cherryPicked := catalog.Cover{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked}
	stale := catalog.Cover{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom}

	randomiserWithOffset := func(offset int) catalog.RandomiserFunc {
		return catalog.RandomiserFunc(func(upperBound, n int) []int {
			if n > upperBound {
				n = upperBound
			}
			indices := make([]int, n)
			for i := range indices {
				indices[i] = (i + offset) % upperBound
			}
			return indices
		})
	}

	type fields struct {
		CoverRepository *CoverRepositoryInMemory
		MediaRepository *MediaReadRepositoryInMemory
		Randomiser      catalog.Randomiser
	}
	tests := []struct {
		name              string
		fields            fields
		args              catalog.AlbumId
		wantCovers        []catalog.Cover
		expectSavedCovers map[catalog.AlbumId][]catalog.Cover
		expectEvents      []catalog.AlbumCoversRandomised
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "it should drop RANDOM covers, keep CHERRY_PICKED, fill empties, persist the new set, and emit the event",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(avengersId, cherryPicked, stale)),
				MediaRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{avengersId: {image1, image2, image3, image4}},
				},
				Randomiser: deterministicRandomiser,
			},
			args: avengersId,
			wantCovers: []catalog.Cover{
				cherryPicked,
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					cherryPicked,
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectEvents: []catalog.AlbumCoversRandomised{
				{
					AlbumId: avengersId,
					Covers: []catalog.Cover{
						cherryPicked,
						{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
						{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
						{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return the current covers without firing the event when the redraw produced the same set",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(avengersId, cherryPicked)),
				MediaRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{avengersId: {image1}},
				},
				Randomiser: deterministicRandomiser,
			},
			args:              avengersId,
			wantCovers:        []catalog.Cover{cherryPicked},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{avengersId: {cherryPicked}},
			expectEvents:      nil,
			wantErr:           assert.NoError,
		},
		{
			name: "it should return the new covers from the event when the set changed",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(avengersId, stale)),
				MediaRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{avengersId: {image3, image4}},
				},
				Randomiser: randomiserWithOffset(0),
			},
			args: avengersId,
			wantCovers: []catalog.Cover{
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectEvents: []catalog.AlbumCoversRandomised{
				{
					AlbumId: avengersId,
					Covers: []catalog.Cover{
						{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
						{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := &albumCoversRandomisedObserverFake{}
			useCase := catalog.NewRandomizeAlbumCovers(
				&catalog.CoverService{
					CoverRepository:     tt.fields.CoverRepository,
					MediaReadRepository: tt.fields.MediaRepository,
					Randomiser:          tt.fields.Randomiser,
				},
				tt.fields.CoverRepository,
				observer,
			)

			got, err := useCase.Randomize(context.Background(), tt.args)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.wantCovers, got, "returned covers")
			assert.Equal(t, tt.expectSavedCovers, tt.fields.CoverRepository.Covers, "canonical covers persisted")
			assert.Equal(t, tt.expectEvents, observer.Events, "observer notifications")
		})
	}
}

func TestRandomizeAlbumCovers_Randomize_propagatesErrors(t *testing.T) {
	wantErr := errors.New("boom")
	useCase := catalog.NewRandomizeAlbumCovers(
		coverServicePortFake{err: wantErr},
		nil,
	)

	_, err := useCase.Randomize(context.Background(), avengersId)
	assert.ErrorIs(t, err, wantErr)
}

type albumCoversRandomisedObserverFake struct {
	Events []catalog.AlbumCoversRandomised
	Err    error
}

func (o *albumCoversRandomisedObserverFake) OnAlbumCoversRandomised(_ context.Context, event catalog.AlbumCoversRandomised) error {
	o.Events = append(o.Events, event)
	return o.Err
}

type coverServicePortFake struct {
	err error
}

func (f coverServicePortFake) Randomise(_ context.Context, _ bool, _ ...catalog.AlbumId) (map[catalog.AlbumId][]catalog.Cover, error) {
	return nil, f.err
}

func (f coverServicePortFake) StableRefresh(_ context.Context, _ catalog.TransferredMedias) (map[catalog.AlbumId][]catalog.Cover, error) {
	return nil, f.err
}
