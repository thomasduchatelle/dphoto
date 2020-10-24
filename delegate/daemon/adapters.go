package daemon

import "duchatelle.io/dphoto/delegate/backup"

var (
	VolumeManager VolumeManagerPort
)

type VolumeManagerPort interface {
	OnMountedVolume(volume backup.RemovableVolume)
	OnUnMountedVolume(uuid string)
}
