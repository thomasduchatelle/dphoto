package scanner

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var mediaDate = time.Date(2021, 4, 27, 10, 16, 22, 0, time.UTC)

func Test_analyseMedia(t *testing.T) {
	a := assert.New(t)
	mockImageDetailsReader := new(MockImageDetailsReaderAdapter)
	ImageDetailsReader = mockImageDetailsReader

	medias := []FoundMedia{
		NewInmemoryMedia("/somewhere/my_image.jpg", 42, mediaDate),
		&InmemoryMediaWithHash{NewInmemoryMedia("/somewhere/my_video.AVI", 42, mediaDate).(*InmemoryMedia), "qwerty"},
		NewInmemoryMedia("/somewhere/my_document.txt", 42, mediaDate),
	}

	details := &MediaDetails{
		Width:        1024,
		Height:       768,
		DateTime:     time.Now(),
		Orientation:  OrientationUpperLeft,
		Make:         "Google",
		Model:        "Pixel 1",
		GPSLatitude:  0.0001,
		GPSLongitude: 0.0002,
	}
	mockImageDetailsReader.On("ReadImageDetails", mock.Anything, mediaDate).Once().Return(details, nil)

	type args struct {
		found FoundMedia
	}
	tests := []struct {
		name    string
		args    args
		want    *AnalysedMedia
		wantErr bool
	}{
		{"it should extract EXIF values from images, and compute a hash", args{medias[0]}, &AnalysedMedia{
			FoundMedia: medias[0],
			Type:       MediaTypeImage,
			Signature:  &FullMediaSignature{Sha256: "07b9bc44acdbbc0926117bb9e284f953060b2da0259b703af3def3841c7f61e8", Size: 42},
			Details:    details,
		}, false},

		{"it should not extract from video, and use pre-computed hash", args{medias[1]}, &AnalysedMedia{
			FoundMedia: medias[1],
			Type:       MediaTypeVideo,
			Signature:  &FullMediaSignature{Sha256: "qwerty", Size: 42},
			Details:    &MediaDetails{DateTime: mediaDate},
		}, false},

		{"it should not extract from other, and compute a hash", args{medias[2]}, &AnalysedMedia{
			FoundMedia: medias[2],
			Type:       MediaTypeOther,
			Signature:  &FullMediaSignature{Sha256: "07b9bc44acdbbc0926117bb9e284f953060b2da0259b703af3def3841c7f61e8", Size: 42},
			Details:    &MediaDetails{DateTime: mediaDate},
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
