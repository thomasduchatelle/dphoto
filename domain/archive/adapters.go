package archive

import (
	"github.com/thomasduchatelle/dphoto/domain/archive/image_resize"
	"io"
	"time"
)

var (
	repositoryPort ARepositoryAdapter
	storePort      StoreAdapter
	cachePort      CacheAdapter
	ResizerPort    ResizerAdapter = image_resize.NewResizer() // ResizerPort can be overrided for testing purpose
)

func Init(repository ARepositoryAdapter, store StoreAdapter, cache CacheAdapter) {
	repositoryPort = repository
	storePort = store
	cachePort = cache
}

// ARepositoryAdapter is storing the mapping between keys in the main storage and the media ids.
type ARepositoryAdapter interface {
	// FindById returns found location (key) of the media, or a NotFoundError
	FindById(owner, id string) (string, error)

	// FindByIds searches multiple physical location at once
	FindByIds(owner string, ids []string) (map[string]string, error)

	// AddLocation adds (or override) the media location with the new key
	AddLocation(owner, id, key string) error

	// UpdateLocations will update or set location for each id
	UpdateLocations(owner string, locations map[string]string) error
}

// StoreAdapter is the adapter where the original medias are stored (cool storage - safe for long term)
type StoreAdapter interface {
	// Download retrieves the file store at this key, raise a NotFoundError if the key doesn't exist
	Download(key string) (io.ReadCloser, error)

	// Upload stores online the file and return the final key used
	Upload(values DestructuredKey, content io.Reader) (string, error)

	// Copy copied the file to a different location, without overriding existing file
	Copy(origin string, destination DestructuredKey) (string, error)

	// Delete permanently stored files (certainly after having been moved.
	Delete(locations []string) error

	// SignedURL returns a pre-authorised URL to download the content
	SignedURL(key string, duration time.Duration) (string, error)
}

// CacheAdapter is the adapter where the re-sized medias are stored (hot storage - not long term safe)
type CacheAdapter interface {
	// Get retrieve the file store at this key, raise a NotFoundError if the key doesn't exist
	Get(key string) (io.ReadCloser, int, string, error)

	// Put stores the content by overriding exiting file if any
	Put(key string, mediaType string, content io.Reader) error

	// SignedURL returns a pre-authorised URL to download the content
	SignedURL(key string, duration time.Duration) (string, error)
}

type ResizerAdapter interface {
	ResizeImage(reader io.Reader, width int, fast bool) ([]byte, string, error)
}