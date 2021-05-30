package importer

import (
	"duchatelle.io/dphoto/dphoto/scanner"
	"path"
	"time"
)

type DirectoryBoundary struct {
	Start time.Time
	End   time.Time
}

type StructureDiscovery struct {
	Directories map[string]DirectoryBoundary
}

// StructureDiscovery is to be used in place of an uploader to collect min/max date of each folder
func (s *StructureDiscovery) StructureDiscovery(buffer []*scanner.AnalysedMedia, progressChannel chan *scanner.ProgressEvent) error {
	for _, media := range buffer {
		dir := path.Dir(buffer[0].FoundMedia.Filename())
		if boundaries, ok := s.Directories[dir]; ok {
			boundaries.Push(media.Details.DateTime)
		} else {
			s.Directories[dir] = DirectoryBoundary{
				Start: media.Details.DateTime,
				End:   media.Details.DateTime,
			}
		}
	}

	return nil
}

// Push move the boundaries if the date is outside the range
func (b *DirectoryBoundary) Push(dateTime time.Time) {
	if b.Start.After(dateTime) {
		b.Start = dateTime
	}

	if b.End.Before(dateTime) {
		b.End = dateTime
	}
}
