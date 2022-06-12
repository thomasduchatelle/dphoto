package backup

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/catalogmodel"
	"image"
	"time"
)

var (
	MediaNotFoundError = errors.Errorf("media not found")
)

var (
	Storage StorageAdapter
)

type StorageAdapter interface {
	FetchFile(owner string, folderName, filename string) ([]byte, error)
	ContentSignedUrl(owner string, folderName string, filename string, expires time.Duration) (string, error)
}

// GetMediaContent is an abomination trying to be the bridge between the catalog domain and the backup domain (temporary since May 2022)
func GetMediaContent(owner string, locations []*catalogmodel.MediaLocation, width int) ([]byte, string, error) {
	start := time.Now()
	content, err := fetchContent(owner, locations)
	if err != nil {
		return nil, "", err
	}

	fetched := time.Now()

	if width == 0 {
		return content, "", nil
	}

	resized, mediaType, err := resizeImage(content, width)

	compressed := time.Now()

	log.WithField("Location", locations[0].Filename).Infof("timing %-6d %s", fetched.Sub(start).Milliseconds(), "fetching")
	log.WithField("Location", locations[0].Filename).Infof("timing %-6d %s", compressed.Sub(fetched).Milliseconds(), "scaling")

	return resized, mediaType, err
}

func resizeImage(content []byte, width int) ([]byte, string, error) {
	img, err := imaging.Decode(bytes.NewReader(content), imaging.AutoOrientation(true))
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to decode the image")
	}

	resized := imaging.Resize(img, width, 0, imaging.Box) // todo should use imaging.Lanczos for cached version

	_, format, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil {
		return nil, "", errors.Wrapf(err, "couldn't determine image format from its content")
	}
	encodingFormat, err := imaging.FormatFromExtension(format)
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to find format from extention %s", format)
	}
	dest := bytes.NewBuffer(nil)
	err = imaging.Encode(dest, resized, encodingFormat)
	return dest.Bytes(), "image/" + format, err
}

func fetchContent(owner string, locations []*catalogmodel.MediaLocation) (content []byte, err error) {
	for _, location := range locations {
		content, err = Storage.FetchFile(owner, location.FolderName, location.Filename)

		if err == nil {
			return
		}
	}

	return
}
