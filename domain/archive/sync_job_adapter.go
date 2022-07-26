package archive

import (
	"context"
)

func NewSyncJobAdapter() AsyncJobAdapter {
	return new(syncJobAdapter)
}

type syncJobAdapter struct {
}

func (s syncJobAdapter) WarmUpCacheByFolder(owner, missedStoreKey string, width int) error {
	return WarmUpCacheByFolder(owner, missedStoreKey, width)
}

func (s syncJobAdapter) LoadImagesInCache(images ...*ImageToResize) error {
	_, err := LoadImagesInCache(context.Background(), images...)
	return err
}
