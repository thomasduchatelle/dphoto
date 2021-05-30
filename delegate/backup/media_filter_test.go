package backup

import (
	"duchatelle.io/dphoto/dphoto/scanner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_filter_filter(t *testing.T) {
	a := assert.New(t)

	mockVolumeRepository := new(MockVolumeRepositoryAdapter)
	VolumeRepository = mockVolumeRepository

	// given
	mockVolumeRepository.On("RestoreLastSnapshot", "volume-1").Return([]scanner.SimpleMediaSignature{
		{RelativePath: "image_002.jpg", Size: 42},
		{RelativePath: "image_003.jpg", Size: 12},
	}, nil)

	mediaFilter, err := newMediaFilter(&scanner.VolumeToBackup{
		UniqueId: "volume-1",
		Type:     scanner.VolumeTypeFileSystem,
		Path:     "/somewhere",
		Local:    false,
	})
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	// and
	type args struct {
		found scanner.FoundMedia
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"it should keep media with different name and size", args{scanner.NewInmemoryMedia("image_001.jpg", 1, mediaDate)}, true},
		{"it should keep media with different name", args{scanner.NewInmemoryMedia("image_001.jpg", 42, mediaDate)}, true},
		{"it should keep media with same name but different size", args{scanner.NewInmemoryMedia("image_002.jpg", 1, mediaDate)}, true},
		{"it should filter out medias matching both name and size", args{scanner.NewInmemoryMedia("image_002.jpg", 42, mediaDate)}, false},
	}

	// when - then
	for _, tt := range tests {
		got := mediaFilter.Filter(tt.args.found)
		a.Equal(tt.want, got, tt.name)
	}
}
