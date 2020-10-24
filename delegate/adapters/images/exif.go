package images

import (
	"duchatelle.io/dphoto/delegate/backup"
	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	log "github.com/sirupsen/logrus"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"
)

func init() {
	exif.RegisterParsers(mknote.All...)

	backup.ImageDetailsReader = new(exifReader)
}

type exifReader struct{}

func (e *exifReader) ReadImageDetails(imagePath string) (*backup.MediaDetails, error) {
	withContext := log.WithField("ImagePath", imagePath)

	file, err := os.Open(imagePath)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file '%s'", imagePath)
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		withContext.WithError(err).Warn("no EXIF data found in file, try another way")
		return e.readImageWithoutExif(imagePath)
	}

	latitude, longitude, err := x.LatLong()
	if err != nil {
		latitude = 0
		longitude = 0
	}

	return &backup.MediaDetails{
		Width:        e.getIntOrIgnore(x, exif.ImageWidth),
		Height:       e.getIntOrIgnore(x, exif.ImageLength),
		Orientation:  e.readOrientation(x),
		DateTime:     e.readDateTime(x, imagePath),
		Make:         e.getStringOrIgnore(x, exif.Make),
		Model:        e.getStringOrIgnore(x, exif.Model),
		GPSLatitude:  latitude,
		GPSLongitude: longitude,
	}, nil
}

func (e *exifReader) readOrientation(x *exif.Exif) backup.ImageOrientation {
	switch e.getIntOrIgnore(x, exif.Orientation) {
	case 3:
		return backup.LOWER_RIGHT
	case 6:
		return backup.UPPER_RIGHT
	case 8:
		return backup.LOWER_LEFT
	default:
		return backup.UPPER_LEFT
	}
}

func (e *exifReader) readDateTime(x *exif.Exif, imagePath string) time.Time {
	datetime := e.getStringOrIgnore(x, exif.DateTime)
	if datetime != "" {
		exifTime, err := time.Parse("2006:01:02 15:04:05", datetime)
		if err == nil {
			return exifTime.UTC()
		}
	}

	if stat, err := os.Stat(imagePath); err == nil {
		return stat.ModTime()
	}

	return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
}

func (e *exifReader) getStringOrIgnore(x *exif.Exif, model exif.FieldName) string {
	if t, err := x.Get(model); err == nil && t != nil {
		if val, err := t.StringVal(); err == nil {
			return val
		}
	}
	return ""
}

func (e *exifReader) getIntOrIgnore(x *exif.Exif, model exif.FieldName) int {
	if t, err := x.Get(model); err == nil && t != nil && t.Count > 0 {
		if val, err := t.Int(0); err == nil {
			return val
		}
	}
	return 0
}

func (e *exifReader) getFloatOrIgnore(x *exif.Exif, model exif.FieldName) float64 {
	if t, err := x.Get(model); err == nil && t != nil && t.Count > 0 {
		if val, err := t.Float(0); err == nil {
			return val
		}
	}
	return 0
}

func (e *exifReader) readImageWithoutExif(imagePath string) (*backup.MediaDetails, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open image file '%s' to extract its dimensions", imagePath)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't read image '%s' to extract its dimensions", imagePath)
	}

	return &backup.MediaDetails{
		Width:       img.Width,
		Height:      img.Height,
		Orientation: backup.UPPER_LEFT,
	}, nil
}
