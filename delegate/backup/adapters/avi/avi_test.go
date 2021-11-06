package avi

import (
	"bytes"
	"github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"time"
)

func TestMp4DetailsExtraction(t *testing.T) {
	a := assert.New(t)

	exifAdapter := &Parser{
		Debug: true,
	}
	var reader io.Reader
	reader, err := os.Open("../../../test_resources/scan/t_images_RIFF.avi")
	if !a.NoError(err) {
		panic(err.Error())
	}

	buffer := make([]byte, 1024)
	_, _ = io.ReadFull(reader, buffer)
	fmt.Printf("\nFile first 1MB:\n%s\n", hex.Dump(buffer))

	details, err := exifAdapter.ReadDetails(io.MultiReader(bytes.NewReader(buffer), reader), backupmodel.DetailsReaderOptions{})
	if a.NoError(err) {
		fmt.Printf("Parsed date is %s\n", details.DateTime.Format("2006-01-02 15:04:05"))
		a.Equal(&backupmodel.MediaDetails{
			DateTime:      time.Date(2003, 3, 10, 15, 4, 43, 0, time.UTC),
			Duration:      15533,
			Height:        240,
			Make:          "CanonMVI01",
			VideoEncoding: "mjpg",
			Width:         320,
		}, details)
	}
}

func TestMp4DetailsExtraction_2(t *testing.T) {
	a := assert.New(t)

	exifAdapter := &Parser{
		Debug: true,
	}
	var reader io.Reader
	reader, err := os.Open("../../../test_resources/scan/MOV02192_CANON.AVI")
	if !a.NoError(err) {
		panic(err.Error())
	}

	details, err := exifAdapter.ReadDetails(reader, backupmodel.DetailsReaderOptions{})
	if a.NoError(err) {
		fmt.Printf("Parsed date is %s\n", details.DateTime.Format("2006-01-02 15:04:05"))
		a.Equal(&backupmodel.MediaDetails{
			DateTime:      time.Date(2011, 4, 15, 15, 5, 34, 0, time.UTC),
			Duration:      7771,
			Height:        240,
			Make:          "SONY DSC MJPEG 0100",
			VideoEncoding: "\a",
			Width:         320,
		}, details)
	}
}
