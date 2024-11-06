package backup

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"io"
	"time"
)

var (
	archivePort       BArchiveAdapter
	referencerFactory CataloguerFactory
	insertMediaPort   InsertMediaPort
)

// Init for scan or backup (but only refFactory is required for scan)
func Init(archive BArchiveAdapter, refFactory CataloguerFactory, insertMedia InsertMediaPort) {
	archivePort = archive
	referencerFactory = refFactory
	insertMediaPort = insertMedia
}

type InsertMediaPort interface {
	IndexMedias(ctx context.Context, owner ownermodel.Owner, requests []*CatalogMediaRequest) error
}

type TimelineAdapter interface {
	FindOrCreateAlbum(mediaTime time.Time) (folderName string, created bool, err error)
	FindAlbum(dateTime time.Time) (folderName string, exists bool, err error)
}

type AlbumLookupPort interface {
	FindOrCreateAlbum(owner ownermodel.Owner, mediaTime time.Time) (folderName string, created bool, err error)
}

type CatalogAdapter interface {
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
	Fast bool // Fast true indicate the parser should focus at extracting the date, nothing else TODO can be retired
}

type CataloguerFactory interface {
	// NewAlbumCreatorCataloguer returns a Referencer that will create the album if the date is not yet covered.
	NewAlbumCreatorCataloguer(ctx context.Context, owner ownermodel.Owner) (Cataloguer, error)
	// NewDryRunCataloguer returns a Referencer that will not create any album.
	NewDryRunCataloguer(ctx context.Context, owner ownermodel.Owner) (Cataloguer, error)
}
