package images

import (
	"bytes"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	log "github.com/sirupsen/logrus"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

type ExifReader struct{}

func (e *ExifReader) ReadImageDetails(reader io.Reader, lastModifiedDate time.Time) (*model.MediaDetails, error) {
	buffer := bytes.NewBuffer(nil)
	teeReader := io.TeeReader(reader, buffer)

	x, err := exif.Decode(teeReader)
	if err != nil {
		log.WithField("Reader", reader).WithError(err).Warn("no EXIF data found in file, try another way")
		return e.readImageWithoutExif(io.MultiReader(buffer, reader), lastModifiedDate)
	}

	latitude, longitude, err := x.LatLong()
	if err != nil {
		latitude = 0
		longitude = 0
	}

	return &model.MediaDetails{
		Width:        e.getIntOrIgnore(x, exif.ImageWidth),
		Height:       e.getIntOrIgnore(x, exif.ImageLength),
		Orientation:  e.readOrientation(x),
		DateTime:     e.readDateTime(x, lastModifiedDate),
		Make:         e.getStringOrIgnore(x, exif.Make),
		Model:        e.getStringOrIgnore(x, exif.Model),
		GPSLatitude:  latitude,
		GPSLongitude: longitude,
	}, nil
}

func (e *ExifReader) readOrientation(x *exif.Exif) model.ImageOrientation {
	switch e.getIntOrIgnore(x, exif.Orientation) {
	case 3:
		return model.OrientationLowerRight
	case 6:
		return model.OrientationUpperRight
	case 8:
		return model.OrientationLowerLeft
	default:
		return model.OrientationUpperLeft
	}
}

func (e *ExifReader) readDateTime(x *exif.Exif, lastModifiedDate time.Time) time.Time {
	datetime := e.getStringOrIgnore(x, exif.DateTime)
	if datetime != "" {
		exifTime, err := time.Parse("2006:01:02 15:04:05", datetime)
		if err == nil {
			return exifTime.UTC()
		}
	}

	return lastModifiedDate
}

func (e *ExifReader) getStringOrIgnore(x *exif.Exif, model exif.FieldName) string {
	if t, err := x.Get(model); err == nil && t != nil {
		if val, err := t.StringVal(); err == nil {
			return val
		}
	}
	return ""
}

func (e *ExifReader) getIntOrIgnore(x *exif.Exif, model exif.FieldName) int {
	if t, err := x.Get(model); err == nil && t != nil && t.Count > 0 {
		if val, err := t.Int(0); err == nil {
			return val
		}
	}
	return 0
}

func (e *ExifReader) getFloatOrIgnore(x *exif.Exif, model exif.FieldName) float64 {
	if t, err := x.Get(model); err == nil && t != nil && t.Count > 0 {
		if val, err := t.Float(0); err == nil {
			return val
		}
	}
	return 0
}

func (e *ExifReader) readImageWithoutExif(reader io.Reader, lastModifiedDate time.Time) (*model.MediaDetails, error) {
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't extract image dimentions")
	}

	return &model.MediaDetails{
		Width:       img.Width,
		Height:      img.Height,
		Orientation: model.OrientationUpperLeft,
		DateTime:    lastModifiedDate,
	}, nil
}
