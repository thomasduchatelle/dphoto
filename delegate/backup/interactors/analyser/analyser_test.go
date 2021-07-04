package analyser

import (
	"duchatelle.io/dphoto/dphoto/backup/interactors"
	"duchatelle.io/dphoto/dphoto/backup/model"
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

	medias := []model.FoundMedia{
		model.NewInmemoryMedia("/somewhere/my_image.jpg", 42, mediaDate),
		model.NewInmemoryMediaWithHash("/somewhere/my_video.AVI", 42, mediaDate, "qwerty"),
		model.NewInmemoryMedia("/somewhere/my_document.txt", 42, mediaDate),
	}

	details := &model.MediaDetails{
		Width:        1024,
		Height:       768,
		DateTime:     time.Now(),
		Orientation:  model.OrientationUpperLeft,
		Make:         "Google",
		Model:        "Pixel 1",
		GPSLatitude:  0.0001,
		GPSLongitude: 0.0002,
	}
	mockImageDetailsReader.On("ReadDetails", mock.Anything, model.DetailsReaderOptions{Fast: true}).Once().Return(details, nil)

	type args struct {
		found model.FoundMedia
	}
	tests := []struct {
		name    string
		args    args
		want    *model.AnalysedMedia
		wantErr bool
	}{
		{"it should extract EXIF values from images, and compute a hash", args{medias[0]}, &model.AnalysedMedia{
			FoundMedia: medias[0],
			Type:       model.MediaTypeImage,
			Signature:  &model.FullMediaSignature{Sha256: "07b9bc44acdbbc0926117bb9e284f953060b2da0259b703af3def3841c7f61e8", Size: 42},
			Details:    details,
		}, false},

		{"it should not extract from video, and use pre-computed hash", args{medias[1]}, &model.AnalysedMedia{
			FoundMedia: medias[1],
			Type:       model.MediaTypeVideo,
			Signature:  &model.FullMediaSignature{Sha256: "qwerty", Size: 42},
			Details:    &model.MediaDetails{DateTime: mediaDate},
		}, false},

		{"it should not extract from other, and compute a hash", args{medias[2]}, &model.AnalysedMedia{
			FoundMedia: medias[2],
			Type:       model.MediaTypeOther,
			Signature:  &model.FullMediaSignature{Sha256: "07b9bc44acdbbc0926117bb9e284f953060b2da0259b703af3def3841c7f61e8", Size: 42},
			Details:    &model.MediaDetails{DateTime: mediaDate},
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
