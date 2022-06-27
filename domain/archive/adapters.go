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

	// AddLocation adds (or override) the media location with the new key
	AddLocation(owner, id, key string) error
}

type StoreAdapter interface {
	// Upload stores online the file and return the final key used
	Upload(values DestructuredKey, content io.Reader) (string, error)
}
