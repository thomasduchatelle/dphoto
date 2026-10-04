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
	avengersId = catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/avengers")}
	stealthId  = catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/stealth")}
	image1     = &catalog.MediaMeta{Id: "media-1", Filename: "photo-1.jpg", Type: catalog.MediaTypeImage}
	image2     = &catalog.MediaMeta{Id: "media-2", Filename: "photo-2.jpg", Type: catalog.MediaTypeImage}
	image3     = &catalog.MediaMeta{Id: "media-3", Filename: "photo-3.jpg", Type: catalog.MediaTypeImage}
	image4     = &catalog.MediaMeta{Id: "media-4", Filename: "photo-4.jpg", Type: catalog.MediaTypeImage}
	image5     = &catalog.MediaMeta{Id: "media-5", Filename: "photo-5.jpg", Type: catalog.MediaTypeImage}
	image6     = &catalog.MediaMeta{Id: "media-6", Filename: "photo-6.jpg", Type: catalog.MediaTypeImage}
	video1     = &catalog.MediaMeta{Id: "video-1", Filename: "clip-1.mp4", Type: catalog.MediaTypeVideo}
	other1     = &catalog.MediaMeta{Id: "other-1", Filename: "note-1.txt", Type: catalog.MediaTypeOther}
)

func TestCoverService_Randomise(t *testing.T) {
	type fields struct {
		CoverRepository     *CoverRepositoryInMemory
		MediaReadRepository *MediaReadRepositoryInMemory
	}
	type args struct {
		albumIds []catalog.AlbumId
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantChanged       map[catalog.AlbumId][]catalog.Cover
		expectSavedCovers map[catalog.AlbumId][]catalog.Cover
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "it should return a nil map when no album is requested",
			fields: fields{
				CoverRepository:     NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args:              args{albumIds: nil},
			wantChanged:       nil,
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{},
			wantErr:           assert.NoError,
		},
		{
			name: "it should fill empty slots from the album's full image set",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{avengersId: {image1, image2, image3, image4}},
				},
			},
			args: args{albumIds: []catalog.AlbumId{avengersId}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should drop existing RANDOM covers and redraw them from the album",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(avengersId,
					catalog.Cover{MediaId: "media-98", Filename: "stale.jpg", Origin: catalog.CoverOriginRandom},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{avengersId: {image1, image2, image3, image4}},
				},
			},
			args: args{albumIds: []catalog.AlbumId{avengersId}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should keep CHERRY_PICKED covers whose media is still in the album",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(avengersId,
					catalog.Cover{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-98", Filename: "random-existing.jpg", Origin: catalog.CoverOriginRandom},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{avengersId: {image1, image2, image3}},
				},
			},
			args: args{albumIds: []catalog.AlbumId{avengersId}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should drop CHERRY_PICKED covers whose media is no longer in the album",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(avengersId,
					catalog.Cover{MediaId: "orphan", Filename: "gone.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{avengersId: {image1, image2}},
				},
			},
			args: args{albumIds: []catalog.AlbumId{avengersId}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should skip videos and OTHER medias when drawing from the album",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{avengersId: {video1, other1, image1, image2}},
				},
			},
			args: args{albumIds: []catalog.AlbumId{avengersId}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should omit an album from the result map when its cover set is unchanged",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(avengersId,
					catalog.Cover{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginCherryPicked},
					catalog.Cover{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{avengersId: {image1, image2, image3, image4, image5}},
				},
			},
			args:        args{albumIds: []catalog.AlbumId{avengersId}},
			wantChanged: nil,
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginCherryPicked},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should process every album in a single call and return only the changed ones",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(stealthId,
					catalog.Cover{MediaId: "media-5", Filename: "photo-5.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
						avengersId: {image1, image2},
						stealthId:  {image5},
					},
				},
			},
			args: args{albumIds: []catalog.AlbumId{avengersId, stealthId}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				},
				stealthId: {
					{MediaId: "media-5", Filename: "photo-5.jpg", Origin: catalog.CoverOriginCherryPicked},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &catalog.CoverService{
				CoverRepository:     tt.fields.CoverRepository,
				MediaReadRepository: tt.fields.MediaReadRepository,
				Randomiser:          deterministicRandomiser,
			}

			got, err := service.Randomise(context.Background(), tt.args.albumIds...)
			if !tt.wantErr(t, err, fmt.Sprintf("Randomise(%v)", tt.args.albumIds)) {
				return
			}
			assert.Equal(t, tt.wantChanged, got, "returned changed map")
			assert.Equal(t, tt.expectSavedCovers, tt.fields.CoverRepository.Covers, "covers stored")
		})
	}
}

