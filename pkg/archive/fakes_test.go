package archive_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/thomasduchatelle/dphoto/pkg/archive"
)

// ARepositoryInMemory is an in-memory implementation of archive.ARepositoryAdapter.
type ARepositoryInMemory struct {
	Media map[string]map[string]string
}

func NewARepositoryInMemory() *ARepositoryInMemory {
	return &ARepositoryInMemory{Media: make(map[string]map[string]string)}
}

func (r *ARepositoryInMemory) FindById(owner, id string) (string, error) {
	if byId, ok := r.Media[owner]; ok {
		if key, ok := byId[id]; ok {
			return key, nil
		}
	}
	return "", archive.NotFoundError
}

func (r *ARepositoryInMemory) FindByIds(owner string, ids []string) (map[string]string, error) {
	result := make(map[string]string)
	byId := r.Media[owner]
	for _, id := range ids {
		if key, ok := byId[id]; ok {
			result[id] = key
		}
	}
	return result, nil
}

func (r *ARepositoryInMemory) AddLocation(owner, id, key string) error {
	if _, ok := r.Media[owner]; !ok {
		r.Media[owner] = make(map[string]string)
	}
	r.Media[owner][id] = key
	return nil
}

func (r *ARepositoryInMemory) UpdateLocations(owner string, locations map[string]string) error {
	if _, ok := r.Media[owner]; !ok {
		r.Media[owner] = make(map[string]string)
	}
	for id, key := range locations {
		r.Media[owner][id] = key
	}
	return nil
}

func (r *ARepositoryInMemory) FindIdsFromKeyPrefix(keyPrefix string) (map[string]string, error) {
	result := make(map[string]string)
	for _, byId := range r.Media {
		for id, key := range byId {
			if strings.HasPrefix(key, keyPrefix) {
				result[id] = key
			}
		}
	}
	return result, nil
}

// StoreInMemory is an in-memory implementation of archive.StoreAdapter.
type StoreInMemory struct {
	Content map[string][]byte
}

func NewStoreInMemory() *StoreInMemory {
	return &StoreInMemory{Content: make(map[string][]byte)}
}

func (s *StoreInMemory) Download(key string) (io.ReadCloser, error) {
	content, ok := s.Content[key]
	if !ok {
		return nil, archive.NotFoundError
	}
	return io.NopCloser(bytes.NewReader(content)), nil
}

func (s *StoreInMemory) Upload(values archive.DestructuredKey, content io.Reader) (string, error) {
	key := values.Prefix + values.Suffix
	data, err := io.ReadAll(content)
	if err != nil {
		return "", err
	}
	s.Content[key] = data
	return key, nil
}

func (s *StoreInMemory) Copy(origin string, destination archive.DestructuredKey) (string, error) {
	data, ok := s.Content[origin]
	if !ok {
		return "", archive.NotFoundError
	}
	key := destination.Prefix + destination.Suffix
	s.Content[key] = data
	return key, nil
}

func (s *StoreInMemory) Delete(locations []string) error {
	for _, key := range locations {
		delete(s.Content, key)
	}
	return nil
}

func (s *StoreInMemory) SignedURL(key string, duration time.Duration) (string, error) {
	if _, ok := s.Content[key]; !ok {
		return "", archive.NotFoundError
	}
	return "signed://" + key, nil
}

func (s *StoreInMemory) Has(key string) bool {
	_, ok := s.Content[key]
	return ok
}

// CacheInMemory is an in-memory implementation of archive.CacheAdapter.
type CacheInMemory struct {
	Content   map[string][]byte
	MediaType map[string]string
}

func NewCacheInMemory() *CacheInMemory {
	return &CacheInMemory{
		Content:   make(map[string][]byte),
		MediaType: make(map[string]string),
	}
}

func (c *CacheInMemory) Get(key string) (io.ReadCloser, int, string, error) {
	content, ok := c.Content[key]
	if !ok {
		return nil, 0, "", archive.NotFoundError
	}
	return io.NopCloser(bytes.NewReader(content)), len(content), c.MediaType[key], nil
}

func (c *CacheInMemory) Put(key string, mediaType string, content io.Reader) error {
	data, err := io.ReadAll(content)
	if err != nil {
		return err
	}
	c.Content[key] = data
	c.MediaType[key] = mediaType
	return nil
}

func (c *CacheInMemory) SignedURL(key string, duration time.Duration) (string, error) {
	return "signed-cache://" + key, nil
}

func (c *CacheInMemory) WalkCacheByPrefix(prefix string, observer func(string)) error {
	for key := range c.Content {
		if strings.HasPrefix(key, prefix) {
			observer(key)
		}
	}
	return nil
}

// AsyncJobInMemory records calls made to archive.AsyncJobAdapter.
type AsyncJobInMemory struct {
	LoadedImages [][]*archive.ImageToResize
	WarmUpCalls  []WarmUpCall
}

type WarmUpCall struct {
	Owner     string
	MissedKey string
	Width     int
}

func NewAsyncJobInMemory() *AsyncJobInMemory {
	return &AsyncJobInMemory{}
}

func (a *AsyncJobInMemory) WarmUpCacheByFolder(owner, missedStoreKey string, width int) error {
	a.WarmUpCalls = append(a.WarmUpCalls, WarmUpCall{
		Owner:     owner,
		MissedKey: missedStoreKey,
		Width:     width,
	})
	return nil
}

func (a *AsyncJobInMemory) LoadImagesInCache(images ...*archive.ImageToResize) error {
	a.LoadedImages = append(a.LoadedImages, images)
	return nil
}

// ResizerInMemory returns deterministic bytes for resize requests, encoding the
// requested width in the payload so tests can distinguish outputs.
type ResizerInMemory struct {
	MediaType string
}

func NewResizerInMemory() *ResizerInMemory {
	return &ResizerInMemory{
		MediaType: "image/jpeg",
	}
}

func (r *ResizerInMemory) ResizeImage(reader io.Reader, width int, fast bool) ([]byte, string, error) {
	_, _ = io.Copy(io.Discard, reader)
	return []byte(fmt.Sprintf("resized-w=%d", width)), r.MediaType, nil
}

func (r *ResizerInMemory) ResizeImageAtDifferentWidths(reader io.Reader, widths []int) (map[int][]byte, string, error) {
	_, _ = io.Copy(io.Discard, reader)
	result := make(map[int][]byte)
	for _, w := range widths {
		result[w] = []byte(fmt.Sprintf("resized-w=%d", w))
	}
	return result, r.MediaType, nil
}
