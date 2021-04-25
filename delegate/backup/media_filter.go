package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/pkg/errors"
)

type filter struct {
	volumeId           string
	lastVolumeSnapshot map[string]int
	currentSnapshot    []model.SimpleMediaSignature
}

func newMediaFilter(volume *model.VolumeToBackup) (*filter, error) {
	snapshot, err := VolumeRepository.RestoreLastSnapshot(volume.UniqueId)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to restore previous snapshot fo volume %s", volume.UniqueId)
	}

	f := &filter{
		volumeId:           volume.UniqueId,
		lastVolumeSnapshot: make(map[string]int),
	}
	for _, m := range snapshot {
		f.lastVolumeSnapshot[m.RelativePath] = m.Size
	}

	return f, nil
}

func (f *filter) Filter(found model.FoundMedia) bool {
	f.currentSnapshot = append(f.currentSnapshot, *found.SimpleSignature())

	size, ok := f.lastVolumeSnapshot[found.SimpleSignature().RelativePath]
	return !ok || size != found.SimpleSignature().Size
}

func (f *filter) StoreState(backupId string) error {
	return VolumeRepository.StoreSnapshot(f.volumeId, backupId, f.currentSnapshot)
}
