package archive

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"path"
	"strings"
)

// WarmUpCacheByFolder list medias missing in the cache and load them
func WarmUpCacheByFolder(owner, missedStoreKey string, width int) error {
	if !strings.HasPrefix(missedStoreKey, owner) {
		return errors.Errorf("cannot load the cache: %s is not the owner of %s media", owner, missedStoreKey)
	}

	parent := path.Dir(missedStoreKey)
	ids, err := repositoryPort.FindIdsFromKeyPrefix(parent)
	if err != nil {
		return errors.Wrapf(err, "finding ids with prefix %s", parent)
	}

	err = cachePort.WalkCacheByPrefix(generateCacheId(owner, "", width), func(cacheKey string) {
		mediaId := strings.TrimSuffix(path.Base(cacheKey), path.Ext(cacheKey))
		delete(ids, mediaId)
	})
	if err != nil {
		return errors.Wrapf(err, "walking cache with prefix %s", generateCacheId(owner, "", width))
	}

	var images []*ImageToResize
	for mediaId, storeKey := range ids {
		images = append(images, &ImageToResize{
			Owner:    owner,
			MediaId:  mediaId,
			StoreKey: storeKey,
			Widths:   []int{width},
			Open:     nil,
		})
	}

	if len(images) > 0 {
		err = asyncJobPort.LoadImagesInCache(images...)
		return errors.Wrapf(err, "loading images in cache (or sending the message)")
	} else {
		log.Infof("No more media to cache in %s for width=%d", parent, width)
	}

	return nil
}

// LoadImagesInCache generates resized images and store them in the cache. Returns the remaining ids if context is cancelled.
func LoadImagesInCache(ctx context.Context, images ...*ImageToResize) ([]*ImageToResize, error) {
	for index, img := range images {
		select {
		case <-ctx.Done():
			log.WithField("Owner", img.Owner).Warnf("time over - generated %d / %d images", index, len(images))
			return images[index:], nil

		default:
			opener := img.Open
			if opener == nil {
				opener = contentOpener(img.Owner, img.MediaId, img.StoreKey)
			}

			reader, err := opener()
			if err != nil {
				return nil, errors.Wrapf(err, "opening %s / %s", img.Owner, img.MediaId)
			}

			err = generateMiniature(img.Owner, img.MediaId, reader, img.Widths)
			_ = reader.Close()
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func contentOpener(owner, mediaId, storeKey string) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		var err error

		if storeKey == "" {
			storeKey, err = repositoryPort.FindById(owner, mediaId)
			if err != nil {
				return nil, errors.Wrapf(err, "finding storeKey for mediaId '%s'", mediaId)
			}
		}

		return storePort.Download(storeKey)
	}
}

func generateMiniature(owner, mediaId string, reader io.Reader, widths []int) error {
	contents, mediaType, err := ResizerPort.ResizeImageAtDifferentWidths(reader, widths)
	if err != nil {
		return errors.Wrapf(err, "failed to generate minature")
	}

	for width, content := range contents {
		cacheId := generateCacheId(owner, mediaId, width)
		err = cachePort.Put(cacheId, mediaType, bytes.NewReader(content))
		if err != nil {
			return errors.Wrapf(err, "inserting in cache %s at width=%d", mediaId, width)
		}
	}

	return nil
}
