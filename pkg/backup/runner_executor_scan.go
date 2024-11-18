package backup

import (
	"context"
	"sync"
)

func newScanReport() *scanReport {
	return &scanReport{
		albums: make(map[string]*ScannedFolder),
	}
}

type scanReport struct {
	lock   sync.Mutex
	albums map[string]*ScannedFolder
}

func (s *scanReport) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, request := range requests {
		scannedFolder := s.getOrCreateScannedFolder(request.AnalysedMedia.FoundMedia)
		scannedFolder.PushBoundaries(request.AnalysedMedia.Details.DateTime, request.AnalysedMedia.FoundMedia.Size())
	}

	return nil
}

func (s *scanReport) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	scannedFolder := s.getOrCreateScannedFolder(found)
	scannedFolder.RejectsCount++
	return nil
}

func (s *scanReport) getOrCreateScannedFolder(foundMedia FoundMedia) *ScannedFolder {
	mediaPath := foundMedia.MediaPath()
	if _, ok := s.albums[mediaPath.Path]; !ok {
		s.albums[mediaPath.Path] = s.newFoundAlbum(mediaPath)
	}

	scannedFolder := s.albums[mediaPath.Path]
	return scannedFolder
}

func (s *scanReport) collect() []*ScannedFolder {
	suggestions := make([]*ScannedFolder, 0, len(s.albums))
	for _, album := range s.albums {
		suggestions = append(suggestions, album)
	}

	slices.SortFunc(suggestions, func(i, j *ScannedFolder) int {
		if i.Start.Equal(j.Start) {
			return i.Start.Compare(j.Start)
		}

		return i.End.Compare(j.End)
	})

	return suggestions
}

func (s *scanReport) newFoundAlbum(mediaPath MediaPath) *ScannedFolder {
	return &ScannedFolder{
		Name:         mediaPath.ParentDir,
		RelativePath: mediaPath.Path,
		FolderName:   mediaPath.ParentDir,
		AbsolutePath: mediaPath.ParentFullPath,
		Distribution: make(map[string]MediaCounter),
	}
}
