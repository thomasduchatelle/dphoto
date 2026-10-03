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

var deterministicRandomiser = catalog.RandomiserFunc(func(upperBound, n int) []int {
	if n > upperBound {
		n = upperBound
	}
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	return indices
})

var (
	coverAlbumId = catalog.AlbumId{
		Owner:      "ironman",
		FolderName: catalog.NewFolderName("/avengers"),
	}
	image1 = &catalog.MediaMeta{Id: "media-1", Filename: "photo-1.jpg", Type: catalog.MediaTypeImage}
	image2 = &catalog.MediaMeta{Id: "media-2", Filename: "photo-2.jpg", Type: catalog.MediaTypeImage}
	image3 = &catalog.MediaMeta{Id: "media-3", Filename: "photo-3.jpg", Type: catalog.MediaTypeImage}
	image4 = &catalog.MediaMeta{Id: "media-4", Filename: "photo-4.jpg", Type: catalog.MediaTypeImage}
	image5 = &catalog.MediaMeta{Id: "media-5", Filename: "photo-5.jpg", Type: catalog.MediaTypeImage}
	video1 = &catalog.MediaMeta{Id: "video-1", Filename: "clip-1.mp4", Type: catalog.MediaTypeVideo}
	other1 = &catalog.MediaMeta{Id: "other-1", Filename: "note-1.txt", Type: catalog.MediaTypeOther}
)

func TestCoverMaintenance_Refresh(t *testing.T) {
	type fields struct {
		CoverRepository     *CoverRepositoryInMemory
		MediaReadRepository *MediaReadRepositoryInMemory
	}
	type args struct {
		removed []catalog.MediaId
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantCovers        []catalog.Cover
		wantChanged       bool
		expectSavedCovers []catalog.Cover
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "it should fill empty slots from the album's full image set",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {image1, image2, image3, image4}},
				},
			},
			args: args{},
			wantCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantChanged: true,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should cap at MaxCoversPerAlbum when more images are available",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {image1, image2, image3, image4, image5}},
				},
			},
			args: args{},
			wantCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantChanged: true,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should drop existing RANDOM covers and redraw them",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-98", Filename: "stale.jpg", Origin: catalog.CoverOriginRandom},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {image1, image2, image3, image4}},
				},
			},
			args: args{},
			wantCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantChanged: true,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should keep CHERRY_PICKED covers untouched and refill the rest",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-98", Filename: "random-existing.jpg", Origin: catalog.CoverOriginRandom},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {image1, image2, image3}},
				},
			},
			args: args{},
			wantCovers: []catalog.Cover{
				{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantChanged: true,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should strip a CHERRY_PICKED cover whose media is in the removed set",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-98", Filename: "stay.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args: args{removed: []catalog.MediaId{"media-99"}},
			wantCovers: []catalog.Cover{
				{MediaId: "media-98", Filename: "stay.jpg", Origin: catalog.CoverOriginCherryPicked},
			},
			wantChanged: true,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-98", Filename: "stay.jpg", Origin: catalog.CoverOriginCherryPicked},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should skip videos and OTHER medias when drawing from the album",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {video1, other1, image1, image2}},
				},
			},
			args: args{},
			wantCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantChanged: true,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should be a no-op when the kept set is already full of CHERRY_PICKED covers",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args: args{},
			wantCovers: []catalog.Cover{
				{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
			},
			wantChanged: false,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should produce an empty set when no eligible candidate is available",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {video1, other1}},
				},
			},
			args:              args{},
			wantCovers:        nil,
			wantChanged:       false,
			expectSavedCovers: nil,
			wantErr:           assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maintenance := &catalog.CoverMaintenance{
				CoverRepository:     tt.fields.CoverRepository,
				MediaReadRepository: tt.fields.MediaReadRepository,
				Randomiser:          deterministicRandomiser,
			}

			covers, changed, err := maintenance.Refresh(context.Background(), coverAlbumId, tt.args.removed)
			if !tt.wantErr(t, err, fmt.Sprintf("Refresh(%v, removed=%v)", coverAlbumId, tt.args.removed)) {
				return
			}
			assert.Equal(t, tt.wantCovers, covers, "returned covers")
			assert.Equal(t, tt.wantChanged, changed, "changed flag")
			assert.Equal(t, tt.expectSavedCovers, tt.fields.CoverRepository.Covers[coverAlbumId], "covers stored for album")
		})
	}
}

