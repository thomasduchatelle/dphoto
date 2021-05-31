package filesystem

import (
	"duchatelle.io/dphoto/dphoto/backup/interactors/analyser"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/dixonwille/skywalker"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	media     chan model.FoundMedia
	count     int64
	sizeSum   int64
}

type fsMedia struct {
	absolutePath         string
	size                 int
	lastModificationDate time.Time
	relativePath         string
}

func (f *FsHandler) FindMediaRecursively(volume model.VolumeToBackup, medias chan model.FoundMedia) (uint, uint, error) {
	worker, err := f.newWorker(volume.Path, medias)
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

	abs, err := filepath.Abs(mediaPath)
	if err != nil {
		logContext.WithError(err).Warnln("Invalid path, no absolute path")
		return
	}
	rel, err := filepath.Rel(w.mountPath, mediaPath)
	if err != nil {
		logContext.WithError(err).WithField("mountPath", w.mountPath).Warnln("Media path is not relative.")
		return
	}

	stat, err := os.Stat(abs)
	if err != nil {
		logContext.WithError(err).Warnln("Failed getting stats for the file")
		return
	}

	w.media <- &fsMedia{
		absolutePath:         abs,
		size:                 int(stat.Size()),
		lastModificationDate: stat.ModTime(),
		relativePath:         rel,
	}
	atomic.AddInt64(&w.count, 1)
	atomic.AddInt64(&w.sizeSum, stat.Size())
}

func (f *FsHandler) newWorker(mountPath string, media chan model.FoundMedia) (*fsWorker, error) {
	absMountPath, err := filepath.Abs(mountPath)
	return &fsWorker{
		mountPath: absMountPath,
		media:     media,
	}, errors.Wrapf(err, "can't get the absolute path of %s", mountPath)
}

func (f *fsMedia) Filename() string {
	return path.Base(f.absolutePath)
}

func (f *fsMedia) LastModificationDate() time.Time {
	return f.lastModificationDate
}

func (f *fsMedia) SimpleSignature() *model.SimpleMediaSignature {
	return &model.SimpleMediaSignature{
		RelativePath: f.relativePath,
		Size:         uint(f.size),
	}
}

func (f *fsMedia) ReadMedia() (io.Reader, error) {
	return os.Open(f.absolutePath)
}

func (f *fsMedia) String() string {
	return f.absolutePath
}
