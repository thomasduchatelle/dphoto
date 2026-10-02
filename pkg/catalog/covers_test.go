package catalog_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func TestCompleteCovers_CompleteCoversFromCandidates(t *testing.T) {
	albumId := catalog.AlbumId{
		Owner:      "ironman",
		FolderName: catalog.NewFolderName("/avengers"),
	}
	image1 := &catalog.MediaMeta{Id: "media-1", Filename: "photo-1.jpg", Type: catalog.MediaTypeImage}
	image2 := &catalog.MediaMeta{Id: "media-2", Filename: "photo-2.jpg", Type: catalog.MediaTypeImage}
	image3 := &catalog.MediaMeta{Id: "media-3", Filename: "photo-3.jpg", Type: catalog.MediaTypeImage}
	image4 := &catalog.MediaMeta{Id: "media-4", Filename: "photo-4.jpg", Type: catalog.MediaTypeImage}
	image5 := &catalog.MediaMeta{Id: "media-5", Filename: "photo-5.jpg", Type: catalog.MediaTypeImage}
	video := &catalog.MediaMeta{Id: "video-1", Filename: "clip-1.mp4", Type: catalog.MediaTypeVideo}
	other := &catalog.MediaMeta{Id: "other-1", Filename: "note-1.txt", Type: catalog.MediaTypeOther}

	// deterministic randomiser: pick the first n indices (in-order), so cases can rely on
	// candidate order to know exactly which medias will be selected.
	deterministicRandomiser := catalog.RandomiserFunc(func(upperBound, n int) []int {
		if n > upperBound {
			n = upperBound
		}
		indices := make([]int, n)
		for i := range indices {
			indices[i] = i
		}
		return indices
	})

	type fields struct {
		CoverRepository *CoverRepositoryInMemory
	}
	type args struct {
		albumId    catalog.AlbumId
		candidates []*catalog.MediaMeta
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		expectSavedCovers []catalog.Cover
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name:   "it should fill 4 empty slots when 4 eligible images are provided",
			fields: fields{CoverRepository: NewCoverRepositoryInMemory()},
			args: args{
				albumId:    albumId,
				candidates: []*catalog.MediaMeta{image1, image2, image3, image4},
			},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name:   "it should cap at 4 covers when more than 4 eligible images are provided",
			fields: fields{CoverRepository: NewCoverRepositoryInMemory()},
			args: args{
				albumId:    albumId,
				candidates: []*catalog.MediaMeta{image1, image2, image3, image4, image5},
			},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should only fill the missing slots, preserving existing covers",
			fields: fields{CoverRepository: NewCoverRepositoryInMemory(coversFor(albumId,
				catalog.Cover{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
				catalog.Cover{MediaId: "media-98", Filename: "random-existing.jpg", Origin: catalog.CoverOriginRandom},
			))},
			args: args{
				albumId:    albumId,
				candidates: []*catalog.MediaMeta{image1, image2, image3},
			},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-98", Filename: "random-existing.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should ignore candidates that are already covers",
			fields: fields{CoverRepository: NewCoverRepositoryInMemory(coversFor(albumId,
				catalog.Cover{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
			))},
			args: args{
				albumId:    albumId,
				candidates: []*catalog.MediaMeta{image1, image2, image3},
			},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name:   "it should skip videos and OTHER medias when selecting covers",
			fields: fields{CoverRepository: NewCoverRepositoryInMemory()},
			args: args{
				albumId:    albumId,
				candidates: []*catalog.MediaMeta{video, other, image1, image2},
			},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should be a no-op when the cover set is already full",
			fields: fields{CoverRepository: NewCoverRepositoryInMemory(coversFor(albumId,
				catalog.Cover{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				catalog.Cover{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				catalog.Cover{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginRandom},
				catalog.Cover{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
			))},
			args: args{
				albumId:    albumId,
				candidates: []*catalog.MediaMeta{image1, image2},
			},
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
			},
			wantErr: assert.NoError,
		},
		{
			name:              "it should be a no-op when no eligible candidate is available",
			fields:            fields{CoverRepository: NewCoverRepositoryInMemory()},
			args:              args{albumId: albumId, candidates: []*catalog.MediaMeta{video, other}},
			expectSavedCovers: nil,
			wantErr:           assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			complete := &catalog.CompleteCovers{
				CoverRepository:     tt.fields.CoverRepository,
				MediaReadRepository: nil,
				Randomiser:          deterministicRandomiser,
			}

			err := complete.CompleteCoversFromCandidates(context.Background(), tt.args.albumId, tt.args.candidates)
			if !tt.wantErr(t, err, fmt.Sprintf("CompleteCoversFromCandidates(%v, %v)", tt.args.albumId, tt.args.candidates)) {
				return
			}
			assert.Equal(t, tt.expectSavedCovers, tt.fields.CoverRepository.Covers[tt.args.albumId], "covers stored for album")
		})
	}
}

func TestCompleteCovers_CompleteCovers_queriesTheAlbumThenCompletes(t *testing.T) {
	albumId := catalog.AlbumId{
		Owner:      "ironman",
		FolderName: catalog.NewFolderName("/avengers"),
	}
	image1 := catalog.MediaMeta{Id: "media-1", Filename: "photo-1.jpg", Type: catalog.MediaTypeImage}
	image2 := catalog.MediaMeta{Id: "media-2", Filename: "photo-2.jpg", Type: catalog.MediaTypeImage}
	video := catalog.MediaMeta{Id: "video-1", Filename: "clip.mp4", Type: catalog.MediaTypeVideo}

	mediaRepo := &MediaReadRepositoryInMemory{
		Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
			albumId: {&image1, &video, &image2},
		},
	}
	coverRepo := NewCoverRepositoryInMemory()

	complete := &catalog.CompleteCovers{
		CoverRepository:     coverRepo,
		MediaReadRepository: mediaRepo,
		Randomiser: catalog.RandomiserFunc(func(upperBound, n int) []int {
			indices := make([]int, n)
			for i := range indices {
				indices[i] = i
			}
			return indices
		}),
	}

	err := complete.CompleteCovers(context.Background(), albumId)
	assert.NoError(t, err)
	assert.Equal(t, []catalog.Cover{
		{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
		{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
	}, coverRepo.Covers[albumId], "images from the album are selected, videos are skipped")
}

func coversFor(albumId catalog.AlbumId, covers ...catalog.Cover) CoverRepositorySeed {
	return CoverRepositorySeed{AlbumId: albumId, Covers: covers}
}
