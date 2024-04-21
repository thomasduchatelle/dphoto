package analysiscache_test

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysiscache"
	"testing"
	"time"
)

func TestDecoratorInstance_Analyse(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	assert.NoError(t, err)
	defer db.Close()

	const computedMediaHash = "qwertyuiop"
	const cachedMediaHash = "cached-sha256-images"
	sometime := time.Date(2024, 3, 9, 23, 10, 11, 0, time.UTC)
	foundMedia := backup.NewInMemoryMedia("/avengers/ironman/stark-tower-01.png", sometime, []byte("some content"))
	analysedWasCached := []backup.ProgressEvent{{Type: "analysed-from-cache", Count: 1, Size: 12, Album: "", MediaType: ""}}

	const badgerKey = "/ram/avengers/ironman/stark-tower-01.png##12"
	const badgerPayload = `{
  "lastModification": "2024-03-09T23:10:11.000Z",
  "type": "IMAGE",
  "sha256Hash": "cached-sha256-images",
  "details": {
    "width": 120,
    "height": 42
  }
}`
	var recordHasBeenStored = map[string]analysiscache.Payload{
		badgerKey: {
			LastModification: sometime,
			Type:             "IMAGE",
			Sha256Hash:       computedMediaHash,
			Details: backup.MediaDetails{
				Width:  120,
				Height: 42,
			},
		},
	}
	var recordHasBeenKept = map[string]analysiscache.Payload{
		badgerKey: {
			LastModification: sometime,
			Type:             "IMAGE",
			Sha256Hash:       cachedMediaHash,
			Details: backup.MediaDetails{
				Width:  120,
				Height: 42,
			},
		},
	}

	type fields struct {
		Delegate backup.RunnerAnalyser
	}
	type args struct {
		found backup.FoundMedia
	}

	doesNotHaveARecordInCache := func(t *testing.T, db *badger.DB) error {
		return nil
	}
	hasARecordInCache := func(t *testing.T, db *badger.DB) error {
		return db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(badgerKey), []byte(badgerPayload))
		})
	}
	doesNotCallTheDelegate := func(t *testing.T) fields {
		return fields{
			Delegate: backup.RunnerAnalyserFunc(func(found backup.FoundMedia, progressChannel chan *backup.ProgressEvent) (*backup.AnalysedMedia, error) {
				assert.Fail(t, "Unexpected call on Analyse", "unexpected Analyse(%+v, ...))", found)
				return nil, nil
			}),
		}
	}
	doesCallTheDelegate := func(t *testing.T) fields {
		return fields{
			Delegate: backup.RunnerAnalyserFunc(func(found backup.FoundMedia, progressChannel chan *backup.ProgressEvent) (*backup.AnalysedMedia, error) {
				return &backup.AnalysedMedia{
					FoundMedia: found,
					Type:       backup.MediaTypeImage,
					Sha256Hash: computedMediaHash,
					Details: &backup.MediaDetails{
						Width:  120,
						Height: 42,
					},
				}, nil
			}),
		}

	}

	tests := []struct {
		name       string
		init       func(t *testing.T, db *badger.DB) error
		mocks      func(t *testing.T) fields
		args       args
		want       *backup.AnalysedMedia
		wantDB     map[string]analysiscache.Payload
		wantEvents []backup.ProgressEvent
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:  "it should call the delegate and cache the result when no cache",
			init:  doesNotHaveARecordInCache,
			mocks: doesCallTheDelegate,
			args: args{
				found: foundMedia,
			},
			want: &backup.AnalysedMedia{
				FoundMedia: foundMedia,
				Type:       backup.MediaTypeImage,
				Sha256Hash: computedMediaHash,
				Details: &backup.MediaDetails{
					Width:  120,
					Height: 42,
				},
			},
			wantDB:  recordHasBeenStored,
			wantErr: assert.NoError,
		},
		{
			name:  "it should use the cache when key and last modified date match",
			init:  hasARecordInCache,
			mocks: doesNotCallTheDelegate,
			args: args{
				found: foundMedia,
			},
			want: &backup.AnalysedMedia{
				FoundMedia: foundMedia,
				Type:       backup.MediaTypeImage,
				Sha256Hash: cachedMediaHash,
				Details: &backup.MediaDetails{
					Width:  120,
					Height: 42,
				},
			},
			wantDB:     recordHasBeenKept,
			wantEvents: analysedWasCached,
			wantErr:    assert.NoError,
		},
		{
			name:  "it should call the delegate and override the result when last modification doesn't match",
			init:  hasARecordInCache,
			mocks: doesCallTheDelegate,
			args: args{
				found: backup.NewInMemoryMedia("/avengers/ironman/stark-tower-01.png", sometime.Add(1*time.Minute), []byte("some content")),
			},
			want: &backup.AnalysedMedia{
				FoundMedia: backup.NewInMemoryMedia("/avengers/ironman/stark-tower-01.png", sometime.Add(1*time.Minute), []byte("some content")),
				Type:       backup.MediaTypeImage,
				Sha256Hash: computedMediaHash,
				Details: &backup.MediaDetails{
					Width:  120,
					Height: 42,
				},
			},
			wantDB: map[string]analysiscache.Payload{
				badgerKey: {
					LastModification: sometime.Add(1 * time.Minute),
					Type:             "IMAGE",
					Sha256Hash:       computedMediaHash,
					Details: backup.MediaDetails{
						Width:  120,
						Height: 42,
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:  "it should use the cache and ignore last modification if ZERO is requested (mean not supported)",
			init:  hasARecordInCache,
			mocks: doesNotCallTheDelegate,
			args: args{
				found: backup.NewInMemoryMedia("/avengers/ironman/stark-tower-01.png", time.Time{}, []byte("some content")),
			},
			want: &backup.AnalysedMedia{
				FoundMedia: backup.NewInMemoryMedia("/avengers/ironman/stark-tower-01.png", time.Time{}, []byte("some content")),
				Type:       backup.MediaTypeImage,
				Sha256Hash: cachedMediaHash,
				Details: &backup.MediaDetails{
					Width:  120,
					Height: 42,
				},
			},
			wantDB:     recordHasBeenKept,
			wantEvents: analysedWasCached,
			wantErr:    assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.init(t, db)
			if !assert.NoError(t, err, "tt.init()") {
				assert.FailNow(t, "tt.init()")
			}

			mockedFields := tt.mocks(t)

			d := &analysiscache.DecoratorInstance{
				DB:       db,
				Delegate: mockedFields.Delegate,
			}
			progressChannel, completion := NewSinkChannel()

			got, err := d.Analyse(tt.args.found, progressChannel)
			close(progressChannel)

			if !tt.wantErr(t, err, fmt.Sprintf("Analyse(%v)", tt.args.found)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Analyse(%v)", tt.args.found)

			gotDB, err := databaseDump(db)
			if assert.NoError(t, err, "databaseDump(db)") {
				assert.Equal(t, tt.wantDB, gotDB)
			}

			gotEvents := <-completion
			assert.Equalf(t, tt.wantEvents, gotEvents, "events: %+v", gotEvents)
		})
	}
}

func databaseDump(db *badger.DB) (map[string]analysiscache.Payload, error) {
	dump := make(map[string]analysiscache.Payload)

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		var (
			buffer []byte
			err    error
		)

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			buffer, err = item.ValueCopy(buffer)
			if err != nil {
				return err
			}

			var content analysiscache.Payload
			err = json.Unmarshal(buffer, &content)
			if err != nil {
				return err
			}

			dump[string(item.Key())] = content
		}

		return nil
	})

	return dump, err
}

func NewSinkChannel() (chan *backup.ProgressEvent, chan []backup.ProgressEvent) {
	channel := make(chan *backup.ProgressEvent)
	completion := make(chan []backup.ProgressEvent, 1)
	holder := struct {
		events []backup.ProgressEvent
	}{}

	go func() {
		for event := range channel {
			holder.events = append(holder.events, *event)
		}

		completion <- holder.events
	}()

	return channel, completion
}
