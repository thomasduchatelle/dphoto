package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	"testing"
	"time"
)

type inmemoryMedia struct {
	filename string
	size     uint
}

type inmemoryMediaWithHash struct {
	*inmemoryMedia
	hash string
}

var mediaDate = time.Date(2021, 4, 27, 10, 16, 22, 0, time.UTC)

func Test_analyseMedia(t *testing.T) {
	a := assert.New(t)
	mockImageDetailsReader := new(MockImageDetailsReaderAdapter)
	ImageDetailsReader = mockImageDetailsReader

	medias := []model.FoundMedia{
		newInmemoryMedia("/somewhere/my_image.jpg", 42),
		&inmemoryMediaWithHash{newInmemoryMedia("/somewhere/my_video.AVI", 42).(*inmemoryMedia), "qwerty"},
		newInmemoryMedia("/somewhere/my_document.txt", 42),
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
	mockImageDetailsReader.On("ReadImageDetails", mock.Anything, mediaDate).Once().Return(details, nil)

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
		got, err := analyseMedia(tt.args.found)
		if !tt.wantErr && a.NoError(err, tt.name) {
			a.Equal(tt.want, got, tt.name)
		} else if tt.wantErr && !a.Error(err) {
			log.Errorln(err)
		}
	}
}

func newInmemoryMedia(name string, size uint) model.FoundMedia {
	return &inmemoryMedia{name, size}
}

func (i *inmemoryMedia) Filename() string {
	return i.filename
}

func (i *inmemoryMedia) LastModificationDate() time.Time {
	return mediaDate
}

func (i *inmemoryMedia) SimpleSignature() *model.SimpleMediaSignature {
	return &model.SimpleMediaSignature{
		RelativePath: i.filename,
		Size:         i.size,
	}
}

func (i *inmemoryMedia) ReadMedia() (io.Reader, error) {
	return os.Open("../test_resources/scan/london_skyline_southbank.jpg")
}

func (m *inmemoryMediaWithHash) Sha256Hash() string {
	return m.hash
}
