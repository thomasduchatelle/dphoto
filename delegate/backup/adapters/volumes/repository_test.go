package volumes

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestStoreAndRetrieve(t *testing.T) {
	a := assert.New(t)

	signatures := []model.SimpleMediaSignature{
		{
			RelativePath: "/somewhere/file-1.jpg",
			Size:         42,
		},
		{
			RelativePath: "/somewhere/file-2.jpg",
			Size:         12,
		},
	}

	repository := FileSystemRepository{Directory: path.Join(os.TempDir(), "dphoto/volumeRepository")}
	err := repository.StoreSnapshot("volume-1", "backupid-1", signatures)
	if a.NoError(err) {
		snapshot, err := repository.RestoreLastSnapshot("volume-1")
		if a.NoError(err) {
			a.Equal(signatures, snapshot, "it should restore the list of signature of the previous run")
		}
	}
}

func TestRestoreEmpty(t *testing.T) {
	a := assert.New(t)

	repository := FileSystemRepository{Directory: path.Join(os.TempDir(), "dphoto/volumeRepository")}
	snapshot, err := repository.RestoreLastSnapshot("volume-2")
	if a.NoError(err) {
		a.Empty(snapshot, "it should return an empty snapshot when the volume wasn't known")
	}
}
