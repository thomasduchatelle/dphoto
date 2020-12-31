package images

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFileWithoutExif(t *testing.T) {
	a := assert.New(t)

	reader := new(exifReader)

	details, err := reader.ReadImageDetails("../../../test_resources/scan/golang-logo.jpeg")

	if a.NoError(err) {
		a.Equal(&backup.MediaDetails{
			Width:       700,
			Height:      307,
			Orientation: backup.UPPER_LEFT,
		}, details)
	}
}

func TestFileWithExif(t *testing.T) {
	a := assert.New(t)

	reader := new(exifReader)

	details, err := reader.ReadImageDetails("../../../test_resources/scan/london_skyline_southbank.jpg")

	if a.NoError(err) {
		a.Equal(&backup.MediaDetails{
			Width:        4048,
			Height:       3036,
			DateTime:     time.Unix(1574694084, 0).UTC(),
			Orientation:  backup.UPPER_LEFT,
			Make:         "Google",
			Model:        "Pixel",
			GPSLatitude:  51.50363055555555,
			GPSLongitude: -0.11583333333333334,
		}, details)
	}
}
