package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_filter_filter(t *testing.T) {
	a := assert.New(t)

	mockVolumeRepository := new(MockVolumeRepositoryAdapter)
	VolumeRepository = mockVolumeRepository

	// given
	mockVolumeRepository.On("RestoreLastSnapshot", "volume-1").Return([]model.SimpleMediaSignature{
		{RelativePath: "image_002.jpg", Size: 42},
		{RelativePath: "image_003.jpg", Size: 12},
	}, nil)

	mediaFilter, err := newMediaFilter(&model.VolumeToBackup{
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
		{"it should keep media with different name and size", args{newInmemoryMedia("image_001.jpg", 1)}, true},
		{"it should keep media with different name", args{newInmemoryMedia("image_001.jpg", 42)}, true},
		{"it should keep media with same name but different size", args{newInmemoryMedia("image_002.jpg", 1)}, true},
		{"it should filter out medias matching both name and size", args{newInmemoryMedia("image_002.jpg", 42)}, false},
	}

	// when - then
	for _, tt := range tests {
		got := mediaFilter.Filter(tt.args.found)
		a.Equal(tt.want, got, tt.name)
	}
}
