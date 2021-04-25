package images

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestFileWithoutExif(t *testing.T) {
	a := assert.New(t)

	exifAdapter := new(ExifReader)
	reader, err := os.Open("../../../test_resources/scan/golang-logo.jpeg")
	if !a.NoError(err) {
		panic(err.Error())
	}

	lastModificationDate := time.Date(2021, 04, 25, 16, 40, 0, 0, time.UTC)
	details, err := exifAdapter.ReadImageDetails(reader, lastModificationDate)
	if a.NoError(err) {
		a.Equal(&model.MediaDetails{
			Width:       700,
			Height:      307,
			Orientation: model.OrientationUpperLeft,
			DateTime:    lastModificationDate,
		}, details)
	}
}

func TestFileWithExif(t *testing.T) {
	a := assert.New(t)

	exifAdapter := new(ExifReader)
	reader, err := os.Open("../../../test_resources/scan/london_skyline_southbank.jpg")
	if !a.NoError(err) {
		panic(err.Error())
	}

	lastModificationDate := time.Date(2021, 04, 25, 16, 40, 0, 0, time.UTC)
	details, err := exifAdapter.ReadImageDetails(reader, lastModificationDate)

	if a.NoError(err) {
		a.Equal(&model.MediaDetails{
			Width:        4048,
			Height:       3036,
			DateTime:     time.Unix(1574694084, 0).UTC(),
			Orientation:  model.OrientationUpperLeft,
			Make:         "Google",
			Model:        "Pixel",
			GPSLatitude:  51.50363055555555,
			GPSLongitude: -0.11583333333333334,
		}, details)
	}
}
