// Package storefs use local filesystem to store medias locally
package storefs

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

func NewStore(storeDir string) (archive.StoreAdapter, error) {
	err := os.MkdirAll(storeDir, 0744)

	abs, _ := filepath.Abs(storeDir)
	return &store{
		rootDir: abs,
	}, errors.Wrapf(err, "creating %s directory", storeDir)
}

type store struct {
	rootDir string
}

func (s *store) Download(key string) (io.ReadCloser, error) {
	return os.Open(path.Join(s.rootDir, key))
}

func (s *store) Upload(keyHint archive.DestructuredKey, reader io.Reader) (string, error) {
	absMediaPath, key := s.findUniqueFileName(keyHint)

	writer, err := os.Create(absMediaPath)
	if err != nil {
		return "", errors.Wrapf(err, "creating file %s", absMediaPath)
	}

	_, err = io.Copy(writer, reader)
	return key, err
}

func (s *store) Copy(origin string, destination archive.DestructuredKey) (string, error) {
	absDestinationPath, key := s.findUniqueFileName(destination)

	absSourcePath := path.Join(s.rootDir, origin)
	reader, err := os.Open(absSourcePath)
	if err != nil {
		return "", errors.Wrapf(err, "reading %s", origin)
	}

	writer, err := os.Create(absDestinationPath)
	if err != nil {
		return "", errors.Wrapf(err, "creating %s", absDestinationPath)
	}

	_, err = io.Copy(writer, reader)
	return key, errors.Wrapf(err, "copy %s -> %s", absSourcePath, absDestinationPath)
}

func (s *store) Delete(locations []string) error {
	for _, location := range locations {
		absPath := path.Join(s.rootDir, location)
		err := os.Remove(absPath)
		if err != nil && !os.IsNotExist(err) {
			return errors.Wrapf(err, "removing %s", absPath)
		}
	}

	return nil
}

func (s *store) SignedURL(key string, duration time.Duration) (string, error) {
	return "", errors.Errorf("SignedURL is not supported by filesystem implementation")
}

func (s *store) findUniqueFileName(keyHint archive.DestructuredKey) (string, string) {
	relativePath := keyHint.Prefix + keyHint.Suffix

	count := 1
	for _, err := os.Stat(path.Join(s.rootDir, relativePath)); !os.IsNotExist(err); _, err = os.Stat(path.Join(s.rootDir, relativePath)) {
		relativePath = fmt.Sprintf("%s_%02d%s", keyHint.Prefix, count, keyHint.Suffix)
		count++
	}

	return path.Join(s.rootDir, relativePath), relativePath
}
