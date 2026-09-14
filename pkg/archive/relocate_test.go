package archive_test

import (
	"sort"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
)

func TestRelocate(t *testing.T) {
	const relocateOwner = "ironman@avenger.marvel"

	type fields struct {
		Repository *ARepositoryInMemory
		Store      *StoreInMemory
	}
	type args struct {
		ids          []string
		targetFolder string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         assert.ErrorAssertionFunc
		expectLocations map[string]string
		expectStoreKeys []string
	}{
		{
			name: "it should relocate an image from both physical store and index",
			fields: fields{
				Repository: repositoryWithLocations(relocateOwner, map[string]string{
					"id-01": relocateOwner + "/deep/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg",
				}),
				Store: storeWithKeys(map[string][]byte{
					relocateOwner + "/deep/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg": []byte("content-01"),
				}),
			},
			args:    args{ids: []string{"id-01"}, targetFolder: "/newFolder"},
			wantErr: assert.NoError,
			expectLocations: map[string]string{
				"id-01": relocateOwner + "/newFolder/2022-06-19_15-02-10_16c6dfa0.jpg",
			},
			expectStoreKeys: []string{relocateOwner + "/newFolder/2022-06-19_15-02-10_16c6dfa0.jpg"},
		},
		{
			name: "it should not do anything if the image belongs to someone else",
			fields: fields{
				Repository: repositoryWithLocations(relocateOwner, map[string]string{
					"id-01": "captainamerica@avenger.marvel/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg",
				}),
				Store: storeWithKeys(map[string][]byte{
					"captainamerica@avenger.marvel/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg": []byte("content-01"),
				}),
			},
			args:    args{ids: []string{"id-01"}, targetFolder: "/newFolder"},
			wantErr: assert.NoError,
			expectLocations: map[string]string{
				"id-01": "captainamerica@avenger.marvel/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg",
			},
			expectStoreKeys: []string{"captainamerica@avenger.marvel/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg"},
		},
		{
			name: "it should ignore unknown media ids and extra locations returned by the repository",
			fields: fields{
				Repository: repositoryWithLocations(relocateOwner, map[string]string{
					"id-01": relocateOwner + "/01.jpg",
					"id-03": "03.jpg",
				}),
				Store: storeWithKeys(map[string][]byte{
					relocateOwner + "/01.jpg": []byte("content-01"),
				}),
			},
			args:    args{ids: []string{"id-01", "id-02"}, targetFolder: "/newFolder"},
			wantErr: assert.NoError,
			expectLocations: map[string]string{
				"id-01": relocateOwner + "/newFolder/01.jpg",
				"id-03": "03.jpg",
			},
			expectStoreKeys: []string{relocateOwner + "/newFolder/01.jpg"},
		},
		{
			name: "it should clean the location from any extra suffix appended to make the filename unique",
			fields: fields{
				Repository: repositoryWithLocations(relocateOwner, map[string]string{
					"id-01": relocateOwner + "/oldFolder1/2022-06-19_15-02-10_16c6dfa0_something_might_have_had_been_added_to_make_it_unique.jpg",
				}),
				Store: storeWithKeys(map[string][]byte{
					relocateOwner + "/oldFolder1/2022-06-19_15-02-10_16c6dfa0_something_might_have_had_been_added_to_make_it_unique.jpg": []byte("content-01"),
				}),
			},
			args:    args{ids: []string{"id-01"}, targetFolder: "/newFolder"},
			wantErr: assert.NoError,
			expectLocations: map[string]string{
				"id-01": relocateOwner + "/newFolder/2022-06-19_15-02-10_16c6dfa0.jpg",
			},
			expectStoreKeys: []string{relocateOwner + "/newFolder/2022-06-19_15-02-10_16c6dfa0.jpg"},
		},
		{
			name: "it should support files not following the standard filename format",
			fields: fields{
				Repository: repositoryWithLocations(relocateOwner, map[string]string{
					"id-01": relocateOwner + "//this/is/a_really-strange^format",
				}),
				Store: storeWithKeys(map[string][]byte{
					relocateOwner + "//this/is/a_really-strange^format": []byte("content-01"),
				}),
			},
			args:    args{ids: []string{"id-01"}, targetFolder: "/newFolder"},
			wantErr: assert.NoError,
			expectLocations: map[string]string{
				"id-01": relocateOwner + "/newFolder/a_really-strange^format",
			},
			expectStoreKeys: []string{relocateOwner + "/newFolder/a_really-strange^format"},
		},
		{
			name: "it should batch finding, indexing, and store deletion operations",
			fields: fields{
				Repository: repositoryWithLocations(relocateOwner, map[string]string{
					"id-01": relocateOwner + "/01.jpg",
					"id-02": relocateOwner + "/02.jpg",
					"id-03": relocateOwner + "/03.jpg",
				}),
				Store: storeWithKeys(map[string][]byte{
					relocateOwner + "/01.jpg": []byte("content-01"),
					relocateOwner + "/02.jpg": []byte("content-02"),
					relocateOwner + "/03.jpg": []byte("content-03"),
				}),
			},
			args:    args{ids: []string{"id-01", "id-02", "id-03"}, targetFolder: "/newFolder"},
			wantErr: assert.NoError,
			expectLocations: map[string]string{
				"id-01": relocateOwner + "/newFolder/01.jpg",
				"id-02": relocateOwner + "/newFolder/02.jpg",
				"id-03": relocateOwner + "/newFolder/03.jpg",
			},
			expectStoreKeys: []string{
				relocateOwner + "/newFolder/01.jpg",
				relocateOwner + "/newFolder/02.jpg",
				relocateOwner + "/newFolder/03.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive.Init(tt.fields.Repository, tt.fields.Store, NewCacheInMemory(), NewAsyncJobInMemory())

			err := archive.Relocate(relocateOwner, tt.args.ids, tt.args.targetFolder)

			if !tt.wantErr(t, err) {
				return
			}

			for id, expectedKey := range tt.expectLocations {
				got, findErr := tt.fields.Repository.FindById(relocateOwner, id)
				if !assert.NoErrorf(t, findErr, "expected id %s to be indexed", id) {
					continue
				}
				assert.Equalf(t, expectedKey, got, "unexpected location for id %s", id)
			}

			gotKeys := make([]string, 0, len(tt.fields.Store.Content))
			for key := range tt.fields.Store.Content {
				gotKeys = append(gotKeys, key)
			}
			sort.Strings(gotKeys)
			expectedKeys := append([]string(nil), tt.expectStoreKeys...)
			sort.Strings(expectedKeys)
			assert.Equal(t, expectedKeys, gotKeys)
		})
	}
}

