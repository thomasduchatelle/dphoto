package daemon

var (
	VolumeManager VolumeManagerPort
)

type VolumeManagerPort interface {
	OnMountedVolume(volume VolumeToBackup)
	OnUnMountedVolume(uuid string)
}

type VolumeToBackup struct {
	UniqueId string
	Path     string
}
