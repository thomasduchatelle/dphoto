package filesystem

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"sort"
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	a := assert.New(t)

	fsHandler := &FsHandler{
		ImageExtensions: []string{".jpeg"},
		VideoExtensions: []string{".TxT"},
	}

	media := make(chan backup.FoundMedia, 42)
	err := fsHandler.FindMediaRecursively("../../test_resources", media)
	close(media)

	if a.NoError(err) {
		var found []backup.FoundMedia
		for m := range media {
			found = append(found, m)
		}
		sort.Slice(found, func(i, j int) bool {
			return found[i].LocalAbsolutePath < found[j].LocalAbsolutePath
		})

		if a.Len(found, 2) {
			a.Equal(found[0].Type, backup.VIDEO)
			a.True(strings.HasSuffix(found[0].LocalAbsolutePath, "a_text.TXT"), found[0].LocalAbsolutePath, "has suffix")
			a.Equal(backup.SimpleMediaSignature{
				RelativePath: "scan/a_text.TXT",
				Size:         6,
			}, found[0].SimpleSignature)

			a.Equal(found[1].Type, backup.IMAGE)
			a.Equal(backup.SimpleMediaSignature{
				RelativePath: "scan/golang-logo.jpeg",
				Size:         22601,
			}, found[1].SimpleSignature)
		}
	}
}

func TestFsHandler_CopyToLocal(t *testing.T) {
	a := assert.New(t)

	fsHandler := FsHandler{
		DirMode: 0744,
	}

	localCache := path.Join(os.TempDir(), "dphoto_filesystem")
	_ = os.RemoveAll(localCache)

	mediaHash, err := fsHandler.CopyToLocal("../../test_resources/scan/golang-logo.jpeg", path.Join(localCache, "logo.jpg"))

	if a.NoError(err) {
		a.Equal("2921a1e238f8fd3231a8b04c8bd1433e49ef66e3737c49cd01d34ca4cd4a97ad", mediaHash)
	}
}
