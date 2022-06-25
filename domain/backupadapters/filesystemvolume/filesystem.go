// Package filesystemvolume scan a local filesystem to find medias in it
package filesystemvolume

import (
	"github.com/dixonwille/skywalker"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"io"
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

func (v *volume) FindMedias() ([]backup.FoundMedia, error) {
	// TODO - Could be simplified with https://pkg.go.dev/io/fs#WalkDir
	absRootPath, err := filepath.Abs(v.path)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid volume path")
	}

	extensions := make([]string, len(v.supportedExtensions)*2)
	index := 0
	for typ, _ := range v.supportedExtensions {
		extensions[index*2] = "." + strings.ToLower(typ)
		extensions[index*2+1] = "." + strings.ToUpper(typ)
		index++
	}

	worker := v.newWorker(absRootPath)
	walker := skywalker.New(absRootPath, worker)
	walker.ExtListType = skywalker.LTWhitelist
	walker.ExtList = extensions

	medias, caughtErrors, done := v.collect(log.WithField("mediaPath", v.path), worker.medias, worker.errors)

	err = walker.Walk()
	worker.walkComplete()
	<-done

	err = v.wrapErrors(err, *caughtErrors)

	return *medias, err
}

func (v *volume) Children(path backup.MediaPath) (backup.SourceVolume, error) {
	return New(path.ParentFullPath), nil
}

func (v *volume) newWorker(rootPath string) *fsWorker {
	return &fsWorker{
		rootPath: rootPath,
		medias:   make(chan backup.FoundMedia, 256),
		errors:   make(chan error, 32),
	}
}

func (v *volume) collect(mcd *log.Entry, mediaChannel chan backup.FoundMedia, errorsChannel chan error) (*[]backup.FoundMedia, *[]error, chan interface{}) {
	medias := make([]backup.FoundMedia, 0)
	caughtErrors := make([]error, 0)
	done := make(chan interface{})

	go func() {
		defer close(done)

		closed := 0
		for closed < 0x11 {
			select {
			case err, more := <-errorsChannel:
				if more {
					mcd.WithError(err).Warnln("scan error")
					caughtErrors = append(caughtErrors, err)
				} else {
					closed = closed | 0x01
				}

			case media, more := <-mediaChannel:
				if more {
					medias = append(medias, media)
				} else {
					closed = closed | 0x10
				}
			}
		}
	}()

	return &medias, &caughtErrors, done
}

func (v *volume) wrapErrors(err error, caughtErrors []error) error {
	if err != nil {
		caughtErrors = append([]error{err}, caughtErrors...)
	}

	if len(caughtErrors) > 0 {
		var messages []string
		for _, e := range caughtErrors[1:] {
			messages = append(messages, e.Error())
		}

		return errors.Wrapf(err, "scanning caused %d errors (%s)", len(caughtErrors), strings.Join(messages, ", "))
	}

	return nil
}

type fsWorker struct {
	rootPath string
	medias   chan backup.FoundMedia
	errors   chan error
}

func (w *fsWorker) Work(mediaPath string) {
	for _, p := range strings.Split(mediaPath, "/") {
		if strings.HasPrefix(p, ".") && p != "." && p != ".." {
			// skip hidden files
			return
		}
	}

	abs, err := filepath.Abs(mediaPath)
	if err != nil {
		w.errors <- errors.Wrapf(err, "invalid path, no absolute path")
		return
	}

	stat, err := os.Stat(abs)
	if err != nil {
		w.errors <- errors.Wrapf(err, "file stats cannot be read")
		return
	}

	w.medias <- &fsMedia{
		volumeAbsolutePath:   w.rootPath,
		absolutePath:         abs,
		size:                 int(stat.Size()),
		lastModificationDate: stat.ModTime(),
	}
}

func (w *fsWorker) walkComplete() {
	close(w.medias)
	close(w.errors)
}

type fsMedia struct {
	volumeAbsolutePath   string
	absolutePath         string
	size                 int
	lastModificationDate time.Time
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
