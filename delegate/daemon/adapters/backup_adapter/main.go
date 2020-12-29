package backup_adapter

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/daemon"
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
