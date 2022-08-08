// Package arepositoryjson stores in JSON indexes of all medias stored locally
package arepositoryjson

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func New(archiveLocalIndex string) (archive.ARepositoryAdapter, error) {
	patchesDir := path.Join(archiveLocalIndex, "patches")
	repo := &repository{
		patchesDir: patchesDir,
		indexSync:  sync.Mutex{},
		index:      make(map[string]map[string]string),
	}

	err := os.MkdirAll(patchesDir, 0744)
	if err != nil {
		return nil, errors.Wrapf(err, "creating dir for patches")
	}
	err = filepath.WalkDir(patchesDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".json") {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrapf(err, "reading file %s", path)
			}

			patch := make(map[string]map[string]string)
			err = json.Unmarshal(content, &patch)
			if err != nil {
				return errors.Wrapf(err, "unmarshaling JSON %s from file %s", string(content), path)
			}

			repo.applyPatch(patch)
		}

		return nil
	})

	return repo, errors.Wrapf(err, "loading previous updates")
}

type repository struct {
	patchesDir    string
	indexSync     sync.Mutex
	index         map[string]map[string]string
	updateCounter int
}

func (r *repository) FindById(owner, id string) (string, error) {
	location, ok := r.indexForOwner(owner)[id]
	if !ok {
		return "", archive.NotFoundError
	}

	return location, nil
}

func (r *repository) FindByIds(owner string, ids []string) (map[string]string, error) {
	subset := make(map[string]string)
	for _, id := range ids {
		if location, ok := r.indexForOwner(owner)[id]; ok {
			subset[id] = location
		}
	}

	return subset, nil
}

func (r *repository) AddLocation(owner, id, key string) error {
	return r.UpdateLocations(owner, map[string]string{
		id: key,
	})
}

func (r *repository) UpdateLocations(owner string, locations map[string]string) error {
	r.indexSync.Lock()
	defer r.indexSync.Unlock()

	patch := map[string]map[string]string{
		owner: locations,
	}
	err := r.persistPath(patch)
	if err != nil {
		return err
	}

	r.applyPatch(patch)

	return nil
}

func (r *repository) FindIdsFromKeyPrefix(keyPrefix string) (map[string]string, error) {
	panic("FindIdsFromKeyPrefix is not supported locally (limited by multi-tenancy support)")
}

func (r *repository) indexForOwner(owner string) map[string]string {
	r.indexSync.Lock()
	defer r.indexSync.Unlock()

	ownerIndex, ok := r.index[owner]
	if !ok {
		ownerIndex = make(map[string]string)
		r.index[owner] = ownerIndex
	}

	return ownerIndex
}

func (r *repository) persistPath(patch map[string]map[string]string) error {
	jsonUpdate, err := json.Marshal(patch)
	if err != nil {
		return errors.Wrapf(err, "marshaling %+v locations", patch)
	}

	err = ioutil.WriteFile(path.Join(r.patchesDir, fmt.Sprintf("%s-%04d.json", time.Now().Format("2006-01-02T15-04-05"), r.updateCounter)), jsonUpdate, 0644)
	r.updateCounter++
	return errors.Wrapf(err, "writing archive index with content %s", string(jsonUpdate))
}

func (r *repository) applyPatch(patch map[string]map[string]string) {
	for owner, locations := range patch {
		for key, location := range locations {
			r.indexForOwner(owner)[key] = location
		}
	}
}
