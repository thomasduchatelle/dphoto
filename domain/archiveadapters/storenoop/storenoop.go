package storefs

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"io"
	"time"
)

type CacheAndStore interface {
	archive.StoreAdapter
	archive.CacheAdapter
}

func NewStore() CacheAndStore {
	return &store{}
}

type store struct {
}

func (s *store) Get(key string) (io.ReadCloser, int, string, error) {
	return nil, 0, "", archive.NotFoundError
}

func (s *store) Put(key string, mediaType string, content io.Reader) error {
	return nil
}

func (s *store) WalkCacheByPrefix(prefix string, observer func(string)) error {
	return nil
}

func (s *store) Download(key string) (io.ReadCloser, error) {
	return nil, archive.NotFoundError
}

func (s *store) Upload(keyHint archive.DestructuredKey, content io.Reader) (string, error) {
	return keyHint.Prefix + keyHint.Suffix, nil
}

func (s *store) Copy(origin string, keyHint archive.DestructuredKey) (string, error) {
	return keyHint.Prefix + keyHint.Suffix, nil
}

func (s *store) Delete(locations []string) error {
	return nil
}

func (s *store) SignedURL(key string, duration time.Duration) (string, error) {
	return "", errors.Errorf("SignedURL is not supported in offline / local mode")
}
