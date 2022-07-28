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

// SupportResize will return true if the file format can be resized (and cached)
func SupportResize(filename string) bool {
	_, supported := supportedExtensionsForResizing[strings.ToLower(path.Ext(filename))]
	return supported
}

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

	totalCount := len(ids)

	err = cachePort.WalkCacheByPrefix(generateCacheId(owner, "", width), func(cacheKey string) {
		mediaId := strings.TrimSuffix(path.Base(cacheKey), path.Ext(cacheKey))
		delete(ids, mediaId)
	})
	if err != nil {
		return errors.Wrapf(err, "walking cache with prefix %s", generateCacheId(owner, "", width))
	}

	var images []*ImageToResize
	for mediaId, storeKey := range ids {
		if SupportResize(storeKey) {
			images = append(images, &ImageToResize{
				Owner:    owner,
				MediaId:  mediaId,
				StoreKey: storeKey,
				Widths:   []int{width},
				Open:     nil,
			})
		}
	}

	if len(images) > 0 {
		log.WithField("Owner", owner).Infof("%d / %d medias have been flagged to be cached at width=%d, %d are not supported", len(images), totalCount, width, len(ids)-len(images))
		err = asyncJobPort.LoadImagesInCache(images...)
		return errors.Wrapf(err, "loading images in cache (or sending the message)")
	} else {
		log.Infof("No more media to cache in %s out of %d for width=%d, %d are not supported", parent, totalCount, width, len(ids)-len(images))
	}

	return nil
}

// LoadImagesInCache generates resized images and store them in the cache. Returns how many has been processed
func LoadImagesInCache(ctx context.Context, images ...*ImageToResize) (int, error) {
	for index, img := range images {
		select {
		case <-ctx.Done():
			log.WithField("Owner", img.Owner).Warnf("time over - generated %d / %d images", index, len(images))
			return index, nil

		default:
			opener := img.Open
			storeKey := img.StoreKey // might be updated after opener() is called
			if opener == nil {
				opener = contentOpener(img.Owner, img.MediaId, &storeKey)
			}

			reader, err := opener()
			if err != nil {
				log.WithField("Owner", img.Owner).WithError(err).Errorf("opening %s/%s [%s] failed: %s", img.Owner, img.MediaId, storeKey, err.Error())
				continue
			}

			err = generateMiniature(img.Owner, img.MediaId, reader, img.Widths)
			_ = reader.Close()
			if err != nil {
				log.WithField("Owner", img.Owner).WithError(err).Errorf("failed to cache resized version of %s/%s from store '%s'", img.Owner, img.MediaId, storeKey)
				continue
			}
			log.WithField("Owner", img.Owner).Infof("[%d/%d] miniaturised %s into %v", index+1, len(images), storeKey, img.Widths)
		}
	}

	return len(images), nil
}

func contentOpener(owner, mediaId string, storeKey *string) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		var err error

		if *storeKey == "" {
			*storeKey, err = repositoryPort.FindById(owner, mediaId)
			if err != nil {
				return nil, errors.Wrapf(err, "finding storeKey for mediaId '%s'", mediaId)
			}
		}

		return storePort.Download(*storeKey)
	}
}

func generateMiniature(owner, mediaId string, reader io.Reader, widths []int) error {
	contents, mediaType, err := ResizerPort.ResizeImageAtDifferentWidths(reader, widths)
	if err != nil {
		return err
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
