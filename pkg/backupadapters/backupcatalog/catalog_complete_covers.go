package backupcatalog

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type CatalogCompleteCovers interface {
	CompleteCoversFromCandidates(ctx context.Context, albumId catalog.AlbumId, candidates []*catalog.MediaMeta) error
}

type CompleteCoversAdapter struct {
	CatalogCompleteCovers CatalogCompleteCovers
}

func (a *CompleteCoversAdapter) CompleteCoversFromCandidates(ctx context.Context, owner ownermodel.Owner, albumFolderName string, candidates []backup.CoverCandidate) error {
	mediaMetas := make([]*catalog.MediaMeta, len(candidates))
	for i, candidate := range candidates {
		mediaMetas[i] = &catalog.MediaMeta{
			Id:       catalog.MediaId(candidate.MediaId),
			Filename: candidate.Filename,
			Type:     catalog.MediaTypeImage,
		}
	}

	return a.CatalogCompleteCovers.CompleteCoversFromCandidates(ctx, catalog.AlbumId{
		Owner:      owner,
		FolderName: catalog.FolderName(albumFolderName),
	}, mediaMetas)
}
