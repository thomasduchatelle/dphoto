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
	cacheKey := generateCacheId(owner, mediaId, width)

	// note: retention and storage class is managed at infrastructure level
	return NewCache().GetOrStore(
		cacheKey,
		func() ([]byte, string, error) {
			cachedSize := width
			if cachedSize < MiniatureCachedWidth {
				cachedSize = MiniatureCachedWidth
			}

			key, err := repositoryPort.FindById(owner, mediaId)
			if err != nil {
				return nil, "", err
			}

			log.WithFields(log.Fields{
				"Owner":    owner,
				"Width":    width,
				"StoreKey": key,
			}).Infof("Missed archive cache")

			err = asyncJobPort.WarmUpCacheByFolder(owner, key, cachedSize)
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

			return ResizerPort.ResizeImage(originalReader, cachedSize, false)
		},
		func(reader io.ReadCloser, size int, mediaType string, err error) ([]byte, string, error) {
			defer func() {
				if reader != nil {
					_ = reader.Close()
				}
			}()

			switch {
			case err != nil:
				return nil, "", err

			case width < MiniatureCachedWidth:
				return ResizerPort.ResizeImage(reader, width, true)

			case maxBytes > 0 && size > maxBytes:
				return nil, mediaType, MediaOverflowError

			default:
				content, err := ioutil.ReadAll(reader)
				return content, mediaType, errors.Wrapf(err, "failed to read content from key '%s'", cacheKey)
			}
		},
	)
}

// GetResizedImageURL returns a pre-signed URL to download the resized image ; GetResizedImage must have been called before.
func GetResizedImageURL(owner, mediaId string, width int) (string, error) {
	cacheId := generateCacheId(owner, mediaId, width)
	return cachePort.SignedURL(cacheId, DownloadUrlValidityDuration)
}

func generateCacheId(owner, id string, width int) string {
	size := fmt.Sprintf("w=%d", width)
	if width <= 400 {
		size = "miniatures"
	}

	return strings.Join([]string{size, owner, id}, "/")
}
