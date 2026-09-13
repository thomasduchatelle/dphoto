package archive_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
)

func TestRelocate(t *testing.T) {
	const owner = "ironman@avenger.marvel"

	type indexEntry struct {
		id       string
		location string
	}
	type relocated struct {
		id     string
		newKey string
	}
	tests := []struct {
		name              string
		ids               []string
		targetFolder      string
		seed              []indexEntry
		wantRelocated     []relocated
		wantOldGone       []string
		wantUnchangedKeep []string
	}{
		{
			name:         "it should relocate an image from both physical store and index",
			ids:          []string{"id-01"},
			targetFolder: "/newFolder",
			seed: []indexEntry{
				{"id-01", owner + "/deep/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg"},
			},
			wantRelocated: []relocated{
				{"id-01", owner + "/newFolder/2022-06-19_15-02-10_16c6dfa0.jpg"},
			},
			wantOldGone: []string{owner + "/deep/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg"},
		},
		{
			name:         "it should not do anything if the image belongs to someone else",
			ids:          []string{"id-01"},
			targetFolder: "/newFolder",
			seed: []indexEntry{
				{"id-01", "captainamerica@avenger.marvel/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg"},
			},
			wantUnchangedKeep: []string{"captainamerica@avenger.marvel/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg"},
		},
		{
			name:         "it should ignore extra responses from GetLocation and ignore (log) unknown media ids",
			ids:          []string{"id-01", "id-02"},
			targetFolder: "/newFolder",
			seed: []indexEntry{
				{"id-01", owner + "/01.jpg"},
			},
			wantRelocated: []relocated{
				{"id-01", owner + "/newFolder/01.jpg"},
			},
			wantOldGone: []string{owner + "/01.jpg"},
		},
		{
			name:         "it should clean the location from any suffix",
			ids:          []string{"id-01"},
			targetFolder: "/newFolder",
			seed: []indexEntry{
				{"id-01", owner + "/oldFolder1/2022-06-19_15-02-10_16c6dfa0_something_might_have_had_been_added_to_make_it_unique.jpg"},
			},
			wantRelocated: []relocated{
				{"id-01", owner + "/newFolder/2022-06-19_15-02-10_16c6dfa0.jpg"},
			},
			wantOldGone: []string{owner + "/oldFolder1/2022-06-19_15-02-10_16c6dfa0_something_might_have_had_been_added_to_make_it_unique.jpg"},
		},
		{
			name:         "it should support files now following a proper format",
			ids:          []string{"id-01"},
			targetFolder: "/newFolder",
			seed: []indexEntry{
				{"id-01", owner + "//this/is/a_really-strange^format"},
			},
			wantRelocated: []relocated{
				{"id-01", owner + "/newFolder/a_really-strange^format"},
			},
			wantOldGone: []string{owner + "//this/is/a_really-strange^format"},
		},
		{
			name:         "it should batch finding, indexing, and s3 deletion operations",
			ids:          []string{"id-01", "id-02", "id-03"},
			targetFolder: "/newFolder",
			seed: []indexEntry{
				{"id-01", owner + "/01.jpg"},
				{"id-02", owner + "/02.jpg"},
				{"id-03", owner + "/03.jpg"},
			},
			wantRelocated: []relocated{
				{"id-01", owner + "/newFolder/01.jpg"},
				{"id-02", owner + "/newFolder/02.jpg"},
				{"id-03", owner + "/newFolder/03.jpg"},
			},
			wantOldGone: []string{owner + "/01.jpg", owner + "/02.jpg", owner + "/03.jpg"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := NewARepositoryInMemory()
			store := NewStoreInMemory()
			for _, entry := range tt.seed {
				_ = repository.AddLocation(owner, entry.id, entry.location)
				store.Content[entry.location] = []byte("content-" + entry.id)
			}
			archive.Init(repository, store, NewCacheInMemory(), NewAsyncJobInMemory())

			err := archive.Relocate(owner, tt.ids, tt.targetFolder)

			if assert.NoError(t, err) {
				for _, r := range tt.wantRelocated {
					got, findErr := repository.FindById(owner, r.id)
					if assert.NoError(t, findErr) {
						assert.Equal(t, r.newKey, got, "index should point to new key for %s", r.id)
					}
					assert.True(t, store.Has(r.newKey), "content should exist at new key %s", r.newKey)
				}
				for _, oldKey := range tt.wantOldGone {
					assert.False(t, store.Has(oldKey), "content at old key %s should have been deleted", oldKey)
				}
				for _, key := range tt.wantUnchangedKeep {
					assert.True(t, store.Has(key), "content at %s should remain untouched", key)
				}
			}
		})
	}

	// E1 — data-loss safety: no deletion when index update fails.
	// Uses testify/mock inline because a Fake cannot inject an UpdateLocations failure.
	t.Run("it should not delete anything if the index cannot be updated", func(t *testing.T) {
		const previousLocation = owner + "/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg"

		backing := NewARepositoryInMemory()
		_ = backing.AddLocation(owner, "id-01", previousLocation)
		repository := &failingUpdateRepository{ARepositoryInMemory: backing}
		repository.On("UpdateLocations", owner, mock.Anything).Return(errors.Errorf("TEST - should abort deletion"))

		store := NewStoreInMemory()
		store.Content[previousLocation] = []byte("content")

		archive.Init(repository, store, NewCacheInMemory(), NewAsyncJobInMemory())

		err := archive.Relocate(owner, []string{"id-01"}, "/newFolder")

		assert.Error(t, err)
		assert.True(t, store.Has(previousLocation), "old content must NOT be deleted when index update fails")
	})
}

// failingUpdateRepository embeds a real Fake but overrides UpdateLocations with
// a testify/mock stub so we can inject an error mid-flow.
type failingUpdateRepository struct {
	*ARepositoryInMemory
	mock.Mock
}

func (f *failingUpdateRepository) UpdateLocations(owner string, locations map[string]string) error {
	args := f.Called(owner, locations)
	return args.Error(0)
}
