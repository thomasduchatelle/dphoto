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
			"scan/a_text.TXT",
			"scan/another.txt",
			"scan/golang-logo-resized.jpeg",
			"scan/golang-logo.jpeg",
		}, relativePaths, "it should find all the files with matching extension, sub folders, while ignoring hiden files")

		a.Equal(backup.MediaPath{
			ParentFullPath: path.Join(abspath, "scan"),
			Root:           abspath,
			Path:           "scan",
			Filename:       "a_text.TXT",
			ParentDir:      "scan",
		}, medias[0].MediaPath())

		name := "it should read the content of the file"
		a.Equal(6, medias[0].Size(), name)
		reader, err := medias[0].ReadMedia()
		if a.NoError(err, name) {
			content, err := ioutil.ReadAll(reader)
			if a.NoError(err, name) {
				a.Equal([]byte("Hello."), content, name)
			}
		}

		a.Equal(path.Join(abspath, "scan/a_text.TXT"), medias[0].String(), "it should have String() returning the full URL")
	}
}