func TestCoverService_StableRefresh(t *testing.T) {
	type fields struct {
		CoverRepository     *CoverRepositoryInMemory
		MediaReadRepository *MediaReadRepositoryInMemory
	}
	type args struct {
		transferred catalog.TransferredMedias
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantChanged       map[catalog.AlbumId][]catalog.Cover
		expectSavedCovers map[catalog.AlbumId][]catalog.Cover
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "it should return a nil map when the transfer is empty",
			fields: fields{
				CoverRepository:     NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{},
			},
			args:              args{transferred: catalog.NewTransferredMedias()},
			wantChanged:       nil,
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{},
			wantErr:           assert.NoError,
		},
		{
			name: "it should be a no-op on a source whose covers were not touched by the transfer",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(avengersId,
					catalog.Cover{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
						avengersId: {image1, image2, image3, image4, image5},
						stealthId:  {image6},
					},
				},
			},
			args: args{transferred: catalog.TransferredMedias{
				Transfers:  map[catalog.AlbumId][]catalog.MediaId{stealthId: {"media-99"}},
				FromAlbums: []catalog.AlbumId{avengersId},
			}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				stealthId: {
					{MediaId: "media-6", Filename: "photo-6.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
				stealthId: {
					{MediaId: "media-6", Filename: "photo-6.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should strip source covers whose media moved and backfill the empty slots",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(avengersId,
					catalog.Cover{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					catalog.Cover{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
						avengersId: {image2, image3, image4},
						stealthId:  {image1},
					},
				},
			},
			args: args{transferred: catalog.TransferredMedias{
				Transfers:  map[catalog.AlbumId][]catalog.MediaId{stealthId: {"media-1"}},
				FromAlbums: []catalog.AlbumId{avengersId},
			}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
				stealthId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
				stealthId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should carry CHERRY_PICKED covers to the destination album as CHERRY_PICKED",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(coversFor(avengersId,
					catalog.Cover{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
				)),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
						avengersId: {image2, image3},
						stealthId:  {image1, image4},
					},
				},
			},
			args: args{transferred: catalog.TransferredMedias{
				Transfers:  map[catalog.AlbumId][]catalog.MediaId{stealthId: {"media-1"}},
				FromAlbums: []catalog.AlbumId{avengersId},
			}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				},
				stealthId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
				},
				stealthId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should displace the oldest RANDOM cover on the destination to make room for an inherited CHERRY_PICKED",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(
					coversFor(avengersId, catalog.Cover{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked}),
					coversFor(stealthId,
						catalog.Cover{MediaId: "media-10", Filename: "stealth-10.jpg", Origin: catalog.CoverOriginRandom},
						catalog.Cover{MediaId: "media-11", Filename: "stealth-11.jpg", Origin: catalog.CoverOriginCherryPicked},
						catalog.Cover{MediaId: "media-12", Filename: "stealth-12.jpg", Origin: catalog.CoverOriginRandom},
						catalog.Cover{MediaId: "media-13", Filename: "stealth-13.jpg", Origin: catalog.CoverOriginCherryPicked},
					),
				),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
						avengersId: {image2},
						stealthId: {
							image1,
							{Id: "media-10", Filename: "stealth-10.jpg", Type: catalog.MediaTypeImage},
							{Id: "media-11", Filename: "stealth-11.jpg", Type: catalog.MediaTypeImage},
							{Id: "media-12", Filename: "stealth-12.jpg", Type: catalog.MediaTypeImage},
							{Id: "media-13", Filename: "stealth-13.jpg", Type: catalog.MediaTypeImage},
						},
					},
				},
			},
			args: args{transferred: catalog.TransferredMedias{
				Transfers:  map[catalog.AlbumId][]catalog.MediaId{stealthId: {"media-1"}},
				FromAlbums: []catalog.AlbumId{avengersId},
			}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				},
				stealthId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-11", Filename: "stealth-11.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-12", Filename: "stealth-12.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-13", Filename: "stealth-13.jpg", Origin: catalog.CoverOriginCherryPicked},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				},
				stealthId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-11", Filename: "stealth-11.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-12", Filename: "stealth-12.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-13", Filename: "stealth-13.jpg", Origin: catalog.CoverOriginCherryPicked},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should silently drop an inherited CHERRY_PICKED when the destination is full of CHERRY_PICKED covers",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(
					coversFor(avengersId, catalog.Cover{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginCherryPicked}),
					coversFor(stealthId,
						catalog.Cover{MediaId: "media-10", Filename: "stealth-10.jpg", Origin: catalog.CoverOriginCherryPicked},
						catalog.Cover{MediaId: "media-11", Filename: "stealth-11.jpg", Origin: catalog.CoverOriginCherryPicked},
						catalog.Cover{MediaId: "media-12", Filename: "stealth-12.jpg", Origin: catalog.CoverOriginCherryPicked},
						catalog.Cover{MediaId: "media-13", Filename: "stealth-13.jpg", Origin: catalog.CoverOriginCherryPicked},
					),
				),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
						avengersId: {image2},
						stealthId: {
							image1,
							{Id: "media-10", Filename: "stealth-10.jpg", Type: catalog.MediaTypeImage},
							{Id: "media-11", Filename: "stealth-11.jpg", Type: catalog.MediaTypeImage},
							{Id: "media-12", Filename: "stealth-12.jpg", Type: catalog.MediaTypeImage},
							{Id: "media-13", Filename: "stealth-13.jpg", Type: catalog.MediaTypeImage},
						},
					},
				},
			},
			args: args{transferred: catalog.TransferredMedias{
				Transfers:  map[catalog.AlbumId][]catalog.MediaId{stealthId: {"media-1"}},
				FromAlbums: []catalog.AlbumId{avengersId},
			}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				avengersId: {
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
				},
				stealthId: {
					{MediaId: "media-10", Filename: "stealth-10.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-11", Filename: "stealth-11.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-12", Filename: "stealth-12.jpg", Origin: catalog.CoverOriginCherryPicked},
					{MediaId: "media-13", Filename: "stealth-13.jpg", Origin: catalog.CoverOriginCherryPicked},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should fill empty slots on a destination album with no existing covers",
			fields: fields{
				CoverRepository: NewCoverRepositoryInMemory(),
				MediaReadRepository: &MediaReadRepositoryInMemory{
					Medias: map[catalog.AlbumId][]*catalog.MediaMeta{
						stealthId: {image1, image2, image3, image4, image5},
					},
				},
			},
			args: args{transferred: catalog.TransferredMedias{
				Transfers: map[catalog.AlbumId][]catalog.MediaId{stealthId: {"media-1", "media-2"}},
			}},
			wantChanged: map[catalog.AlbumId][]catalog.Cover{
				stealthId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			expectSavedCovers: map[catalog.AlbumId][]catalog.Cover{
				stealthId: {
					{MediaId: "media-1", Filename: "photo-1.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-2", Filename: "photo-2.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-3", Filename: "photo-3.jpg", Origin: catalog.CoverOriginRandom},
					{MediaId: "media-4", Filename: "photo-4.jpg", Origin: catalog.CoverOriginRandom},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &catalog.CoverService{
				CoverRepository:     tt.fields.CoverRepository,
				MediaReadRepository: tt.fields.MediaReadRepository,
				Randomiser:          deterministicRandomiser,
			}

			got, err := service.StableRefresh(context.Background(), tt.args.transferred)
			if !tt.wantErr(t, err, fmt.Sprintf("StableRefresh(%v)", tt.args.transferred)) {
				return
			}
			assert.Equal(t, tt.wantChanged, got, "returned changed map")
			assert.Equal(t, tt.expectSavedCovers, tt.fields.CoverRepository.Covers, "covers stored")
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

// randomiseCoversPortFake records every albumId Randomise is called with, and lets a
// case pre-register per-album errors or per-album canned cover sets to simulate
// partial failures and no-ops.
type randomiseCoversPortFake struct {
	Calls   []catalog.AlbumId
	Errors  map[catalog.AlbumId]error
	Changed map[catalog.AlbumId][]catalog.Cover
}

func (r *randomiseCoversPortFake) Randomise(_ context.Context, albumIds ...catalog.AlbumId) (map[catalog.AlbumId][]catalog.Cover, error) {
	for _, albumId := range albumIds {
		r.Calls = append(r.Calls, albumId)
		if err, ok := r.Errors[albumId]; ok {
			return nil, err
		}
	}
	result := make(map[catalog.AlbumId][]catalog.Cover)
	for _, albumId := range albumIds {
		if covers, ok := r.Changed[albumId]; ok {
			result[albumId] = covers
		}
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
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
	backfillAvengersId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers")}
	backfillStealthId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/stealth")}
	avengers := &catalog.Album{AlbumId: backfillAvengersId, Name: "Avengers"}
	stealth := &catalog.Album{AlbumId: backfillStealthId, Name: "Stealth"}
	avengersCovers := []catalog.Cover{{MediaId: "m-a", Filename: "a.jpg", Origin: catalog.CoverOriginRandom}}
	stealthCovers := []catalog.Cover{{MediaId: "m-s", Filename: "s.jpg", Origin: catalog.CoverOriginRandom}}

	listFailure := errors.New("list exploded")
	randomiseFailure := errors.New("randomise exploded")
	viewFailure := errors.New("view update exploded")

	ownerWithTwoAlbums := func() *findAlbumByOwnerPortFake {
		return &findAlbumByOwnerPortFake{AlbumsByOwner: map[ownermodel.Owner][]*catalog.Album{
			owner: {avengers, stealth},
		}}
	}

	type fields struct {
		FindAlbumByOwnerPort      *findAlbumByOwnerPortFake
		RandomiseCoversPort       *randomiseCoversPortFake
		BackfillCoversViewUpdater *backfillCoversViewUpdaterFake
	}
	type args struct {
		owner ownermodel.Owner
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantReport        catalog.BackfillReport
		wantErr           assert.ErrorAssertionFunc
		expectCalledIds   []catalog.AlbumId
		expectViewUpdates map[catalog.AlbumId][]catalog.Cover
	}{
		{
			name: "it should randomise every album of the owner and fan the new covers to the view",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				RandomiseCoversPort: &randomiseCoversPortFake{
					Changed: map[catalog.AlbumId][]catalog.Cover{backfillAvengersId: avengersCovers, backfillStealthId: stealthCovers},
				},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
			},
			args:              args{owner: owner},
			wantReport:        catalog.BackfillReport{Albums: 2},
			wantErr:           assert.NoError,
			expectCalledIds:   []catalog.AlbumId{backfillAvengersId, backfillStealthId},
			expectViewUpdates: map[catalog.AlbumId][]catalog.Cover{backfillAvengersId: avengersCovers, backfillStealthId: stealthCovers},
		},
		{
			name: "it should not update the view when Randomise reports no change",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				RandomiseCoversPort: &randomiseCoversPortFake{
					Changed: map[catalog.AlbumId][]catalog.Cover{backfillStealthId: stealthCovers},
				},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
			},
			args:              args{owner: owner},
			wantReport:        catalog.BackfillReport{Albums: 2},
			wantErr:           assert.NoError,
			expectCalledIds:   []catalog.AlbumId{backfillAvengersId, backfillStealthId},
			expectViewUpdates: map[catalog.AlbumId][]catalog.Cover{backfillStealthId: stealthCovers},
		},
		{
			name: "it should return an empty report when the owner has no album",
			fields: fields{
				FindAlbumByOwnerPort:      &findAlbumByOwnerPortFake{},
				RandomiseCoversPort:       &randomiseCoversPortFake{},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
			},
			args:              args{owner: owner},
			wantReport:        catalog.BackfillReport{Albums: 0},
			wantErr:           assert.NoError,
			expectCalledIds:   nil,
			expectViewUpdates: nil,
		},
		{
			name: "it should continue after a per-album Randomise failure and skip the view update for that album",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				RandomiseCoversPort: &randomiseCoversPortFake{
					Errors:  map[catalog.AlbumId]error{backfillAvengersId: randomiseFailure},
					Changed: map[catalog.AlbumId][]catalog.Cover{backfillStealthId: stealthCovers},
				},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
			},
			args: args{owner: owner},
			wantReport: catalog.BackfillReport{
				Albums:   2,
				Failures: []catalog.BackfillFailure{{AlbumId: backfillAvengersId, Err: randomiseFailure}},
			},
			wantErr:           assert.NoError,
			expectCalledIds:   []catalog.AlbumId{backfillAvengersId, backfillStealthId},
			expectViewUpdates: map[catalog.AlbumId][]catalog.Cover{backfillStealthId: stealthCovers},
		},
		{
			name: "it should continue after a per-album view update failure and record it in the report",
			fields: fields{
				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
				RandomiseCoversPort: &randomiseCoversPortFake{
					Changed: map[catalog.AlbumId][]catalog.Cover{backfillAvengersId: avengersCovers, backfillStealthId: stealthCovers},
				},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{
					Errors: map[catalog.AlbumId]error{backfillAvengersId: viewFailure},
				},
			},
			args: args{owner: owner},
			wantReport: catalog.BackfillReport{
				Albums:   2,
				Failures: []catalog.BackfillFailure{{AlbumId: backfillAvengersId, Err: viewFailure}},
			},
			wantErr:           assert.NoError,
			expectCalledIds:   []catalog.AlbumId{backfillAvengersId, backfillStealthId},
			expectViewUpdates: map[catalog.AlbumId][]catalog.Cover{backfillStealthId: stealthCovers},
		},
		{
			name: "it should return an error when listing the owner's albums fails",
			fields: fields{
				FindAlbumByOwnerPort:      &findAlbumByOwnerPortFake{Err: listFailure},
				RandomiseCoversPort:       &randomiseCoversPortFake{},
				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
			},
			args:       args{owner: owner},
			wantReport: catalog.BackfillReport{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, listFailure)
			},
			expectCalledIds:   nil,
			expectViewUpdates: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backfill := &catalog.BackfillCovers{
				FindAlbumByOwnerPort:      tt.fields.FindAlbumByOwnerPort,
				RandomiseCoversPort:       tt.fields.RandomiseCoversPort,
				BackfillCoversViewUpdater: tt.fields.BackfillCoversViewUpdater,
			}

			report, err := backfill.BackfillForOwner(context.Background(), tt.args.owner)
			if !tt.wantErr(t, err, fmt.Sprintf("BackfillForOwner(%s)", tt.args.owner)) {
				return
			}
			assert.Equal(t, tt.wantReport, report, "backfill report")
			assert.Equal(t, tt.expectCalledIds, tt.fields.RandomiseCoversPort.Calls, "albums called")
			assert.Equal(t, tt.expectViewUpdates, tt.fields.BackfillCoversViewUpdater.Updates, "view updates")
		})
	}
}
