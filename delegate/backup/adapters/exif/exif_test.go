package exif

import (
	"github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestFileWithoutExif(t *testing.T) {
	a := assert.New(t)

	exifAdapter := new(Parser)
	reader, err := os.Open("../../../test_resources/scan/golang-logo.jpeg")
	if !a.NoError(err) {
		panic(err.Error())
	}

	details, err := exifAdapter.ReadDetails(reader, backupmodel.DetailsReaderOptions{})
	if a.NoError(err) {
		a.Equal(&backupmodel.MediaDetails{
			Width:       700,
			Height:      307,
			Orientation: backupmodel.OrientationUpperLeft,
		}, details)
	}
}

func TestFileWithExif(t *testing.T) {
	a := assert.New(t)

	exifAdapter := new(Parser)
	reader, err := os.Open("../../../test_resources/scan/london_skyline_southbank.jpg")
	if !a.NoError(err) {
		panic(err.Error())
	}

	details, err := exifAdapter.ReadDetails(reader, backupmodel.DetailsReaderOptions{})

	if a.NoError(err) {
		a.Equal(&backupmodel.MediaDetails{
			Width:        4048,
			Height:       3036,
			DateTime:     time.Unix(1574694084, 0).UTC(),
			Orientation:  backupmodel.OrientationUpperLeft,
			Make:         "Google",
			Model:        "Pixel",
			GPSLatitude:  51.50363055555555,
			GPSLongitude: -0.11583333333333334,
		}, details)
	}
}
