package backup

import (
	"duchatelle.io/dphoto/dphoto/scanner"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type filter struct {
	volumeId           string
	lastVolumeSnapshot map[string]uint
	currentSnapshot    []scanner.SimpleMediaSignature
}

func newMediaFilter(volume *scanner.VolumeToBackup) (*filter, error) {
	snapshot, err := VolumeRepository.RestoreLastSnapshot(volume.UniqueId)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to restore previous snapshot fo volume %s", volume.UniqueId)
	}

	f := &filter{
		volumeId:           volume.UniqueId,
		lastVolumeSnapshot: make(map[string]uint),
	}
	for _, m := range snapshot {
		f.lastVolumeSnapshot[m.RelativePath] = m.Size
	}

	return f, nil
}

func (f *filter) Filter(found scanner.FoundMedia) bool {
	f.currentSnapshot = append(f.currentSnapshot, *found.SimpleSignature())

	size, ok := f.lastVolumeSnapshot[found.SimpleSignature().RelativePath]
	keep := !ok || size != found.SimpleSignature().Size

	if !keep {
		log.Debugf("Filter > filter out media %s", found)
	}
	return keep
}

func (f *filter) StoreState(backupId string) error {
	return VolumeRepository.StoreSnapshot(f.volumeId, backupId, f.currentSnapshot)
}
