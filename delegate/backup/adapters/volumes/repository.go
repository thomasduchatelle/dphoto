package volumes

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

type FileSystemRepository struct {
	Directory string
}

func (r *FileSystemRepository) RestoreLastSnapshot(volumeId string) ([]model.SimpleMediaSignature, error) {
	file, err := os.Open(path.Join(r.Directory, volumeId+".json"))
	if err != nil {
		log.WithField("VolumeId", volumeId).Info("No previous snapshot to load.")
		return nil, nil
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var signatures []model.SimpleMediaSignature
	err = json.Unmarshal(content, &signatures)
	return signatures, err
}

func (r *FileSystemRepository) StoreSnapshot(volumeId string, backupId string, signatures []model.SimpleMediaSignature) error {
	snapshotPath := path.Join(r.Directory, volumeId+".json")
	log.WithField("VolumeId", volumeId).Infof("Storing snapshot in '%s'...", backupId)

	err := os.MkdirAll(path.Dir(snapshotPath), 0766)
	if err != nil {
		return err
	}

	content, err := json.Marshal(signatures)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(snapshotPath, content, 0644)
}
