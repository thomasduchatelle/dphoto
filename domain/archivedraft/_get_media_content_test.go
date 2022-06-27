package archivedraft

import (
	"crypto/md5"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/mocks"
	"io/ioutil"
	"os"
	"testing"
)

func TestResizing(t *testing.T) {
	a := assert.New(t)

	StorageMock := new(mocks.StorageAdapter)
	Storage = StorageMock
	src, err := ioutil.ReadFile("../../dphoto/test_resources/scan/golang-logo.jpeg")
	if !a.NoError(err) {
		return
	}
	expected, err := ioutil.ReadFile("../../dphoto/test_resources/scan/golang-logo-resized.jpeg")
	if !a.NoError(err) {
		return
	}

	StorageMock.On("FetchFile", "ironman", "tony", "stark").Return(src, nil)

	content, format, err := GetMediaContent("ironman", []*catalog.MediaLocation{
		{FolderName: "tony", Filename: "stark"},
	}, 100)

	if a.NoError(err) {
		a.Equal("image/jpeg", format, "it should use JPEG format")
		a.Equal(md5.Sum(expected), md5.Sum(content), "it should scale down the image")

		// for visual confirmation
		_ = os.MkdirAll("../../.build", 0744)
		err = os.WriteFile("../../.build/golang-logo-resized.jpeg", content, 0644)
		a.NoError(err)
	}
}
