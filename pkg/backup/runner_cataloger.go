package backup

import "github.com/pkg/errors"

// newCreatorCataloger creates missing albums
func newCreatorCataloger(owner string) (runnerCataloger, error) {
	timeline, err := catalogPort.GetAlbumsTimeline(owner)

	return func(medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error) {
		idsMap, err := catalogPort.AssignIdsToNewMedias(owner, medias)
		if err != nil {
			return nil, err
		}

		return buildRequestsAndFireEvents(progressChannel, medias, idsMap, func(media *AnalysedMedia) (string, error) {
			folderName, created, err := timeline.FindOrCreateAlbum(media.Details.DateTime)
			if err != nil {
				return "", errors.Wrapf(err, "failed to find or create album for date %s", media.Details.DateTime)
			}
			if created {
				progressChannel <- &ProgressEvent{Type: ProgressEventAlbumCreated, Count: 1, Album: folderName}
			}

			return folderName, nil
		})

	}, err
}

// newAlbumFilterCataloger doesn't create any album, and filters media not listed in 'albums'
func newAlbumFilterCataloger(owner string, albums map[string]interface{}) (runnerCataloger, error) {
	if len(albums) == 0 {
		return nil, errors.Errorf("newAlbumFilterCataloger must be created with at least 1 album (otherwise no media will go through")
	}
	timeline, err := catalogPort.GetAlbumsTimeline(owner)

	return func(medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error) {
		var filteredMedias []*AnalysedMedia
		mediaAlbums := make(map[*AnalysedMedia]string)

		for _, media := range medias {
			folderName, exists, err := timeline.FindAlbum(media.Details.DateTime)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to find album for the date %s", media.Details.DateTime)
			}

			if _, included := albums[folderName]; exists && included {
				filteredMedias = append(filteredMedias, media)
				mediaAlbums[media] = folderName
			} else {
				progressChannel <- &ProgressEvent{Type: ProgressEventWrongAlbum, Count: 1, Size: media.FoundMedia.Size()}
			}
		}

		idsMap, err := catalogPort.AssignIdsToNewMedias(owner, filteredMedias)
		if err != nil {
			return nil, err
		}

		return buildRequestsAndFireEvents(progressChannel, filteredMedias, idsMap, func(media *AnalysedMedia) (string, error) {
			return mediaAlbums[media], nil
		})

	}, err
}

// buildRequestsAndFireEvents populate the list of BackingUpMediaRequest by keeping the same order as medias parameters (to simplify the testing)
func buildRequestsAndFireEvents(progressChannel chan *ProgressEvent, medias []*AnalysedMedia, idsMap map[*AnalysedMedia]string, getAlbum func(*AnalysedMedia) (string, error)) ([]*BackingUpMediaRequest, error) {
	var requests []*BackingUpMediaRequest

	catalogedCount := MediaCounterZero
	skippedCount := MediaCounterZero
	for _, media := range medias {
		if id, newMedia := idsMap[media]; newMedia {
			catalogedCount = catalogedCount.Add(1, media.FoundMedia.Size())

			album, err := getAlbum(media)
			if err != nil {
				return nil, err
			}

			requests = append(requests, &BackingUpMediaRequest{
				AnalysedMedia: media,
				Id:            id,
				FolderName:    album,
			})
		} else {
			skippedCount = skippedCount.Add(1, media.FoundMedia.Size())
		}
	}

	if catalogedCount.Count > 0 {
		progressChannel <- &ProgressEvent{Type: ProgressEventCatalogued, Count: catalogedCount.Count, Size: catalogedCount.Size}
	}
	if skippedCount.Count > 0 {
		progressChannel <- &ProgressEvent{Type: ProgressEventAlreadyExists, Count: skippedCount.Count, Size: skippedCount.Size}
	}

	return requests, nil
}

func newScannerCataloger(owner string) runnerCataloger {
	return func(medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error) {
		idsMap, err := catalogPort.AssignIdsToNewMedias(owner, medias)
		if err != nil {
			return nil, err
		}

		return buildRequestsAndFireEvents(progressChannel, medias, idsMap, func(media *AnalysedMedia) (string, error) {
			return "", nil
		})
	}
}
