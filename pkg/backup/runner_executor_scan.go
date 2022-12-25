package backup

import (
	log "github.com/sirupsen/logrus"
	"sort"
)

func newScanReceiver(mdc *log.Entry, volume SourceVolume) *scanReceiver {
	return &scanReceiver{
		mdc:    mdc,
		volume: volume,
		albums: make(map[string]*ScannedFolder),
	}
}

type scanReceiver struct {
	mdc     *log.Entry
	volume  SourceVolume
	albums  map[string]*ScannedFolder
	rejects []FoundMedia
}

func (s *scanReceiver) receive(buffer []*BackingUpMediaRequest, progressChannel chan *ProgressEvent) error {
	for _, media := range buffer {
		foundMedia := media.AnalysedMedia.FoundMedia
		if media.AnalysedMedia.Details.DateTime.IsZero() {
			log.Warnf("Media timestamp cannot be found within the file %s", foundMedia)
			s.rejects = append(s.rejects, foundMedia)

		} else {
			mediaPath := foundMedia.MediaPath()
			if _, ok := s.albums[mediaPath.Path]; !ok {
				s.albums[mediaPath.Path] = s.newFoundAlbum(mediaPath)
			}

			s.albums[mediaPath.Path].PushBoundaries(media.AnalysedMedia.Details.DateTime, foundMedia.Size())
		}

		progressChannel <- &ProgressEvent{
			Type:      ProgressEventUploaded,
			Count:     1,
			Size:      foundMedia.Size(),
			Album:     media.FolderName,
			MediaType: media.AnalysedMedia.Type,
		}
	}

	return nil
}

func (s *scanReceiver) collect() []*ScannedFolder {
	suggestions := make([]*ScannedFolder, len(s.albums))
	i := 0
	for _, album := range s.albums {
		suggestions[i] = album
		i++
	}
	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Start != suggestions[j].Start {
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
