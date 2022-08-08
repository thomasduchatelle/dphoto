// Package asyncjobnoop discards all jobs meant to be ran asynchronously
package asyncjobnoop

import (
	"github.com/thomasduchatelle/dphoto/domain/archive"
)

func New() archive.AsyncJobAdapter {
	return new(adapter)
}

type adapter struct {
}

func (a adapter) WarmUpCacheByFolder(owner, missedStoreKey string, width int) error {
	return nil
}

func (a adapter) LoadImagesInCache(images ...*archive.ImageToResize) error {
	return nil
}
