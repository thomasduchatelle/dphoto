package filesystem

import (
	"bytes"
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/backup/interactors/analyser"
	"github.com/stretchr/testify/assert"
	"io"
	"path"
	"path/filepath"
	"sort"
	"testing"
)

func TestScanner(t *testing.T) {
	a := assert.New(t)

	fsHandler := new(FsHandler)

	mediaChannel := make(chan backupmodel.FoundMedia, 42)
	volumeMount := "../../../test_resources"
	volumeMountAbs, _ := filepath.Abs(volumeMount)

	analyser.SupportedExtensions = map[string]backupmodel.MediaType{
		"txt":  backupmodel.MediaTypeOther,
		"Jpeg": backupmodel.MediaTypeImage,
	}

	_, _, err := fsHandler.FindMediaRecursively(backupmodel.VolumeToBackup{
		UniqueId: volumeMount,
		Type:     backupmodel.VolumeTypeFileSystem,
		Path:     volumeMount,
		Local:    true,
	}, func(media backupmodel.FoundMedia) {
		mediaChannel <- media
	})
	close(mediaChannel)

	if a.NoError(err) {
		var found []*fsMedia
		for m := range mediaChannel {
			found = append(found, m.(*fsMedia))
		}
		sort.Slice(found, func(i, j int) bool {
			return found[i].absolutePath < found[j].absolutePath
		})

		if a.Len(found, 2) {
			a.Equal(path.Join(volumeMountAbs, "scan/a_text.TXT"), found[0].String())
			a.Equal("a_text.TXT", path.Base(found[0].Filename()))
			a.Equal(&backupmodel.SimpleMediaSignature{
				RelativePath: "scan/a_text.TXT",
				Size:         6,
			}, found[0].SimpleSignature())

			buffer := new(bytes.Buffer)
			content, err := found[0].ReadMedia()
			if a.NoError(err) {
				_, err = io.Copy(buffer, content)
				if a.NoError(err) {
					a.Equal("Hello.", buffer.String())
				}
			}

			a.Equal(&backupmodel.SimpleMediaSignature{
				RelativePath: "scan/golang-logo.jpeg",
				Size:         22601,
			}, found[1].SimpleSignature())
		}
	}
}
