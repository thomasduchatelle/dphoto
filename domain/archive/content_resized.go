package archive

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"strings"
)

const MiniatureCachedWidth = 360 // MiniatureCachedWidth is the minimum size in which images are stored. Under that, the MiniatureCachedWidth is stored and the image will be re-scaled down on the fly

// GetResizedImage returns the image in the requested size (or rounded up), and the media type.
func GetResizedImage(owner, mediaId string, width int, maxBytes int) ([]byte, string, error) {
	cacheKey := generateCacheId(owner, mediaId, width)

	// note: retention and storage class is managed at infrastructure level
	return NewCache().GetOrStore(
		cacheKey,
		func() ([]byte, string, error) {
			key, err := repositoryPort.FindById(owner, mediaId)
			if err != nil {
				return nil, "", err
			}

			originalReader, err := storePort.Download(key)
			if err != nil {
				return nil, "", err
			}
			defer originalReader.Close()

			cachedSize := width
			if cachedSize < MiniatureCachedWidth {
				cachedSize = MiniatureCachedWidth
			}
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

func generateCacheId(owner, id string, width int) string {
	size := fmt.Sprintf("w=%d", width)
	if width <= 400 {
		size = "miniatures"
	}

	return strings.Join([]string{size, owner, id}, "/")
}
