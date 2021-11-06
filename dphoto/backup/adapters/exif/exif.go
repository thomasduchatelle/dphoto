// Package exif parse image files to extract key details.
package exif

import (
	"bytes"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"github.com/pkg/errors"
	"github.com/rwcarlsen/goexif/exif"
	log "github.com/sirupsen/logrus"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"
)

func init() {
	// note - Canon parser is failing on 2007 photos from a Canon camera
	exif.RegisterParsers()
}

type Parser struct{}

func (p *Parser) Supports(media backupmodel.FoundMedia, mediaType backupmodel.MediaType) bool {
	return mediaType == backupmodel.MediaTypeImage
}

func (p *Parser) ReadDetails(reader io.Reader, options backupmodel.DetailsReaderOptions) (*backupmodel.MediaDetails, error) {
	buffer := bytes.NewBuffer(nil)
	teeReader := io.TeeReader(reader, buffer)

	x, err := exif.Decode(teeReader)
	if err != nil {
		log.WithError(err).Warn("no EXIF data found in file, try another way")
		return p.readImageWithoutExif(io.MultiReader(buffer, reader))
	}

	latitude, longitude, err := x.LatLong()
	if err != nil {
		latitude = 0
		longitude = 0
	}

	return &backupmodel.MediaDetails{
		Width:        p.getIntOrIgnore(x, exif.ImageWidth),
		Height:       p.getIntOrIgnore(x, exif.ImageLength),
		Orientation:  p.readOrientation(x),
		DateTime:     p.readDateTime(x, time.Time{}),
		Make:         p.getStringOrIgnore(x, exif.Make),
		Model:        p.getStringOrIgnore(x, exif.Model),
		GPSLatitude:  latitude,
		GPSLongitude: longitude,
	}, nil
}

func (p *Parser) readOrientation(x *exif.Exif) backupmodel.ImageOrientation {
	switch p.getIntOrIgnore(x, exif.Orientation) {
	case 3:
		return backupmodel.OrientationLowerRight
	case 6:
		return backupmodel.OrientationUpperRight
	case 8:
		return backupmodel.OrientationLowerLeft
	default:
		return backupmodel.OrientationUpperLeft
	}
}

func (p *Parser) readDateTime(x *exif.Exif, defaultDate time.Time) time.Time {
	datetime := p.getStringOrIgnore(x, exif.DateTime)
	if datetime == "" {
		datetime = p.getStringOrIgnore(x, exif.DateTimeOriginal)
	}

	if datetime != "" {
		exifTime, err := time.Parse("2006:01:02 15:04:05", datetime)
		if err == nil {
			return exifTime.UTC()
		} else {
			log.WithField("MediaAnalyser", "Exif").Warnf("Unsupported dfate format: %s", datetime)
		}
	}

	return defaultDate
}

func (p *Parser) getStringOrIgnore(x *exif.Exif, model exif.FieldName) string {
	if t, err := x.Get(model); err == nil && t != nil {
		if val, err := t.StringVal(); err == nil {
			return val
		}
	}
	return ""
}

func (p *Parser) getIntOrIgnore(x *exif.Exif, model exif.FieldName) int {
	if t, err := x.Get(model); err == nil && t != nil && t.Count > 0 {
		if val, err := t.Int(0); err == nil {
			return val
		}
	}
	return 0
}

func (p *Parser) getFloatOrIgnore(x *exif.Exif, model exif.FieldName) float64 {
	if t, err := x.Get(model); err == nil && t != nil && t.Count > 0 {
		if val, err := t.Float(0); err == nil {
			return val
		}
	}
	return 0
}

func (p *Parser) readImageWithoutExif(reader io.Reader) (*backupmodel.MediaDetails, error) {
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't extract image dimentions")
	}

	return &backupmodel.MediaDetails{
		Width:       img.Width,
		Height:      img.Height,
		Orientation: backupmodel.OrientationUpperLeft,
	}, nil
}
