package backup

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sort"
	"sync"
)

func newScanReceiver(mdc *log.Entry, volume SourceVolume) *scanReceiver {
	return &scanReceiver{
		mdc:    mdc,
		volume: volume,
		albums: make(map[string]*ScannedFolder),
	}
}

type scanReceiver struct {
	mdc    *log.Entry
	lock   sync.Mutex
	volume SourceVolume
	albums map[string]*ScannedFolder
}

func (s *scanReceiver) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, request := range requests {
		scannedFolder := s.getOrCreateScannedFolder(request.AnalysedMedia.FoundMedia)
		scannedFolder.PushBoundaries(request.AnalysedMedia.Details.DateTime, request.AnalysedMedia.FoundMedia.Size())
	}

	return nil
}

func (s *scanReceiver) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	scannedFolder := s.getOrCreateScannedFolder(found)
	scannedFolder.RejectsCount++
	return nil
}

func (s *scanReceiver) getOrCreateScannedFolder(foundMedia FoundMedia) *ScannedFolder {
	mediaPath := foundMedia.MediaPath()
	if _, ok := s.albums[mediaPath.Path]; !ok {
		s.albums[mediaPath.Path] = s.newFoundAlbum(mediaPath)
	}

	scannedFolder := s.albums[mediaPath.Path]
	return scannedFolder
}

func (s *scanReceiver) collect() []*ScannedFolder {
	suggestions := make([]*ScannedFolder, 0, len(s.albums))
	for _, album := range s.albums {
		suggestions = append(suggestions, album)
	}

	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Start.Equal(suggestions[j].Start) {
			return suggestions[i].Start.Before(suggestions[j].Start)
		}

		return suggestions[i].End.Before(suggestions[j].End)
	})

	return suggestions
}

func (s *scanReceiver) newFoundAlbum(mediaPath MediaPath) *ScannedFolder {
	return &ScannedFolder{
		Name:         mediaPath.ParentDir,
		RelativePath: mediaPath.Path,
		FolderName:   mediaPath.ParentDir,
		AbsolutePath: mediaPath.ParentFullPath,
		Distribution: make(map[string]MediaCounter),
	}
}
