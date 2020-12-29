package filesystem

import (
	"crypto/sha256"
	"duchatelle.io/dphoto/dphoto/backup"
	"encoding/hex"
	"github.com/dixonwille/skywalker"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func init() {
	backup.FileHandler = &FsHandler{
		DirMode:         0744,
		ImageExtensions: []string{".png", ".jpg"},
		VideoExtensions: []string{".avi", ".mpeg", ".mkv", ".mts"},
	}
}

type FsHandler struct {
	DirMode         os.FileMode
	ImageExtensions []string
	VideoExtensions []string
}

type worker struct {
	mountPath string
	media     chan backup.FoundMedia
	imagesExt map[string]interface{}
	videoExt  map[string]interface{}
}

func (f *FsHandler) FindMediaRecursively(mountPath string, media chan backup.FoundMedia) error {
	worker, err := f.newWorker(mountPath, media)
	if err != nil {
		return err
	}

	walker := skywalker.New(mountPath, worker)
	walker.ExtListType = skywalker.LTWhitelist
	walker.ExtList = withUpper(f.ImageExtensions, f.VideoExtensions)

	return errors.Wrapf(walker.Walk(), "failed scanning path %s", mountPath)
}

func (w *worker) Work(mediaPath string) {
	logContext := log.WithField("mediaPath", mediaPath)

	abs, err := filepath.Abs(mediaPath)
	if err != nil {
		logContext.WithError(err).Errorln("Invalid path")
		return
	}
	rel, err := filepath.Rel(w.mountPath, mediaPath)
	if err != nil {
		logContext.WithError(err).WithField("mountPath", w.mountPath).Errorln("Media path is not relative")
		return
	}

	ext := strings.ToLower(path.Ext(mediaPath))
	media := backup.OTHER
	if _, ok := w.imagesExt[ext]; ok {
		media = backup.IMAGE
	} else if _, ok := w.videoExt[ext]; ok {
		media = backup.VIDEO
	}

	stat, err := os.Stat(abs)
	if err != nil {
		logContext.WithError(err).Errorf("Invalid media file.")
		return
	}

	w.media <- backup.FoundMedia{
		Type:              media,
		LocalAbsolutePath: abs,
		SimpleSignature: backup.SimpleMediaSignature{
			RelativePath: rel,
			Size:         int(stat.Size()),
		},
	}
}

func (f *FsHandler) CopyToLocal(originPath string, destPath string) (string, error) {
	mediaHash, err := f.doCopyToLocal(originPath, destPath)
	return mediaHash, errors.Wrapf(err, "failed to copy %s to %s", originPath, destPath)
}

func (f *FsHandler) doCopyToLocal(originPath string, destPath string) (mediaHash string, err error) {
	err = os.MkdirAll(path.Dir(destPath), f.DirMode)
	if err != nil {
		return
	}

	var src, targetFile *os.File
	src, err = os.Open(originPath)
	if err != nil {
		return
	}

	targetFile, err = os.Create(destPath)
	if err != nil {
		return
	}

	h := sha256.New()

	_, err = io.Copy(io.MultiWriter(targetFile, h), src)
	if err != nil {
		return
	}

	mediaHash = hex.EncodeToString(h.Sum(nil))
	return
}

func (f *FsHandler) newWorker(mountPath string, media chan backup.FoundMedia) (skywalker.Worker, error) {
	absMountPath, err := filepath.Abs(mountPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	worker := &worker{
		mountPath: absMountPath,
		media:     media,
		imagesExt: indexStrings(f.ImageExtensions),
		videoExt:  indexStrings(f.VideoExtensions),
	}

	for _, e := range f.ImageExtensions {
		worker.imagesExt[e] = nil
	}

	for _, e := range f.VideoExtensions {
		worker.videoExt[e] = nil
	}

	return worker, nil
}

func indexStrings(values []string) map[string]interface{} {
	m := make(map[string]interface{})

	for _, e := range values {
		m[strings.ToLower(e)] = nil
	}

	return m
}

// join several arrays and add the uppercase version
func withUpper(slices ...[]string) []string {
	size := 0
	for _, s := range slices {
		size += len(s)
	}

	values := make([]string, size*2)
	for _, s := range slices {
		for _, val := range s {
			values = append(values, strings.ToLower(val), strings.ToUpper(val))
		}
	}

	return values
}
