package backup

import "github.com/pkg/errors"

// newCreatorCataloger creates missing albums
func newCreatorCataloger(owner string) (runnerCataloger, error) {
	timeline, err := catalogPort.GetAlbumsTimeline(owner)

	return func(medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error) {
		var requests []*BackingUpMediaRequest
		idsMap, err := catalogPort.AssignIdsToNewMedias(owner, medias)
		if err != nil {
			return nil, err
		}

		for id, media := range idsMap {
			folderName, created, err := timeline.FindOrCreateAlbum(media.Details.DateTime)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to find or create album for date %s", media.Details.DateTime)
			}
			if created {
				progressChannel <- &ProgressEvent{Type: ProgressEventAlbumCreated, Count: 1, Album: folderName}
			}

			requests = append(requests, &BackingUpMediaRequest{
				AnalysedMedia: media,
				Id:            id,
				FolderName:    folderName,
			})
		}

		progressChannel <- &ProgressEvent{Type: ProgressEventCatalogued, Count: len(idsMap)}
		skippedEvent := &ProgressEvent{Type: ProgressEventAlreadyExists, Count: len(medias) - len(idsMap)}
		if skippedEvent.Count > 0 {
			progressChannel <- skippedEvent
		}

		return requests, nil

	}, err
}

// newReadOnlyCataloger doesn't create any album, and filters media not listed in 'albums'
func newReadOnlyCataloger(owner string, albums map[string]interface{}) (runnerCataloger, error) {
	timeline, err := catalogPort.GetAlbumsTimeline(owner)

	return func(medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error) {
		var filteredMedias []*AnalysedMedia
		mediaAlbums := make(map[*AnalysedMedia]string)

		for _, media := range medias {
			folderName, exists, err := timeline.FindAlbum(media.Details.DateTime)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to find album for the date %s", media.Details.DateTime)
			}
			_, included := albums[folderName]

			if exists && (len(albums) == 0 || included) {
				filteredMedias = append(filteredMedias, media)
				mediaAlbums[media] = folderName
			} else {
				progressChannel <- &ProgressEvent{Type: ProgressEventWrongAlbum, Count: 1, Size: media.FoundMedia.Size()}
			}
		}

		var requests []*BackingUpMediaRequest

		idsMap, err := catalogPort.AssignIdsToNewMedias(owner, filteredMedias)
		if err != nil {
			return nil, err
		}

		for id, media := range idsMap {
			requests = append(requests, &BackingUpMediaRequest{
				AnalysedMedia: media,
				Id:            id,
				FolderName:    mediaAlbums[media],
			})
		}

		progressChannel <- &ProgressEvent{Type: ProgressEventCatalogued, Count: len(idsMap)}
		skippedEvent := &ProgressEvent{Type: ProgressEventAlreadyExists, Count: len(filteredMedias) - len(idsMap)}
		if skippedEvent.Count > 0 {
			progressChannel <- skippedEvent
		}

		return requests, nil

	}, err
}
