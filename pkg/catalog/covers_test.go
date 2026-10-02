package catalog_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
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

// findAlbumByOwnerPortFake lets tests return a canned list of albums or an error for a
// specific owner, in a fully deterministic order.
type findAlbumByOwnerPortFake struct {
	AlbumsByOwner map[ownermodel.Owner][]*catalog.Album
	Err           error
}

func (f *findAlbumByOwnerPortFake) FindAlbumsByOwner(_ context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.AlbumsByOwner[owner], nil
}

// completeCoversPortFake records each albumId CompleteCovers is called with, and lets a case
// pre-register per-album errors to simulate partial failures.
type completeCoversPortFake struct {
	Completed []catalog.AlbumId
	Errors    map[catalog.AlbumId]error
}

func (c *completeCoversPortFake) CompleteCovers(_ context.Context, albumId catalog.AlbumId) error {
	c.Completed = append(c.Completed, albumId)
	if err, ok := c.Errors[albumId]; ok {
		return err
	}
	return nil
}

func TestBackfillCovers_BackfillForOwner(t *testing.T) {
	owner := ownermodel.Owner("ironman")
	avengersId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers")}
	stealthId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/stealth")}
	avengers := &catalog.Album{AlbumId: avengersId, Name: "Avengers"}
	stealth := &catalog.Album{AlbumId: stealthId, Name: "Stealth"}

	listFailure := errors.New("list exploded")
	completeFailure := errors.New("complete exploded")

	ownerWithTwoAlbums := func() *findAlbumByOwnerPortFake {
		return &findAlbumByOwnerPortFake{AlbumsByOwner: map[ownermodel.Owner][]*catalog.Album{
			owner: {avengers, stealth},
		}}
	}

	type fields struct {
		FindAlbumByOwnerPort *findAlbumByOwnerPortFake
		CompleteCoversPort   *completeCoversPortFake
	}
	type args struct {
		owner ownermodel.Owner
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantReport         catalog.BackfillReport
		wantErr            assert.ErrorAssertionFunc
		expectCompletedIds []catalog.AlbumId
	}{
		{
			name: "it should complete every album of the owner",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				CompleteCoversPort:   &completeCoversPortFake{},
			},
			args:               args{owner: owner},
			wantReport:         catalog.BackfillReport{Albums: 2},
			wantErr:            assert.NoError,
			expectCompletedIds: []catalog.AlbumId{avengersId, stealthId},
		},
		{
			name: "it should return an empty report when the owner has no album",
			fields: fields{
				FindAlbumByOwnerPort: &findAlbumByOwnerPortFake{},
				CompleteCoversPort:   &completeCoversPortFake{},
			},
			args:               args{owner: owner},
			wantReport:         catalog.BackfillReport{Albums: 0},
			wantErr:            assert.NoError,
			expectCompletedIds: nil,
		},
		{
			name: "it should continue after a per-album failure and record it in the report",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				CompleteCoversPort: &completeCoversPortFake{
					Errors: map[catalog.AlbumId]error{avengersId: completeFailure},
				},
			},
			args: args{owner: owner},
			wantReport: catalog.BackfillReport{
				Albums:   2,
				Failures: []catalog.BackfillFailure{{AlbumId: avengersId, Err: completeFailure}},
			},
			wantErr:            assert.NoError,
			expectCompletedIds: []catalog.AlbumId{avengersId, stealthId},
		},
		{
			name: "it should return an error when listing the owner's albums fails",
			fields: fields{
				FindAlbumByOwnerPort: &findAlbumByOwnerPortFake{Err: listFailure},
				CompleteCoversPort:   &completeCoversPortFake{},
			},
			args:       args{owner: owner},
			wantReport: catalog.BackfillReport{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, listFailure)
			},
			expectCompletedIds: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backfill := &catalog.BackfillCovers{
				FindAlbumByOwnerPort: tt.fields.FindAlbumByOwnerPort,
				CompleteCoversPort:   tt.fields.CompleteCoversPort,
			}

			report, err := backfill.BackfillForOwner(context.Background(), tt.args.owner)
			if !tt.wantErr(t, err, fmt.Sprintf("BackfillForOwner(%s)", tt.args.owner)) {
				return
			}
			assert.Equal(t, tt.wantReport, report, "backfill report")
			assert.Equal(t, tt.expectCompletedIds, tt.fields.CompleteCoversPort.Completed, "albums completed")
		})
	}
}
