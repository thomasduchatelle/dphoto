package image_resize

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"image"
	"io"
)

func NewResizer() *Resizer {
	return new(Resizer)
}

type Resizer struct{}

func (r Resizer) ResizeImage(reader io.Reader, width int, fast bool) ([]byte, string, error) {
	// ResizeImage downscales (or upscales) the dimensions of an image to fit the requested width.
	return ResizeImage(reader, width, fast)
}

// ResizeImage downscales (or upscales) the dimensions of an image to fit the requested width.
func ResizeImage(reader io.Reader, width int, fast bool) ([]byte, string, error) {
	onHold := bytes.NewBuffer(nil)
	teeReader := io.TeeReader(reader, onHold)

	_, format, err := image.DecodeConfig(teeReader)
	if err != nil {
		return nil, "", errors.Wrapf(err, "couldn't determine image format from its content")
	}

	img, err := imaging.Decode(io.MultiReader(bytes.NewReader(onHold.Bytes()), reader), imaging.AutoOrientation(true))
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to decode the image")
	}

	algorithm := imaging.Lanczos
	if fast {
		algorithm = imaging.Box
	}
	resized := imaging.Resize(img, width, 0, algorithm)

	encodingFormat, err := imaging.FormatFromExtension(format)
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to find format from extention '%s'", format)
	}

	dest := bytes.NewBuffer(nil)
	err = imaging.Encode(dest, resized, encodingFormat)
	return dest.Bytes(), "image/" + format, err
}
