package analyser

import (
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/backup/interactors"
	"duchatelle.io/dphoto/dphoto/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var mediaDate = time.Date(2021, 4, 27, 10, 16, 22, 0, time.UTC)

func Test_analyseMedia(t *testing.T) {
	a := assert.New(t)
	mockImageDetailsReader := new(mocks.DetailsReaderAdapter)
	interactors.DetailsReaders[interactors.DetailsReaderTypeImage] = mockImageDetailsReader

	medias := []backupmodel.FoundMedia{
		backupmodel.NewInmemoryMedia("/somewhere/my_image.jpg", 42, mediaDate),
		backupmodel.NewInmemoryMediaWithHash("/somewhere/my_video.AVI", 42, mediaDate, "qwerty"),
		backupmodel.NewInmemoryMedia("/somewhere/my_document.txt", 42, mediaDate),
	}

	details := &backupmodel.MediaDetails{
		Width:        1024,
		Height:       768,
		DateTime:     time.Now(),
		Orientation:  backupmodel.OrientationUpperLeft,
		Make:         "Google",
		Model:        "Pixel 1",
		GPSLatitude:  0.0001,
		GPSLongitude: 0.0002,
	}
	mockImageDetailsReader.On("ReadDetails", mock.Anything, backupmodel.DetailsReaderOptions{Fast: true}).Once().Return(details, nil)

	type args struct {
		found backupmodel.FoundMedia
	}
	tests := []struct {
		name    string
		args    args
		want    *backupmodel.AnalysedMedia
		wantErr bool
	}{
		{"it should extract EXIF values from images, and compute a hash", args{medias[0]}, &backupmodel.AnalysedMedia{
			FoundMedia: medias[0],
			Type:       backupmodel.MediaTypeImage,
			Signature:  &backupmodel.FullMediaSignature{Sha256: "07b9bc44acdbbc0926117bb9e284f953060b2da0259b703af3def3841c7f61e8", Size: 42},
			Details:    details,
		}, false},

		{"it should not extract from video, and use pre-computed hash", args{medias[1]}, &backupmodel.AnalysedMedia{
			FoundMedia: medias[1],
			Type:       backupmodel.MediaTypeVideo,
			Signature:  &backupmodel.FullMediaSignature{Sha256: "qwerty", Size: 42},
			Details:    &backupmodel.MediaDetails{DateTime: mediaDate},
		}, false},

		{"it should not extract from other, and compute a hash", args{medias[2]}, &backupmodel.AnalysedMedia{
			FoundMedia: medias[2],
			Type:       backupmodel.MediaTypeOther,
			Signature:  &backupmodel.FullMediaSignature{Sha256: "07b9bc44acdbbc0926117bb9e284f953060b2da0259b703af3def3841c7f61e8", Size: 42},
			Details:    &backupmodel.MediaDetails{DateTime: mediaDate},
		}, false},
	}

	for _, tt := range tests {
		got, err := AnalyseMedia(tt.args.found)
		if !tt.wantErr && a.NoError(err, tt.name) {
			a.Equal(tt.want, got, tt.name)
		} else if tt.wantErr && !a.Error(err) {
			log.Errorln(err)
		}
	}
}
