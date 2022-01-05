// Package filesystem scan a local filesystem to find medias on it
package filesystem

import (
	"github.com/dixonwille/skywalker"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/interactors/analyser"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

type FsHandler struct{}

type fsWorker struct {
	mountPath string
	callback  func(backupmodel.FoundMedia)
	count     int64
	sizeSum   int64
}

type fsMedia struct {
	volumeAbsolutePath   string
	absolutePath         string
	size                 int
	lastModificationDate time.Time
}

func (f *FsHandler) FindMediaRecursively(volume backupmodel.VolumeToBackup, callback func(backupmodel.FoundMedia)) (uint, uint, error) {
	worker, err := f.newWorker(volume.Path, callback)
	if err != nil {
		return 0, 0, err
	}

	extensions := make([]string, len(analyser.SupportedExtensions)*2)
	index := 0
	for typ, _ := range analyser.SupportedExtensions {
		extensions[index*2] = "." + strings.ToLower(typ)
		extensions[index*2+1] = "." + strings.ToUpper(typ)
		index++
	}

	walker := skywalker.New(volume.Path, worker)
	walker.ExtListType = skywalker.LTWhitelist
	walker.ExtList = extensions

	return uint(worker.count), uint(worker.sizeSum), errors.Wrapf(walker.Walk(), "failed scanning path %s", volume.Path)
}

func (w *fsWorker) Work(mediaPath string) {
	logContext := log.WithField("mediaPath", mediaPath)

	for _, p := range strings.Split(mediaPath, "/") {
		if strings.HasPrefix(p, ".") && p != "." && p != ".." {
			// skip hidden files
			return
		}
	}

	abs, err := filepath.Abs(mediaPath)
	if err != nil {
		logContext.WithError(err).Warnln("Invalid path, no absolute path")
		return
	}

	stat, err := os.Stat(abs)
	if err != nil {
		logContext.WithError(err).Warnln("Failed getting stats for the file")
		return
	}

	w.callback(&fsMedia{
		volumeAbsolutePath:   w.mountPath,
		absolutePath:         abs,
		size:                 int(stat.Size()),
		lastModificationDate: stat.ModTime(),
	})
	atomic.AddInt64(&w.count, 1)
	atomic.AddInt64(&w.sizeSum, stat.Size())
}

func (f *FsHandler) newWorker(mountPath string, callback func(backupmodel.FoundMedia)) (*fsWorker, error) {
	absMountPath, err := filepath.Abs(mountPath)
	return &fsWorker{
		mountPath: absMountPath,
		callback:  callback,
	}, errors.Wrapf(err, "can't get the absolute path of %s", mountPath)
}

func (f *fsMedia) MediaPath() backupmodel.MediaPath {
	return backupmodel.MediaPath{
		ParentFullPath: path.Dir(f.absolutePath),
		Root:           f.volumeAbsolutePath,
		Path:           strings.TrimPrefix(strings.TrimPrefix(path.Dir(f.absolutePath), f.volumeAbsolutePath), "/"),
		Filename:       path.Base(f.absolutePath),
		ParentDir:      path.Base(path.Dir(f.absolutePath)),
	}
}

func (f *fsMedia) LastModificationDate() time.Time {
	return f.lastModificationDate
}

func (f *fsMedia) SimpleSignature() *backupmodel.SimpleMediaSignature {
	mediaPath := f.MediaPath()
	return &backupmodel.SimpleMediaSignature{
		RelativePath: path.Join(mediaPath.Path, mediaPath.Filename),
		Size:         uint(f.size),
	}
}

func (f *fsMedia) ReadMedia() (io.ReadCloser, error) {
	return os.Open(f.absolutePath)
}

func (f *fsMedia) String() string {
	return f.absolutePath
}
