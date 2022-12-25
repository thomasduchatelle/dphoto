package archive

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"strings"
)

// GetResizedImage returns the image in the requested size (or rounded up), and the media type.
func GetResizedImage(owner, mediaId string, width int, maxBytes int) ([]byte, string, error) {
	cachedWidth, err := findCacheableSize(width)
	if err != nil {
		return nil, "", err
	}

	cacheKey := generateCacheId(owner, mediaId, cachedWidth)

	// note: retention and storage class is managed at infrastructure level
	return NewCache().GetOrStore(
		cacheKey,
		func() ([]byte, string, error) {
			key, err := repositoryPort.FindById(owner, mediaId)
			if err != nil {
				return nil, "", err
			}

			log.WithFields(log.Fields{
				"Owner": owner,
			}).Infof("%s [%s] is missing in the cache at size %d (requested %d)", key, cacheKey, cachedWidth, width)

			err = asyncJobPort.WarmUpCacheByFolder(owner, key, cachedWidth)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"Owner":    owner,
					"StoreKey": key,
				}).Warnf("Queuing message for misfired archive cache failed: %s", err.Error())
			}

			originalReader, err := storePort.Download(key)
			if err != nil {
				return nil, "", err
			}
			defer originalReader.Close()

			return ResizerPort.ResizeImage(originalReader, cachedWidth, false)
		},
		func(reader io.ReadCloser, size int, mediaType string, err error) ([]byte, string, error) {
			defer func() {
				if reader != nil {
					_ = reader.Close()
				}
			}()

			if err != nil {
				return nil, "", err
			}

			var content []byte
			if width < cachedWidth {
				content, _, err = ResizerPort.ResizeImage(reader, width, true)
				size = len(content)
			} else if maxBytes == 0 || size <= maxBytes {
				content, err = ioutil.ReadAll(reader)
			}
			if err != nil {
				return nil, "", err
			}

			if maxBytes > 0 && size > maxBytes {
				return nil, mediaType, MediaOverflowError
			}

			return content, mediaType, nil
		},
	)
}

// GetResizedImageURL returns a pre-signed URL to download the resized image ; GetResizedImage must have been called before.
func GetResizedImageURL(owner, mediaId string, width int) (string, error) {
	cachedWidth, err := findCacheableSize(width)
	if err != nil {
		return "", err
	}

	cacheId := generateCacheId(owner, mediaId, cachedWidth)
	return cachePort.SignedURL(cacheId, DownloadUrlValidityDuration)
}

func findCacheableSize(width int) (int, error) {
	if len(CacheableWidths) == 0 {
		return width, nil
	}

	i := 0
	for i+1 < len(CacheableWidths) && width <= CacheableWidths[i+1] {
		i++
	}

	if width > CacheableWidths[i] {
		return 0, errors.Errorf("width %d is not supported ; maximum size is %d", width, CacheableWidths[0])
	}
	return CacheableWidths[i], nil
}

func generateCacheId(owner, id string, width int) string {
	size := fmt.Sprintf("w=%d", width)
	if width == MiniatureCachedWidth {
		size = "miniatures"
	}

	return strings.Join([]string{size, owner, id}, "/")
}
