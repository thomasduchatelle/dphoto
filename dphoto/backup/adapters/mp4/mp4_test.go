package mp4

import (
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
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

func TestMp4DetailsExtraction_android8(t *testing.T) {
	a := assert.New(t)

	exifAdapter := &Parser{
		Debug: true,
	}
	reader, err := os.Open("../../../test_resources/scan/MOV_0152.mp4")

	if !a.NoError(err) {
		panic(err.Error())
	}

	details, err := exifAdapter.ReadDetails(reader, backupmodel.DetailsReaderOptions{})
	if a.NoError(err) {
		fmt.Printf("Parsed date is %s\n", details.DateTime.Format("2006-01-02 15:04:05"))
		a.Equal(&backupmodel.MediaDetails{
			DateTime:      time.Date(2018, 3, 6, 15, 22, 25, 0, time.UTC),
			Duration:      13071,
			Height:        1080,
			VideoEncoding: "MP4",
			Width:         1920,
			GPSLatitude:   0.2952,
			GPSLongitude:  42.8127,
		}, details)
	}
}

func TestParseISO6709(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name    string
		gps     string
		wantLon float64
		wantLat float64
	}{
		{"it should parse simple GPS lon / lat", "+48.8577+002.295/", 48.8577, 2.295},
		{"it should not parse an invalid coordinates", "POLE NORTH > Santa Village", 0, 0},
		{"it should parse DDMM.M format", "+4851.462+00217.7/", 48.8577, 2.295},
		{"it should parse DDMMSS.S format", "+485127.72+0021742/", 48.8577, 2.295},
	}

	for _, tt := range tests {
		gotLon, gotLat := parseISO6709(tt.gps)
		a.Equal(tt.wantLon, gotLon, tt.name)
		a.Equal(tt.wantLat, gotLat, tt.name)
	}

}
