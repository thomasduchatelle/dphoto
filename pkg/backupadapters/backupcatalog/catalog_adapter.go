package backupcatalog

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"sync"
	"time"
)

func New() backup.CatalogAdapter {
	return &adapter{}
}

type adapter struct {
	dry bool
}

func (a *adapter) GetAlbumsTimeline(owner string) (backup.TimelineAdapter, error) {
	ctx := context.TODO()

	albums, err := catalog.FindAllAlbums(catalog.Owner(owner))
	if err != nil {
		return nil, err
	}

	timeline, err := catalog.NewTimeline(albums)
	if err != nil {
		return nil, err
	}

	return &timelineAdapter{
		owner:        owner,
		timeline:     timeline,
		timelineLock: sync.Mutex{},
		createAlbum:  pkgfactory.CreateAlbumCase(ctx),
	}, nil
}

func (a *adapter) AssignIdsToNewMedias(owner string, medias []*backup.AnalysedMedia) (map[*backup.AnalysedMedia]string, error) {
	ctx := context.TODO()

	signatures := make([]*catalog.MediaSignature, len(medias), len(medias))
	for i, media := range medias {
		signatures[i] = &catalog.MediaSignature{
			SignatureSha256: media.Sha256Hash,
			SignatureSize:   media.FoundMedia.Size(),
		}
	}

	assignedIds, err := pkgfactory.InsertMediasCase(ctx).AssignIdsToNewMedias(ctx, catalog.Owner(owner), signatures)

	mediasWithId := make(map[*backup.AnalysedMedia]string)
	for _, media := range medias {
		sign := catalog.MediaSignature{
			SignatureSha256: media.Sha256Hash,
			SignatureSize:   media.FoundMedia.Size(),
		}
		if id, found := assignedIds[sign]; found {
			mediasWithId[media] = id.Value()
		}
	}

	return mediasWithId, errors.Wrapf(err, "failed to assign ids")
}

func (a *adapter) IndexMedias(owner string, requests []*backup.CatalogMediaRequest) error {
	ctx := context.TODO()

	catalogRequests := make([]catalog.CreateMediaRequest, len(requests), len(requests))
	for i, request := range requests {
		details := request.BackingUpMediaRequest.AnalysedMedia.Details

		catalogRequests[i] = catalog.CreateMediaRequest{
			Id: catalog.MediaId(request.BackingUpMediaRequest.Id),
			Signature: catalog.MediaSignature{
				SignatureSha256: request.BackingUpMediaRequest.AnalysedMedia.Sha256Hash,
				SignatureSize:   request.BackingUpMediaRequest.AnalysedMedia.FoundMedia.Size(),
			},
			FolderName: catalog.NewFolderName(request.BackingUpMediaRequest.FolderName),
			Filename:   request.ArchiveFilename,
			Type:       catalog.MediaType(request.BackingUpMediaRequest.AnalysedMedia.Type),
			Details: catalog.MediaDetails{
				Width:         details.Width,
				Height:        details.Height,
				DateTime:      details.DateTime,
				Orientation:   catalog.MediaOrientation(details.Orientation),
				Make:          details.Make,
				Model:         details.Model,
				GPSLatitude:   details.GPSLatitude,
				GPSLongitude:  details.GPSLongitude,
				Duration:      details.Duration,
				VideoEncoding: details.VideoEncoding,
			},
		}
	}

	err := pkgfactory.InsertMediasCase(ctx).Insert(ctx, catalog.Owner(owner), catalogRequests)
	return errors.Wrapf(err, "failed to insert %d medias", len(catalogRequests))
}

type timelineAdapter struct {
	owner        string
	timeline     *catalog.Timeline
	timelineLock sync.Mutex
	createAlbum  *catalog.CreateAlbum
}

func (u *timelineAdapter) FindAlbum(mediaTime time.Time) (string, bool, error) {
	if album, exists := u.timeline.FindAt(mediaTime); exists {
		return album.FolderName.String(), exists, nil
	}

	return "", false, nil
}

func (u *timelineAdapter) FindOrCreateAlbum(mediaTime time.Time) (string, bool, error) {
	ctx := context.TODO()

	u.timelineLock.Lock()
	defer u.timelineLock.Unlock()

	if album, ok := u.timeline.FindAt(mediaTime); ok {
		return album.FolderName.String(), false, nil
	}

	year := mediaTime.Year()
	quarter := (mediaTime.Month() - 1) / 3

	createRequest := catalog.CreateAlbumRequest{
		Owner:            catalog.Owner(u.owner),
		Name:             fmt.Sprintf("Q%d %d", quarter+1, year),
		Start:            time.Date(year, quarter*3+1, 1, 0, 0, 0, 0, time.UTC),
		End:              time.Date(year, (quarter+1)*3+1, 1, 0, 0, 0, 0, time.UTC),
		ForcedFolderName: fmt.Sprintf("/%d-Q%d", year, quarter+1),
	}

	log.Infof("Creates new album '%s' to accommodate media at %s", createRequest.ForcedFolderName, mediaTime.Format(time.RFC3339))

	_, err := u.createAlbum.Create(ctx, createRequest)
	if err != nil {
		return "", false, errors.Wrapf(err, "failed to create album containing %s [%s]", mediaTime.Format(time.RFC3339), createRequest.String())
	}

	u.timeline, err = u.timeline.AppendAlbum(&catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      createRequest.Owner,
			FolderName: catalog.NewFolderName(createRequest.ForcedFolderName),
		},
		Name:  createRequest.Name,
		Start: createRequest.Start,
		End:   createRequest.End,
	})
	return createRequest.ForcedFolderName, true, errors.Wrapf(err, "failed to update internal timeline")
}
