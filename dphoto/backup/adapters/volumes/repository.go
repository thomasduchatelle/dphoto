// Package volumes is storing snapshot of the last backup so most media do not need to be re-analysed on next backup.
package volumes

import (
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

type FileSystemRepository struct {
	Directory string
}

func (r *FileSystemRepository) RestoreLastSnapshot(volumeId string) ([]backupmodel.SimpleMediaSignature, error) {
	storageFile := r.getStorageFile(volumeId)
	mdc := log.WithFields(log.Fields{
		"VolumeId":    volumeId,
		"StorageFile": storageFile,
	})

	file, err := os.Open(storageFile)
	if err != nil && os.IsNotExist(err) {
		mdc.Infoln("No previous snapshot to load.")
		return nil, nil

	} else if err != nil {
		return nil, errors.Wrapf(err, "Couldn't load previous volume snapshot from file %s", storageFile)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var signatures []backupmodel.SimpleMediaSignature
	err = json.Unmarshal(content, &signatures)

	mdc.Debugf("FileSystemRepository > restored snaphot with %d medias", len(signatures))
	return signatures, err
}

func (r *FileSystemRepository) StoreSnapshot(volumeId string, backupId string, signatures []backupmodel.SimpleMediaSignature) error {
	storageFile := r.getStorageFile(volumeId)
	mdc := log.WithFields(log.Fields{
		"VolumeId":    volumeId,
		"StorageFile": storageFile,
		"BackupID":    backupId,
	})
	mdc.Infof("FileSystemRepository > Storing volume snapshot with %d signatures", len(signatures))

	err := os.MkdirAll(path.Dir(storageFile), 0766)
	if err != nil {
		return err
	}

	content, err := json.Marshal(signatures)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(storageFile, content, 0644)
}

func (r *FileSystemRepository) getStorageFile(volumeId string) string {
	unsafeChar := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	safeVolumeId := strings.Trim(unsafeChar.ReplaceAllString(volumeId, "_"), "_")
	return os.ExpandEnv(path.Join(r.Directory, safeVolumeId+".json"))
}
