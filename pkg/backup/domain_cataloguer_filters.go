package backup

import (
	"context"
	"github.com/pkg/errors"
	"slices"
	"sync"
)

var (
	ErrCatalogerFilterMustBeInAlbum        = errors.New("media must be in album")
	ErrCatalogerFilterMustNotAlreadyExists = errors.New("media must not already exists")
	ErrMediaMustNotBeDuplicated            = errors.New("media is present twice in the volume")
)

type cataloguerFilter interface {
	// FilterOut returns an error if the media must be filtered out
	FilterOut(ctx context.Context, media AnalysedMedia, reference CatalogReference) error
}

func mustBeInAlbum(albumFolderNames ...string) cataloguerFilter {
	return &mustBeInAlbumCatalogerFilter{albumFolderNames: albumFolderNames}
}

func mustNotExists() cataloguerFilter {
	return new(mustNotAlreadyExistsCatalogerFilter)
}

func mustBeUniqueInVolume() cataloguerFilter {
	return &uniqueFilter{
		lock:          sync.Mutex{},
		uniqueIndexes: make(map[string]interface{}),
	}
}

type mustBeInAlbumCatalogerFilter struct {
	albumFolderNames []string
}

func (m mustBeInAlbumCatalogerFilter) FilterOut(ctx context.Context, media AnalysedMedia, reference CatalogReference) error {
	if slices.Contains(m.albumFolderNames, reference.AlbumFolderName()) {
		return nil
	}

	return ErrCatalogerFilterMustBeInAlbum
}

type mustNotAlreadyExistsCatalogerFilter struct{}

func (m mustNotAlreadyExistsCatalogerFilter) FilterOut(ctx context.Context, media AnalysedMedia, reference CatalogReference) error {
	if reference.Exists() {
		return ErrCatalogerFilterMustNotAlreadyExists
	}
	return nil
}

type uniqueFilter struct {
	lock          sync.Mutex
	uniqueIndexes map[string]interface{}
}

func (u *uniqueFilter) FilterOut(ctx context.Context, media AnalysedMedia, reference CatalogReference) error {
	u.lock.Lock()
	defer u.lock.Unlock()

	uniqueId := reference.UniqueIdentifier()
	if _, filterOut := u.uniqueIndexes[uniqueId]; filterOut {
		return ErrMediaMustNotBeDuplicated
	}

	u.uniqueIndexes[uniqueId] = nil
	return nil
}
