package filesystemvolume

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	a := assert.New(t)

	fs := &volume{
		path: "../../test_resources",
		supportedExtensions: map[string]backup.MediaType{
			"txt":  backup.MediaTypeOther,
			"Jpeg": backup.MediaTypeImage,
		},
	}
	abspath, _ := filepath.Abs(fs.path)

	medias, err := fs.FindMedias()

	sort.Slice(medias, func(i, j int) bool {
		return strings.Compare(medias[i].String(), medias[j].String()) < 0
	})

	if a.NoError(err) {
		relativePaths := make([]string, len(medias), len(medias))
		for i, media := range medias {
			relativePaths[i] = path.Join(media.MediaPath().Path, media.MediaPath().Filename)
		}

		a.Equal([]string{
			"a_text.TXT",
			"scan/folder1/another.txt",
			"scan/golang-logo-resized.jpeg",
			"scan/golang-logo.jpeg",
		}, relativePaths, "it should find all the files with matching extension, sub folders, while ignoring hidden files")

		name := "it should generate proper media path"
		a.Equal(backup.MediaPath{
			ParentFullPath: abspath,
			Root:           abspath,
			Path:           "",
			Filename:       "a_text.TXT",
			ParentDir:      "test_resources",
		}, medias[0].MediaPath(), name+" when in the root folder")
		a.Equal(backup.MediaPath{
			ParentFullPath: path.Join(abspath, "scan/folder1"),
			Root:           abspath,
			Path:           "scan/folder1",
			Filename:       "another.txt",
			ParentDir:      "folder1",
		}, medias[1].MediaPath(), name+" when in the sub-folder")

		name = "it should read the content of the file"
		a.Equal(6, medias[0].Size(), name)
		reader, err := medias[0].ReadMedia()
		if a.NoError(err, name) {
			content, err := ioutil.ReadAll(reader)
			if a.NoError(err, name) {
				a.Equal([]byte("Hello."), content, name)
			}
		}

		a.Equal(path.Join(abspath, "a_text.TXT"), medias[0].String(), "it should have String() returning the full URL")

		name = "it should generate a valid children"
		childVolume, err := fs.Children(medias[1].MediaPath())
		if a.NoError(err, name) {
			fsChildVolume := childVolume.(*volume)
			fsChildVolume.supportedExtensions = fs.supportedExtensions

			subFolderMedias, err := childVolume.FindMedias()
			if a.NoError(err, name) {
				if a.Len(subFolderMedias, 1, name) {
					a.Equal("another.txt", subFolderMedias[0].MediaPath().Filename, name)
				}

			}
		}
	}
}
