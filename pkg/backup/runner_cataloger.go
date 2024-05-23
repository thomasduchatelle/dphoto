package backup

import (
	"context"
)

// TODO Delete this file and all the function called that are not used anymore.

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

func newScannerCataloger(owner string) RunnerCatalogerFunc {
	return func(ctx context.Context, medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error) {
		idsMap, err := catalogPort.AssignIdsToNewMedias(owner, medias)
		if err != nil {
			return nil, err
		}

		return buildRequestsAndFireEvents(progressChannel, medias, idsMap, func(media *AnalysedMedia) (string, error) {
			return "", nil
		})
	}
}
