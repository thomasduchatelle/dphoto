package backup_adapter

import (
	"github.com/thomasduchatelle/dphoto/dphoto/backup"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"github.com/thomasduchatelle/dphoto/dphoto/daemon"
	log "github.com/sirupsen/logrus"
)

func init() {
	daemon.VolumeManager = new(backupAdapter)
}

type backupAdapter struct{}

// OnMountedVolume can be used to a daemon to automatically start a backup
func (b *backupAdapter) OnMountedVolume(volume backupmodel.VolumeToBackup) {
	withContext := log.WithFields(log.Fields{
		"VolumeId": volume.UniqueId,
		"Path":     volume.Path,
	})

	metadata, err := b.findVolumeMetadata(volume.UniqueId)
	if err != nil {
		withContext.WithError(err).Errorln("Failed to find volume metadata...")
		return
	}
	if metadata == nil {
		err := b.createNewVolume(volume)
		if err != nil {
			withContext.WithError(err).Errorln("CreateNewVolume failed")
		}

	} else if volume.Path != "" {
		withContext.WithField("Name", metadata.Name).Infoln("Disk plugged")

		if metadata.AutoBackup {
			_, err = backup.StartBackupRunner("root", volume, backup.Options{})
			if err != nil {
				withContext.WithError(err).Errorf("Backup failed to start")
			}
		}
	}
}

// OnUnMountedVolume is used by daemon to mark the
func (b *backupAdapter) OnUnMountedVolume(uniqueId string) {
	withContext := log.WithFields(log.Fields{
		"VolumeId": uniqueId,
	})

	withContext.Infoln("Disk unplugged")
}

// return nil when not found
func (b *backupAdapter) findVolumeMetadata(string) (*backupmodel.VolumeMetadata, error) {
	panic("Not implemented")
}
func (b *backupAdapter) createNewVolume(volume backupmodel.VolumeToBackup) error {
	panic("Not implemented")
}