// TestRelocate_shouldNotDeleteIfIndexUpdateFails is a data-loss safety test (E1). It uses
// testify/mock inline to inject an UpdateLocations failure on the repository — the Fake would
// otherwise succeed unconditionally and we could not verify the "no delete after failed index
// update" invariant.
func TestRelocate_shouldNotDeleteIfIndexUpdateFails(t *testing.T) {
	const relocateOwner = "ironman@avenger.marvel"
	const previousLocation = relocateOwner + "/oldFolder1/2022-06-19_15-02-10_16c6dfa0.jpg"

	repository := &failingUpdateRepository{
		ARepositoryInMemory: *repositoryWithLocations(relocateOwner, map[string]string{
			"id-01": previousLocation,
		}),
	}
	repository.On("UpdateLocations", mock.Anything, mock.Anything).Return(errors.Errorf("TEST - should abort deletion"))

	store := storeWithKeys(map[string][]byte{previousLocation: []byte("content-01")})

	archive.Init(repository, store, NewCacheInMemory(), NewAsyncJobInMemory())

	err := archive.Relocate(relocateOwner, []string{"id-01"}, "/newFolder")

	assert.Error(t, err)
	assert.True(t, store.Has(previousLocation), "old key must still be present when index update fails")
}

func repositoryWithLocations(owner string, locations map[string]string) *ARepositoryInMemory {
	repository := NewARepositoryInMemory()
	for id, key := range locations {
		_ = repository.AddLocation(owner, id, key)
	}
	return repository
}

func storeWithKeys(content map[string][]byte) *StoreInMemory {
	store := NewStoreInMemory()
	for key, data := range content {
		store.Content[key] = data
	}
	return store
}

type failingUpdateRepository struct {
	ARepositoryInMemory
	mock.Mock
}

func (f *failingUpdateRepository) UpdateLocations(owner string, locations map[string]string) error {
	args := f.Called(owner, locations)
	return args.Error(0)
}
