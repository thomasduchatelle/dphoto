// Package filesystemvolume scans a local filesystem to find medias in it
package filesystemvolume

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func New(path string) backup.SourceVolume {
	return &volume{
		path:                path,
		supportedExtensions: backup.SupportedExtensions,
	}
}

type volume struct {
	path                string
	supportedExtensions map[string]backup.MediaType
}

func (v *volume) String() string {
	return v.path
}

func (v *volume) FindMedias(context.Context) ([]backup.FoundMedia, error) {
	absRootPath, err := filepath.Abs(v.path)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid volume path")
	}

	var medias []backup.FoundMedia
	err = filepath.WalkDir(absRootPath, func(filePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if absRootPath != filePath && strings.HasPrefix(entry.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
		if _, ok := v.supportedExtensions[ext]; !ok {
			return nil
		}

		medias = append(medias, &fsMedia{
			volumeAbsolutePath:   absRootPath,
			absolutePath:         filePath,
			size:                 int(info.Size()),
			lastModificationDate: info.ModTime(),
		})
		return nil
	})

	return medias, errors.Wrapf(err, "walking directory %s", v.path)
}

func (v *volume) Children(path backup.MediaPath) (backup.SourceVolume, error) {
	return New(path.ParentFullPath), nil
}

type fsMedia struct {
	volumeAbsolutePath   string
	absolutePath         string
	size                 int
	lastModificationDate time.Time
}

func (f *fsMedia) LastModification() time.Time {
	return f.lastModificationDate
}

func (f *fsMedia) Size() int {
	return f.size
}

func (f *fsMedia) MediaPath() backup.MediaPath {
	return backup.MediaPath{
		ParentFullPath: path.Dir(f.absolutePath),
		Root:           f.volumeAbsolutePath,
		Path:           strings.TrimPrefix(strings.TrimPrefix(path.Dir(f.absolutePath), f.volumeAbsolutePath), "/"),
		Filename:       path.Base(f.absolutePath),
		ParentDir:      path.Base(path.Dir(f.absolutePath)),
	}
}

func (f *fsMedia) ReadMedia() (io.ReadCloser, error) {
	return os.Open(f.absolutePath)
}

func (f *fsMedia) String() string {
	return f.absolutePath
}
