package backup_adapter

import (
	"duchatelle.io/dphoto/delegate/backup"
	"duchatelle.io/dphoto/delegate/daemon"
)

func init() {
	daemon.VolumeManager = new(backupAdapter)
}

type backupAdapter struct{}

func (b *backupAdapter) OnMountedVolume(volume backup.RemovableVolume) {
	backup.OnMountedVolume(volume)
}

func (b *backupAdapter) OnUnMountedVolume(uuid string) {
	backup.OnUnMountedVolume(uuid)
}
