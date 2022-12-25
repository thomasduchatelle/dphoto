package image_resize

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"image"
	"io"
)

const (
	resizeMaxDimensionBreakPoint = 360
)

func NewResizer() *Resizer {
	return new(Resizer)
}

type Resizer struct{}

func (r Resizer) ResizeImageAtDifferentWidths(reader io.Reader, widths []int) (map[int][]byte, string, error) {
	resized := make(map[int][]byte)

	img, format, err := readImage(reader)
	if err != nil {
		return nil, "", err
	}

	for _, width := range widths {
		resizedImage := resizeImage(img, width, false)

		encodingFormat, err := imaging.FormatFromExtension(format)
		if err != nil {
			return nil, "", errors.Wrapf(err, "failed to find format from extention '%s'", format)
		}

		dest := bytes.NewBuffer(nil)
		err = imaging.Encode(dest, resizedImage, encodingFormat)

		resized[width] = dest.Bytes()
	}

	return resized, "image/" + format, err
}

func (r Resizer) ResizeImage(reader io.Reader, width int, fast bool) ([]byte, string, error) {
	// ResizeImage downscales (or upscales) the dimensions of an image to fit the requested width.
	return ResizeImage(reader, width, fast)
}

// ResizeImage downscales (or upscales) the dimensions of an image to fit the requested width.
func ResizeImage(reader io.Reader, width int, fast bool) ([]byte, string, error) {
	img, format, err := readImage(reader)
	if err != nil {
		return nil, "", err
	}

	resized := resizeImage(img, width, fast)

	encodingFormat, err := imaging.FormatFromExtension(format)
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to find format from extention '%s'", format)
	}

	dest := bytes.NewBuffer(nil)
	err = imaging.Encode(dest, resized, encodingFormat)
	return dest.Bytes(), "image/" + format, err
}

func resizeImage(img image.Image, width int, fast bool) image.Image {
	resizedWidth := width
	resizedHeight := 0
	if width > resizeMaxDimensionBreakPoint && img.Bounds().Dx() < img.Bounds().Dy() {
		// portrait images weight are HUGE when resized by their small dimension,
		// but it's ok for miniatures that are used on mobile phone
		resizedWidth = 0
		resizedHeight = width
	}

	if resizedWidth > img.Bounds().Dx() || resizedHeight > img.Bounds().Dy() {
		return img
	}

	algorithm := imaging.Lanczos
	if fast {
		algorithm = imaging.Box
	}

	return imaging.Resize(img, resizedWidth, resizedHeight, algorithm)
}

func readImage(reader io.Reader) (image.Image, string, error) {
	onHold := bytes.NewBuffer(nil)
	teeReader := io.TeeReader(reader, onHold)

	_, format, err := image.DecodeConfig(teeReader)
	if err != nil {
		return nil, "", errors.Wrapf(err, "couldn't determine image format from its content")
	}

	img, err := imaging.Decode(io.MultiReader(bytes.NewReader(onHold.Bytes()), reader), imaging.AutoOrientation(true))

	return img, format, errors.Wrapf(err, "failed to decode the image")
}
