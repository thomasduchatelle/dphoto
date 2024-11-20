package backup

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"io"
	"time"
)

type InsertMediaPort interface {
	IndexMedias(ctx context.Context, owner ownermodel.Owner, requests []*CatalogMediaRequest) error
}

type TimelinePort interface {
	FindOrCreateAlbum(mediaTime time.Time) (folderName string, created bool, err error)
	FindAlbum(dateTime time.Time) (folderName string, exists bool, err error)
}

type AlbumLookupPort interface {
	FindOrCreateAlbum(owner ownermodel.Owner, mediaTime time.Time) (folderName string, created bool, err error)
}

type IndexMediaPort interface {
	// IndexMedias add to the catalog following medias
	IndexMedias(owner string, requests []*CatalogMediaRequest) error
}

type ArchiveMediaPort interface {
	// ArchiveMedia uploads the file in the right folder but might change the name to avoid clash with other existing files. Use files name is always returned.
	ArchiveMedia(owner string, media *BackingUpMediaRequest) (string, error)
}

type DetailsReader interface {
	// Supports returns true if the file can be parsed with this reader. False otherwise.
	Supports(media FoundMedia, mediaType MediaType) bool
	// ReadDetails extracts metadata from the content of the file.
	ReadDetails(reader io.Reader, options DetailsReaderOptions) (*MediaDetails, error)
}

type DetailsReaderOptions struct {
	Fast bool // Fast true indicate the parser should focus at extracting the date, nothing else TODO can be retired
}

// CataloguerFactory returns a Cataloguer scoped to a single owner. Implementations can be read-only, or have the behaviour of creating missing albums.
type CataloguerFactory interface {
	NewOwnerScopedCataloguer(ctx context.Context, owner ownermodel.Owner) (Cataloguer, error)
}
