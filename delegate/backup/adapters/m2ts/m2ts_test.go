package m2ts

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestM2TSReader(t *testing.T) {
	a := assert.New(t)

	detailsReader := &Parser{Debug: true}

	reader, err := os.Open("../../../test_resources/scan/sample.MTS")
	if a.NoError(err) {
		details, err := detailsReader.ReadDetails(reader, model.DetailsReaderOptions{})
		if a.NoError(err) {
			// debug - details
			fmt.Printf("Details: %+v\n", details)

			expectedDateTime, err := time.Parse("2006-01-02 15:04:05 -0700", "2019-08-16 11:03:11 +0100")
			a.NoError(err)
			a.Equal(expectedDateTime, details.DateTime, "it should read date from H264 video stream")
			a.Equal(int64(449), details.Duration, "it should compute duration from M2TS format")
			a.Equal("Sony", details.Make, "it should read Make from H264 video stream")
			a.Equal("DSC-RX100", details.Model, "it should read Model from H264 video stream (Sony tested only)")
		}
	}
}