func TestCoverMaintenance_Stabilise(t *testing.T) {
	type fields struct {
		CoverRepository     *CoverRepositoryInMemory
		MediaReadRepository *MediaReadRepositoryInMemory
	}
	type args struct {
		removed []catalog.MediaId
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantCovers        []catalog.Cover
		wantChanged       bool
		expectSavedCovers []catalog.Cover
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "it should fill empty slots from the album's full image set on an empty cover set",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {image1, image2, image3, image4}},
				},
			},
			args: args{},
			wantCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantChanged: true,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should preserve existing RANDOM covers (not redraw them)",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {image1, image2, image3, image4}},
				},
			},
			args: args{},
			wantCovers: []catalog.Cover{
				{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantChanged: true,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should preserve CHERRY_PICKED covers and refill the rest",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {image1, image2, image3}},
				},
			},
			args: args{},
			wantCovers: []catalog.Cover{
				{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantChanged: true,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-99", Filename: "starred.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should strip the covers whose medias are in the removed set and refill from the album",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {image1, image2}},
				},
			},
			args: args{removed: []catalog.MediaId{"media-91", "media-92"}},
			wantCovers: []catalog.Cover{
				{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantChanged: true,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
				{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should be a no-op when the cover set is already full and nothing is removed",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {image1, image2, image3, image4, image5}},
				},
			},
			args: args{},
			wantCovers: []catalog.Cover{
				{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
			},
			wantChanged: false,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should be a no-op when removed matches no cover and the set is already full",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args: args{removed: []catalog.MediaId{"not-a-cover"}},
			wantCovers: []catalog.Cover{
				{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
			},
			wantChanged: false,
			expectSavedCovers: []catalog.Cover{
				{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-92", Filename: "b.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-93", Filename: "c.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: "media-94", Filename: "d.jpg", Origin: catalog.CoverOriginCherryPicked},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should leave empty slots empty when no eligible candidate is available after stripping",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(coverAlbumId,
					catalog.Cover{MediaId: "media-91", Filename: "a.jpg", Origin: catalog.CoverOriginRandom},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{coverAlbumId: {video1, other1}},
				},
			},
			args:              args{removed: []catalog.MediaId{"media-91"}},
			wantCovers:        nil,
			wantChanged:       true,
			expectSavedCovers: nil,
			wantErr:           assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maintenance := &catalog.CoverMaintenance{
				CoverRepository:     tt.fields.CoverRepository,
				MediaReadRepository: tt.fields.MediaReadRepository,
				Randomiser:          deterministicRandomiser,
			}

			covers, changed, err := maintenance.Stabilise(context.Background(), coverAlbumId, tt.args.removed)
			if !tt.wantErr(t, err, fmt.Sprintf("Stabilise(%v, removed=%v)", coverAlbumId, tt.args.removed)) {
				return
			}
			assert.Equal(t, tt.wantCovers, covers, "returned covers")
			assert.Equal(t, tt.wantChanged, changed, "changed flag")
			assert.Equal(t, tt.expectSavedCovers, tt.fields.CoverRepository.Covers[coverAlbumId], "covers stored for album")
		})
	}
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

// refreshCoversPortFake records every (albumId, removed) tuple Refresh is called with,
// and lets a case pre-register per-album errors or per-album "unchanged" responses to
// simulate partial failures and no-ops.
type refreshCoversPortFake struct {
	Refreshed []refreshCoversCall
	Errors    map[catalog.AlbumId]error
	Unchanged map[catalog.AlbumId]bool
	Covers    map[catalog.AlbumId][]catalog.Cover
}

type refreshCoversCall struct {
	AlbumId catalog.AlbumId
	Removed []catalog.MediaId
}

func (r *refreshCoversPortFake) Refresh(_ context.Context, albumId catalog.AlbumId, removed []catalog.MediaId) ([]catalog.Cover, bool, error) {
	r.Refreshed = append(r.Refreshed, refreshCoversCall{AlbumId: albumId, Removed: removed})
	if err, ok := r.Errors[albumId]; ok {
		return nil, false, err
	}
	if r.Unchanged[albumId] {
		return nil, false, nil
	}
	return r.Covers[albumId], true, nil
}

type backfillCoversViewUpdaterFake struct {
	Updates map[catalog.AlbumId][]catalog.Cover
	Errors  map[catalog.AlbumId]error
}

func (b *backfillCoversViewUpdaterFake) UpdateCovers(_ context.Context, albumId catalog.AlbumId, covers []catalog.Cover) error {
	if err, ok := b.Errors[albumId]; ok {
		return err
	}
	if b.Updates == nil {
		b.Updates = make(map[catalog.AlbumId][]catalog.Cover)
	}
	b.Updates[albumId] = covers
	return nil
}

