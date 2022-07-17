package archive

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
)

type Cache struct {
	storage CacheAdapter
}

func NewCache() *Cache {
	return &Cache{
		storage: cachePort,
	}
}

// GetOrStore will return cached value if present, or call contentGenerator to cache and return its result.
func (c *Cache) GetOrStore(cacheKey string, contentGenerator func() ([]byte, string, error), postProcess func(io.ReadCloser, int, string, error) ([]byte, string, error)) ([]byte, string, error) {
	reader, size, mediaType, err := c.storage.Get(cacheKey)
	if err != NotFoundError {
		return postProcess(reader, size, mediaType, err)
	}

	content, mediaType, err := contentGenerator()
	if err != nil {
		return postProcess(nil, 0, "", err)
	}

	err = c.storage.Put(cacheKey, mediaType, bytes.NewReader(content))
	return postProcess(io.NopCloser(bytes.NewReader(content)), len(content), mediaType, errors.Wrapf(err, "failed caching media %s", cacheKey))
}

// Store calls and cache the value only if it wasn't already present.
func (c *Cache) Store(cacheKey string, contentGenerator func() ([]byte, string, error)) error {
	_, _, _, err := c.storage.Get(cacheKey)
	if err != NotFoundError {
		return nil
	}

	content, mediaType, err := contentGenerator()
	if err != nil {
		return err
	}

	err = c.storage.Put(cacheKey, mediaType, bytes.NewReader(content))
	return errors.Wrapf(err, "failed caching media %s", cacheKey)
}
