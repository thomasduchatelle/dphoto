// Package localstorage provides a temporary location where files can be copied to be quickly read for analyse.
package localstorage

import (
	"context"
	"crypto/sha256"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

// LocalStorage stores remote files locally so they can be analysed quickly
// TODO - keep files under 10MB in a 512MB in-memory cache
type LocalStorage struct {
	localMediaPath string
	semaphore      *semaphore.Weighted
}

type localMedia struct {
	store    *LocalStorage
	delegate backupmodel.FoundMedia
	path     string
	sha256   string
}

func (m *localMedia) MediaPath() backupmodel.MediaPath {
	return m.delegate.MediaPath()
}

func (m *localMedia) LastModificationDate() time.Time {
	return m.delegate.LastModificationDate()
}

func (m *localMedia) SimpleSignature() *backupmodel.SimpleMediaSignature {
	return m.delegate.SimpleSignature()
}

func (m *localMedia) ReadMedia() (io.ReadCloser, error) {
	return os.Open(m.path)
}

func (m *localMedia) Sha256Hash() string {
	return m.sha256
}

func (m *localMedia) String() string {
	return fmt.Sprint(m.delegate) + " [local=" + m.path + "]"
}

func (m *localMedia) Close() error {
	log.Infof("Releasing %d for %s", m.SimpleSignature().Size, m)
	err := os.Remove(m.path)
	m.store.release(m.SimpleSignature().Size)

	if err != nil {
		log.WithError(err).Warnf("Failed to release temporary file %s", m.path)
	}
	return errors.Wrapf(err, "Failed to release temporary file %s", m.path)
}

func NewLocalStorage(localDir string, bufferAreaSizeInBytes int) (*LocalStorage, error) {
	cleanedDir, err := filepath.Abs(os.ExpandEnv(localDir))
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(cleanedDir, 0744)
	return &LocalStorage{
		localMediaPath: cleanedDir,
		semaphore:      semaphore.NewWeighted(int64(bufferAreaSizeInBytes)),
	}, err
}

func (l *LocalStorage) DownloadMedia(found backupmodel.FoundMedia) (backupmodel.FoundMedia, error) {
	err := l.take(found.SimpleSignature().Size)
	if err != nil {
		return nil, err
	}

	reader, err := found.ReadMedia()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	key := path.Join(l.localMediaPath, uuid.New().String()+path.Ext(found.MediaPath().Filename))
	log.Debugf("DownloaderPort > download locally %s to %s", found, key)

	writer, err := os.Create(key)
	if err != nil {
		return nil, err
	}

	hash := sha256.New()

	_, err = io.Copy(io.MultiWriter(writer, hash), reader)
	return &localMedia{
		store:    l,
		delegate: found,
		path:     key,
		sha256:   hex.EncodeToString(hash.Sum(nil)),
	}, err
}

func (l *LocalStorage) take(size uint) error {
	return l.semaphore.Acquire(context.TODO(), int64(size))
}

func (l *LocalStorage) release(size uint) {
	l.semaphore.Release(int64(size))
}
