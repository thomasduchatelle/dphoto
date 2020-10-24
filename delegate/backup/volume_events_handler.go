package backup

import (
	log "github.com/sirupsen/logrus"
)

func OnMountedVolume(volume RemovableVolume) {
	withContext := log.WithFields(log.Fields{
		"VolumeId":   volume.UniqueId,
		"MountPaths": volume.MountPaths,
	})

	metadata, err := VolumeRepository.FindVolumeMetadata(volume.UniqueId)
	if err != nil {
		withContext.WithError(err).Errorln("Failed to find volume metadata...")
		return
	}
	if metadata == nil {
		err := VolumeRepository.CreateNewVolume(volume)
		if err != nil {
			withContext.WithError(err).Errorln("CreateNewVolume failed")
		}

	} else if len(volume.MountPaths) > 0 {
		withContext.WithField("Name", metadata.Name).Infoln("Disk plugged")

		if metadata.AutoBackup {
			_, err := StartBackupRunner(volume)
			if err != nil {
				withContext.WithError(err).Errorf("Backup failed to start")
			}
		}
	}
}

func OnUnMountedVolume(uniqueId string) {
	withContext := log.WithFields(log.Fields{
		"VolumeId": uniqueId,
	})

	withContext.Infoln("Disk unplugged")
}
