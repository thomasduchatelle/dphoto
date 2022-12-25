package archive

import (
	"github.com/pkg/errors"
	"io"
	"time"
)

const (
	MiniatureCachedWidth     = 360  // MiniatureCachedWidth is the minimum size in which images are stored. Under that, the MiniatureCachedWidth is stored and the image will be re-scaled down on the fly
	MediumQualityCachedWidth = 2400 // MediumQualityCachedWidth is the highest cacheable resolution, consumer should not request above that
)

var (
	NotFoundError      = errors.New("media is not present in the archive")
	MediaOverflowError = errors.New("media at the requested width is bigger that what the consumer can support")
	CacheableWidths    = []int{MediumQualityCachedWidth, MiniatureCachedWidth} // CacheableWidths are the only resolution cached, array must be sorted DESC.

	supportedExtensionsForResizing = map[string]interface{}{
		".jpg":  nil,
		".jpeg": nil,
		".png":  nil,
	}
)

// DestructuredKey indicates the preferred key: Prefix + Suffix. A counter can be added between the 2 to make the name unique.
type DestructuredKey struct {
	Prefix string
	Suffix string
}

// StoreRequest is used to archive a media
type StoreRequest struct {
	DateTime         time.Time                     // DateTime is used to name the final file
	FolderName       string                        // FolderName is the location where the file must be physically stored
	Id               string                        // Id is the unique identifier from 'catalog'
	Open             func() (io.ReadCloser, error) // Open creates a new reader to this file
	OriginalFilename string                        // OriginalFilename is used to preserve the right extension
	Owner            string                        // Owner is he tenant to which this media belongs
	SignatureSha256  string                        // SignatureSha256 is the hash of file content, mainly used to generate the filename
}

// ImageToResize is a request to cache resized image ; optional arguments might be used when processed synchronously but not when queued.
type ImageToResize struct {
	Owner    string                        // Owner is mandatory
	MediaId  string                        // MediaId is mandatory
	StoreKey string                        // StoreKey is optional
	Widths   []int                         // Widths must have at least 1 value
	Open     func() (io.ReadCloser, error) // Open is optional
}
