package mp4

import (
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMp4DetailsExtraction(t *testing.T) {
	a := assert.New(t)

	exifAdapter := &Parser{
		Debug: true,
	}
	reader, err := os.Open("../../../test_resources/scan/MOV_0897.mp4")

	if !a.NoError(err) {
		panic(err.Error())
	}

	details, err := exifAdapter.ReadDetails(reader, backupmodel.DetailsReaderOptions{})
	if a.NoError(err) {
		fmt.Printf("Parsed date is %s\n", details.DateTime.Format("2006-01-02 15:04:05"))
		a.Equal(&backupmodel.MediaDetails{
			DateTime:      time.Date(2019, 8, 18, 9, 57, 55, 0, time.UTC),
			Duration:      1801,
			Height:        1080,
			VideoEncoding: "MP4",
			Width:         1920,
		}, details)
	}
}
