package filter

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"duchatelle.io/dphoto/dphoto/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var mediaDate = time.Date(2021, 4, 27, 10, 16, 22, 0, time.UTC)

func Test_filter_filter(t *testing.T) {
	a := assert.New(t)

	mockVolumeRepository := new(mocks.VolumeRepositoryAdapter)
	model.VolumeRepositoryPort = mockVolumeRepository

	// given
	mockVolumeRepository.On("RestoreLastSnapshot", "volume-1").Return([]model.SimpleMediaSignature{
		{RelativePath: "image_002.jpg", Size: 42},
		{RelativePath: "image_003.jpg", Size: 12},
	}, nil)

	mediaFilter, err := NewMediaFilter(&model.VolumeToBackup{
		UniqueId: "volume-1",
		Type:     model.VolumeTypeFileSystem,
		Path:     "/somewhere",
		Local:    false,
	})
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	// and
	type args struct {
		found model.FoundMedia
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"it should keep media with different name and size", args{model.NewInmemoryMedia("image_001.jpg", 1, mediaDate)}, true},
		{"it should keep media with different name", args{model.NewInmemoryMedia("image_001.jpg", 42, mediaDate)}, true},
		{"it should keep media with same name but different size", args{model.NewInmemoryMedia("image_002.jpg", 1, mediaDate)}, true},
		{"it should filter out medias matching both name and size", args{model.NewInmemoryMedia("image_002.jpg", 42, mediaDate)}, false},
	}

	// when - then
	for _, tt := range tests {
		got := mediaFilter.Filter(tt.args.found)
		a.Equal(tt.want, got, tt.name)
	}
}
