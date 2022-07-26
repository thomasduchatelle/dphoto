package archive_test

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"github.com/thomasduchatelle/dphoto/mocks"
	"testing"
)

func TestRelocate(t *testing.T) {
	const owner = "ironman@avenger.marvel"

	tests := []struct {
		name         string
		ids          []string
		targetFolder string
		spec         func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter)
		wantErr      bool
	}{
		{
			name:         "it should relocate an image from both physical store and index",
			ids:          []string{"id-01"},
			targetFolder: "/newFolder",
			spec: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter) {
				const previousLocation = owner + "/deep/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg"

				repository.On("FindByIds", owner, []string{"id-01"}).Once().Return(map[string]string{
					"id-01": previousLocation,
				}, nil)

				store.On("Copy", previousLocation, archive.DestructuredKey{
					Prefix: owner + "/newFolder/2022-06-19_15-02-10_16c6dfa0",
					Suffix: ".jpg",
				}).Once().Return("newkey-01", nil)

				repository.On("UpdateLocations", owner, map[string]string{
					"id-01": "newkey-01",
				}).Once().Return(nil)

				store.On("Delete", []string{previousLocation}).Once().Return(nil)
			},
		},
		{
			name:         "it should not do anything if the image belongs to someone else",
			ids:          []string{"id-01"},
			targetFolder: "/newFolder",
			spec: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter) {
				const previousLocation = "captainamerica@avenger.marvel/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg"

				repository.On("FindByIds", owner, []string{"id-01"}).Once().Return(map[string]string{
					"id-01": previousLocation,
				}, nil)
			},
		},
		{
			name:         "it should ignore extra responses from GetLocation and ignore (log) unknown media ids",
			ids:          []string{"id-01", "id-02"},
			targetFolder: "/newFolder",
			spec: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter) {
				const previousLocation = owner + "/01.jpg"
				repository.On("FindByIds", owner, []string{"id-01", "id-02"}).Once().Return(map[string]string{
					"id-01": previousLocation,
					"id-03": "03.jpg",
				}, nil)

				store.On("Copy", previousLocation, archive.DestructuredKey{
					Prefix: owner + "/newFolder/01",
					Suffix: ".jpg",
				}).Once().Return("newkey-01", nil)

				repository.On("UpdateLocations", owner, map[string]string{
					"id-01": "newkey-01",
				}).Once().Return(nil)

				store.On("Delete", []string{previousLocation}).Once().Return(nil)
			},
		},
		{
			name:         "it should clean the location from any suffix",
			ids:          []string{"id-01"},
			targetFolder: "/newFolder",
			spec: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter) {
				const previousLocation = owner + "/oldFolder1/2022-06-19_15-02-10_16c6dfa0_something_might_have_had_been_added_to_make_it_unique.jpg"
				repository.On("FindByIds", owner, []string{"id-01"}).Once().Return(map[string]string{
					"id-01": previousLocation,
				}, nil)

				store.On("Copy", previousLocation, archive.DestructuredKey{
					Prefix: owner + "/newFolder/2022-06-19_15-02-10_16c6dfa0",
					Suffix: ".jpg",
				}).Once().Return("newkey-01", nil)

				repository.On("UpdateLocations", owner, map[string]string{
					"id-01": "newkey-01",
				}).Once().Return(nil)

				store.On("Delete", []string{previousLocation}).Once().Return(nil)
			},
		},
		{
			name:         "it should support files now following a proper format",
			ids:          []string{"id-01"},
			targetFolder: "/newFolder",
			spec: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter) {
				const previousLocation = owner + "//this/is/a_really-strange^format"
				repository.On("FindByIds", owner, []string{"id-01"}).Once().Return(map[string]string{
					"id-01": previousLocation,
				}, nil)

				store.On("Copy", previousLocation, archive.DestructuredKey{
					Prefix: owner + "/newFolder/a_really-strange^format",
					Suffix: "",
				}).Once().Return("newkey-01", nil)

				repository.On("UpdateLocations", owner, map[string]string{
					"id-01": "newkey-01",
				}).Once().Return(nil)

				store.On("Delete", []string{previousLocation}).Once().Return(nil)
			},
		},
		{
			name:         "it should batch finding, indexing, and s3 deletion operations",
			ids:          []string{"id-01", "id-02", "id-03"},
			targetFolder: "/newFolder",
			spec: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter) {
				previousLocations := []string{owner + "/01.jpg", owner + "/02.jpg", owner + "/03.jpg"}
				repository.On("FindByIds", owner, []string{"id-01", "id-02", "id-03"}).Once().Return(map[string]string{
					"id-01": previousLocations[0],
					"id-02": previousLocations[1],
					"id-03": previousLocations[2],
				}, nil)

				store.On("Copy", previousLocations[0], archive.DestructuredKey{
					Prefix: owner + "/newFolder/01",
					Suffix: ".jpg",
				}).Once().Return("newkey-01", nil)
				store.On("Copy", previousLocations[1], archive.DestructuredKey{
					Prefix: owner + "/newFolder/02",
					Suffix: ".jpg",
				}).Once().Return("newkey-02", nil)
				store.On("Copy", previousLocations[2], archive.DestructuredKey{
					Prefix: owner + "/newFolder/03",
					Suffix: ".jpg",
				}).Once().Return("newkey-03", nil)

				repository.On("UpdateLocations", owner, map[string]string{
					"id-01": "newkey-01",
					"id-02": "newkey-02",
					"id-03": "newkey-03",
				}).Once().Return(nil)

				store.On("Delete", previousLocations).Once().Return(nil)
			},
		},
		{
			name:         "it should not delete anything if the index cannot be updated",
			ids:          []string{"id-01"},
			targetFolder: "/newFolder",
			spec: func(repository *mocks.ARepositoryAdapter, store *mocks.StoreAdapter) {
				const previousLocation = owner + "/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg"

				repository.On("FindByIds", owner, []string{"id-01"}).Once().Return(map[string]string{
					"id-01": previousLocation,
				}, nil)

				store.On("Copy", previousLocation, archive.DestructuredKey{
					Prefix: owner + "/newFolder/2022-06-19_15-02-10_16c6dfa0",
					Suffix: ".jpg",
				}).Once().Return("newkey-01", nil)

				repository.On("UpdateLocations", owner, map[string]string{
					"id-01": "newkey-01",
				}).Once().Return(errors.Errorf("TEST - should abort deletion"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryAdapter := mocks.NewARepositoryAdapter(t)
			storeAdapter := mocks.NewStoreAdapter(t)

			tt.spec(repositoryAdapter, storeAdapter)

			archive.Init(repositoryAdapter, storeAdapter, mocks.NewCacheAdapter(t), mocks.NewAsyncJobAdapter(t))

			err := archive.Relocate(owner, tt.ids, tt.targetFolder)

			if tt.wantErr {
				assert.Errorf(t, err, tt.name)
			} else {
				assert.NoError(t, err, tt.name)
			}
		})
	}
}
