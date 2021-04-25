package daemon

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
)

var (
	VolumeManager VolumeManagerPort
)

type VolumeManagerPort interface {
	OnMountedVolume(volume model.VolumeToBackup)
	OnUnMountedVolume(uuid string)
}
