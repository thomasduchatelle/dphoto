package archive

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path"
	"regexp"
	"strings"
)

// Relocate physically moves media to the target folder
func Relocate(owner string, ids []string, targetFolder string) error {
	locations, err := repositoryPort.FindByIds(owner, ids)
	if err != nil {
		return errors.Wrapf(err, "cannot relocate medias")
	}

	newLocations := make(map[string]string)
	var oldLocations []string
	for _, id := range ids {
		if previousLocation, ok := locations[id]; ok && strings.HasPrefix(previousLocation, owner+"/") {
			baseFilename := strings.TrimSuffix(path.Base(previousLocation), path.Ext(previousLocation))

			matches := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}_\\d{2}-\\d{2}-\\d{2}_\\w{6,8}").Find([]byte(path.Base(previousLocation)))
			if len(matches) > 0 {
				baseFilename = string(matches)
			}

			destinationKey := DestructuredKey{
				Prefix: path.Join(owner, targetFolder, baseFilename),
				Suffix: path.Ext(previousLocation),
			}
			newLocations[id], err = storePort.Copy(locations[id], destinationKey)
			if err != nil {
				return errors.Wrapf(err, "failed to physically move file %s from %s to %s", id, previousLocation, destinationKey)
			}

			oldLocations = append(oldLocations, previousLocation)

		} else {
			log.WithFields(log.Fields{
				"Owner": owner,
				"Id":    id,
			}).Warnf("file location for media %s does not exist", id)
		}
	}

	if len(newLocations) == 0 || len(oldLocations) == 0 {
		return nil
	}

	err = repositoryPort.UpdateLocations(owner, newLocations)
	if err != nil {
		return errors.Wrapf(err, "failed to update media locations")
	}

	err = storePort.Delete(oldLocations)
	return errors.Wrapf(err, "failed to remove moved files")
}
