package daemon

import (
	"github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"
)

var (
	VolumeManager VolumeManagerPort
)

type VolumeManagerPort interface {
	OnMountedVolume(volume backupmodel.VolumeToBackup)
	OnUnMountedVolume(uuid string)
}
