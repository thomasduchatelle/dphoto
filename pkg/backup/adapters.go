package backup

import (
	"io"
	"time"
)

var (
	catalogPort    CatalogAdapter
	archivePort    BArchiveAdapter
	detailsReaders []DetailsReaderAdapter // DetailsReaders is a list of specific details extractor can auto-register
)

func RegisterDetailsReader(reader DetailsReaderAdapter) {
	detailsReaders = append(detailsReaders, reader)
}

func Init(catalog CatalogAdapter, archive BArchiveAdapter) {
	catalogPort = catalog
	archivePort = archive
}

type TimelineAdapter interface {
	FindOrCreateAlbum(mediaTime time.Time) (folderName string, created bool, err error)
	FindAlbum(dateTime time.Time) (folderName string, exists bool, err error)
}

type CatalogAdapter interface {
	// GetAlbumsTimeline create and initialise a timeline optimised to find album from a date, and create missing albums.
	GetAlbumsTimeline(owner string) (TimelineAdapter, error)

	// AssignIdsToNewMedias filter out existing medias and generate an ID for new ones.
	AssignIdsToNewMedias(owner string, medias []*AnalysedMedia) (map[*AnalysedMedia]string, error)

	// IndexMedias add to the catalog following medias
	IndexMedias(owner string, requests []*CatalogMediaRequest) error
}

type BArchiveAdapter interface {
	// ArchiveMedia uploads the file in the right folder but might change the name to avoid clash with other existing files. Use files name is always returned.
	ArchiveMedia(owner string, media *BackingUpMediaRequest) (string, error)
}

type DetailsReaderAdapter interface {
	// Supports returns true if the file can be parsed with this reader. False otherwise.
	Supports(media FoundMedia, mediaType MediaType) bool
	// ReadDetails extracts metadata from the content of the file.
	ReadDetails(reader io.Reader, options DetailsReaderOptions) (*MediaDetails, error)
}

type DetailsReaderOptions struct {
	Fast bool // Fast true indicate the parser should focus at extracting the date, nothing else
}