func TestBackfillCovers_BackfillForOwner(t *testing.T) {
	owner := ownermodel.Owner("ironman")
	avengersId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers")}
	stealthId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/stealth")}
	avengers := &catalog.Album{AlbumId: avengersId, Name: "Avengers"}
	stealth := &catalog.Album{AlbumId: stealthId, Name: "Stealth"}
	avengersCovers := []catalog.Cover{{MediaId: "m-a", Filename: "a.jpg", Origin: catalog.CoverOriginRandom}}
	stealthCovers := []catalog.Cover{{MediaId: "m-s", Filename: "s.jpg", Origin: catalog.CoverOriginRandom}}

	listFailure := errors.New("list exploded")
	refreshFailure := errors.New("refresh exploded")
	viewFailure := errors.New("view update exploded")

	ownerWithTwoAlbums := func() *findAlbumByOwnerPortFake {
		return &findAlbumByOwnerPortFake{AlbumsByOwner: map[ownermodel.Owner][]*catalog.Album{
			owner: {avengers, stealth},
		}}
	}

	type fields struct {
		FindAlbumByOwnerPort      *findAlbumByOwnerPortFake
		RefreshCoversPort         *refreshCoversPortFake
		BackfillCoversViewUpdater *backfillCoversViewUpdaterFake
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
		expectRefreshedIds []catalog.AlbumId
		expectViewUpdates  map[catalog.AlbumId][]catalog.Cover
	}{
		{
			name: "it should refresh every album of the owner and fan the new covers to the view",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				RefreshCoversPort: &refreshCoversPortFake{
					Covers: map[catalog.AlbumId][]catalog.Cover{avengersId: avengersCovers, stealthId: stealthCovers},
				},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
			},
			args:               args{owner: owner},
			wantReport:         catalog.BackfillReport{Albums: 2},
			wantErr:            assert.NoError,
			expectRefreshedIds: []catalog.AlbumId{avengersId, stealthId},
			expectViewUpdates:  map[catalog.AlbumId][]catalog.Cover{avengersId: avengersCovers, stealthId: stealthCovers},
		},
		{
			name: "it should not update the view when Refresh reports no change",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				RefreshCoversPort: &refreshCoversPortFake{
					Covers:    map[catalog.AlbumId][]catalog.Cover{stealthId: stealthCovers},
					Unchanged: map[catalog.AlbumId]bool{avengersId: true},
				},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
			},
			args:               args{owner: owner},
			wantReport:         catalog.BackfillReport{Albums: 2},
			wantErr:            assert.NoError,
			expectRefreshedIds: []catalog.AlbumId{avengersId, stealthId},
			expectViewUpdates:  map[catalog.AlbumId][]catalog.Cover{stealthId: stealthCovers},
		},
		{
			name: "it should return an empty report when the owner has no album",
			fields: fields{
				FindAlbumByOwnerPort:      &findAlbumByOwnerPortFake{},
				RefreshCoversPort:         &refreshCoversPortFake{},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
			},
			args:               args{owner: owner},
			wantReport:         catalog.BackfillReport{Albums: 0},
			wantErr:            assert.NoError,
			expectRefreshedIds: nil,
			expectViewUpdates:  nil,
		},
		{
			name: "it should continue after a per-album Refresh failure and skip the view update for that album",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				RefreshCoversPort: &refreshCoversPortFake{
					Errors: map[catalog.AlbumId]error{avengersId: refreshFailure},
					Covers: map[catalog.AlbumId][]catalog.Cover{stealthId: stealthCovers},
				},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
			},
			args: args{owner: owner},
			wantReport: catalog.BackfillReport{
				Albums:   2,
				Failures: []catalog.BackfillFailure{{AlbumId: avengersId, Err: refreshFailure}},
			},
			wantErr:            assert.NoError,
			expectRefreshedIds: []catalog.AlbumId{avengersId, stealthId},
			expectViewUpdates:  map[catalog.AlbumId][]catalog.Cover{stealthId: stealthCovers},
		},
		{
			name: "it should continue after a per-album view update failure and record it in the report",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				RefreshCoversPort: &refreshCoversPortFake{
					Covers: map[catalog.AlbumId][]catalog.Cover{avengersId: avengersCovers, stealthId: stealthCovers},
				},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{
					Errors: map[catalog.AlbumId]error{avengersId: viewFailure},
				},
			},
			args: args{owner: owner},
			wantReport: catalog.BackfillReport{
				Albums:   2,
				Failures: []catalog.BackfillFailure{{AlbumId: avengersId, Err: viewFailure}},
			},
			wantErr:            assert.NoError,
			expectRefreshedIds: []catalog.AlbumId{avengersId, stealthId},
			expectViewUpdates:  map[catalog.AlbumId][]catalog.Cover{stealthId: stealthCovers},
		},
		{
			name: "it should return an error when listing the owner's albums fails",
			fields: fields{
				FindAlbumByOwnerPort:      &findAlbumByOwnerPortFake{Err: listFailure},
				RefreshCoversPort:         &refreshCoversPortFake{},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
			},
			args:       args{owner: owner},
			wantReport: catalog.BackfillReport{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, listFailure)
			},
			expectRefreshedIds: nil,
			expectViewUpdates:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backfill := &catalog.BackfillCovers{
				FindAlbumByOwnerPort:      tt.fields.FindAlbumByOwnerPort,
				RefreshCoversPort:         tt.fields.RefreshCoversPort,
				BackfillCoversViewUpdater: tt.fields.BackfillCoversViewUpdater,
			}

			report, err := backfill.BackfillForOwner(context.Background(), tt.args.owner)
			if !tt.wantErr(t, err, fmt.Sprintf("BackfillForOwner(%s)", tt.args.owner)) {
				return
			}
			assert.Equal(t, tt.wantReport, report, "backfill report")

			var refreshedIds []catalog.AlbumId
			for _, call := range tt.fields.RefreshCoversPort.Refreshed {
				refreshedIds = append(refreshedIds, call.AlbumId)
			}
			assert.Equal(t, tt.expectRefreshedIds, refreshedIds, "albums refreshed")
			assert.Equal(t, tt.expectViewUpdates, tt.fields.BackfillCoversViewUpdater.Updates, "view updates")
		})
	}
}
