package backupadaptor

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/config"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// TODO extract that as a key-value storage per volume ?
var (
	storeFile string
)

type stateContent struct {
	VolumeId   string
	ScanResult []*backup.FoundAlbum
}

func init() {
	config.Listen(func(cfg config.Config) {
		dir := cfg.GetStringOrDefault("local.home", os.ExpandEnv("$HOME/.dphoto"))

		storeDir, err := filepath.Abs(dir)
		if err != nil {
			panic(err)
		}

		storeFile = path.Join(storeDir, "last_scan.json")
	})
}

func Store(volumeId string, result []*backup.FoundAlbum) error {
	if storeFile == "" {
		return errors.Errorf("local.home must have been set before using this function.")
	}

	jsonValue, err := json.Marshal(stateContent{
		VolumeId:   volumeId,
		ScanResult: result,
	})
	if err != nil {
		return err
	}

	return ioutil.WriteFile(storeFile, jsonValue, 0644)
}

func restore(volumeId string) ([]*backup.FoundAlbum, error) {
	content, err := ioutil.ReadFile(storeFile)
	if err != nil && os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	state := stateContent{}
	err = json.Unmarshal(content, &state)
	if err != nil {
		return nil, err
	}

	if state.VolumeId != volumeId {
		return nil, nil
	}
	return state.ScanResult, nil
}
