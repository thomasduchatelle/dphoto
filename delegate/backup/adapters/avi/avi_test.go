package avi

import (
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"time"
)

func TestMp4DetailsExtraction_exiftool(t *testing.T) {
	a := assert.New(t)

	exifAdapter := &Parser{
		Debug: true,
	}
	var reader io.Reader
	reader, err := os.Open("../../../test_resources/scan/t_images_RIFF.avi")
	if !a.NoError(err) {
		panic(err.Error())
	}

	details, err := exifAdapter.ReadDetails(reader, backupmodel.DetailsReaderOptions{})
	if a.NoError(err) {
		fmt.Printf("Parsed date is %s\n", details.DateTime.Format("2006-01-02 15:04:05"))
		a.Equal(&backupmodel.MediaDetails{
			DateTime:      time.Date(2003, 3, 10, 15, 4, 43, 0, time.UTC),
			Duration:      15533,
			Height:        240,
			VideoEncoding: "mjpg",
			Width:         320,
		}, details)
	}
}
