package backupproxy

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
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
	VolumeId    string
	ScanResult  []*backup.ScannedFolder
	RejectCount int
}

func init() {
	config.Listen(func(cfg config.Config) {
		dir := cfg.GetStringOrDefault(config.LocalHome, os.ExpandEnv("$HOME/.dphoto"))

		storeDir, err := filepath.Abs(dir)
		if err != nil {
			panic(err)
		}

		storeFile = path.Join(storeDir, "last_scan.json")

		err = os.MkdirAll(path.Dir(storeFile), 0744)
		if err != nil {
			panic(err)
		}
	})
}

func Store(volumeId string, folders []*backup.ScannedFolder, rejectCount int) error {
	if storeFile == "" {
		return errors.Errorf("local.home must have been set before using this function.")
	}

	jsonValue, err := json.Marshal(stateContent{
		VolumeId:    volumeId,
		ScanResult:  folders,
		RejectCount: rejectCount,
	})
	if err != nil {
		return err
	}

	return ioutil.WriteFile(storeFile, jsonValue, 0644)
}

func restore(volumeId string) ([]*backup.ScannedFolder, int, error) {
	content, err := ioutil.ReadFile(storeFile)
	if err != nil && os.IsNotExist(err) {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	}

	state := stateContent{}
	err = json.Unmarshal(content, &state)
	if err != nil {
		return nil, 0, err
	}

	if state.VolumeId != volumeId {
		return nil, 0, nil
	}
	return state.ScanResult, state.RejectCount, nil
}
