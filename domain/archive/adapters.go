package archive

import (
	"io"
)

var (
	repositoryPort ARepositoryAdapter
	storePort      StoreAdapter
)

func Init(repository ARepositoryAdapter, store StoreAdapter) {
	repositoryPort = repository
	storePort = store
}

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

type StoreAdapter interface {
	// Upload stores online the file and return the final key used
	Upload(values DestructuredKey, content io.Reader) (string, error)

	// Copy copied the file to a different location, without overriding existing file
	Copy(origin string, destination DestructuredKey) (string, error)

	// Delete permanently stored files (certainly after having been moved.
	Delete(locations []string) error
}
