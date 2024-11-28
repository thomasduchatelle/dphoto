package backup

import (
	"context"
	"slices"
	"sync"
)

func newScanReportBuilder() *scanReportBuilder {
	return &scanReportBuilder{
		albums: make(map[string]*ScannedFolder),
	}
}

type scanReportBuilder struct {
	lock   sync.Mutex
	albums map[string]*ScannedFolder
}

func (s *scanReportBuilder) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, request := range requests {
		scannedFolder := s.getOrCreateScannedFolder(request.AnalysedMedia.FoundMedia)
		scannedFolder.PushBoundaries(request.AnalysedMedia.Details.DateTime, request.AnalysedMedia.FoundMedia.Size())
	}

	return nil
}

func (s *scanReportBuilder) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	scannedFolder := s.getOrCreateScannedFolder(found)
	scannedFolder.RejectsCount++
	return nil
}

func (s *scanReportBuilder) getOrCreateScannedFolder(foundMedia FoundMedia) *ScannedFolder {
	mediaPath := foundMedia.MediaPath()
	if _, ok := s.albums[mediaPath.Path]; !ok {
		s.albums[mediaPath.Path] = s.newFoundAlbum(mediaPath)
	}

	scannedFolder := s.albums[mediaPath.Path]
	return scannedFolder
}

func (s *scanReportBuilder) build() []*ScannedFolder {
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

func (s *scanReportBuilder) newFoundAlbum(mediaPath MediaPath) *ScannedFolder {
	return &ScannedFolder{
		Name:         mediaPath.ParentDir,
		RelativePath: mediaPath.Path,
		FolderName:   mediaPath.ParentDir,
		AbsolutePath: mediaPath.ParentFullPath,
		Distribution: make(map[string]MediaCounter),
	}
}
