package catalogviews

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// AlbumSummaryReprojector rebuilds AlbumSummaryForUsers projections from canonical albums,
// combining media counts, viewer availabilities and the album cover set queried per call.
// The albums are supplied by the caller so it can decide the scope (per-owner, per-album,
// ...) without a second round-trip.
type AlbumSummaryReprojector struct {
	ListUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
	MediaCounterPort              MediaCounterPort
	FindCoversByAlbumPort         FindCoversByAlbumPort
}

// Reproject returns the expected AlbumSummaryForUsers for every album.
func (r *AlbumSummaryReprojector) Reproject(ctx context.Context, albums []*catalog.Album) ([]AlbumSummaryForUsers, error) {
	if len(albums) == 0 {
		return nil, nil
	}

	albumIds := make([]catalog.AlbumId, len(albums))
	for i, album := range albums {
		albumIds[i] = album.AlbumId
	}

	availabilities, err := r.ListUserWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, albumIds...)
	if err != nil {
		return nil, err
	}

	counts, err := r.MediaCounterPort.CountMedia(ctx, albumIds...)
	if err != nil {
		return nil, err
	}

	covers := make(map[catalog.AlbumId][]catalog.Cover, len(albums))
	for _, albumId := range albumIds {
		albumCovers, err := r.FindCoversByAlbumPort.FindCoversByAlbum(ctx, albumId)
		if err != nil {
			return nil, err
		}
		covers[albumId] = albumCovers
	}

	summaries := make([]AlbumSummaryForUsers, 0, len(albums))
	for _, album := range albums {
		summaries = append(summaries, AlbumSummaryForUsers{
			AlbumSummary: AlbumSummary{
				AlbumId:    album.AlbumId,
				MediaCount: counts[album.AlbumId],
				Name:       album.Name,
				Start:      album.Start,
				End:        album.End,
				Covers:     covers[album.AlbumId],
			},
			Users: availabilities[album.AlbumId],
		})
	}

	return summaries, nil
}
