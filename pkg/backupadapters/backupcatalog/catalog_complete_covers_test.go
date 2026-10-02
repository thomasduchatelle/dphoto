package backupcatalog_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/backupcatalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func TestCompleteCoversAdapter_CompleteCoversFromCandidates(t *testing.T) {
	const owner = ownermodel.Owner("ironman")
	const folderName = "/album-1"
	albumId := catalog.AlbumId{Owner: owner, FolderName: catalog.FolderName(folderName)}

	type args struct {
		owner           ownermodel.Owner
		albumFolderName string
		candidates      []backup.CoverCandidate
	}
	tests := []struct {
		name              string
		args              args
		expectedAlbumId   catalog.AlbumId
		expectedCandidate []*catalog.MediaMeta
	}{
		{
			name: "it should translate the backup candidates into catalog MediaMeta tagged as IMAGE",
			args: args{
				owner:           owner,
				albumFolderName: folderName,
				candidates: []backup.CoverCandidate{
					{MediaId: "media-1", Filename: "photo-1.jpg"},
					{MediaId: "media-2", Filename: "photo-2.jpg"},
				},
			},
			expectedAlbumId: albumId,
			expectedCandidate: []*catalog.MediaMeta{
				{Id: "media-1", Filename: "photo-1.jpg", Type: catalog.MediaTypeImage},
				{Id: "media-2", Filename: "photo-2.jpg", Type: catalog.MediaTypeImage},
			},
		},
		{
			name: "it should forward an empty candidate list unchanged",
			args: args{
				owner:           owner,
				albumFolderName: folderName,
				candidates:      nil,
			},
			expectedAlbumId:   albumId,
			expectedCandidate: []*catalog.MediaMeta{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := &catalogCompleteCoversSpy{}
			adapter := &backupcatalog.CompleteCoversAdapter{CatalogCompleteCovers: spy}

			err := adapter.CompleteCoversFromCandidates(context.Background(), tt.args.owner, tt.args.albumFolderName, tt.args.candidates)

			if assert.NoError(t, err) {
				assert.Equal(t, tt.expectedAlbumId, spy.gotAlbumId)
				assert.Equal(t, tt.expectedCandidate, spy.gotCandidates)
			}
		})
	}
}

type catalogCompleteCoversSpy struct {
	gotAlbumId    catalog.AlbumId
	gotCandidates []*catalog.MediaMeta
}

func (s *catalogCompleteCoversSpy) CompleteCoversFromCandidates(ctx context.Context, albumId catalog.AlbumId, candidates []*catalog.MediaMeta) error {
	s.gotAlbumId = albumId
	s.gotCandidates = candidates
	return nil
}
