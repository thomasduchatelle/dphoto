package backup

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"io"
	"time"
)

var (
	catalogPort       CatalogAdapter
	archivePort       BArchiveAdapter
	detailsReaders    []DetailsReaderAdapter // DetailsReaders is a list of specific details extractor can auto-register
	referencerFactory ReferencerFactory
)

func ClearDetailsReader() {
	detailsReaders = nil
}

func RegisterDetailsReader(reader DetailsReaderAdapter) {
	detailsReaders = append(detailsReaders, reader)
}

func Init(catalog CatalogAdapter, archive BArchiveAdapter, refFactory ReferencerFactory) {
	catalogPort = catalog // TODO is it still required ?
	archivePort = archive
	referencerFactory = refFactory
}

type TimelineAdapter interface {
	FindOrCreateAlbum(mediaTime time.Time) (folderName string, created bool, err error)
	FindAlbum(dateTime time.Time) (folderName string, exists bool, err error)
}

type AlbumLookupPort interface {
	FindOrCreateAlbum(owner ownermodel.Owner, mediaTime time.Time) (folderName string, created bool, err error)
}

type CatalogAdapter interface {
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

type ReferencerFactory interface {
	// NewCreatorReferencer returns a Referencer that will create the album if the date is not yet covered.
	NewCreatorReferencer(ctx context.Context, owner ownermodel.Owner) (CatalogReferencer, error)
	// NewDryRunReferencer returns a Referencer that will not create any album.
	NewDryRunReferencer(ctx context.Context, owner ownermodel.Owner) (CatalogReferencer, error)
}
