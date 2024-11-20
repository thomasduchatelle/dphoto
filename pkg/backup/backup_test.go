package backup

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"maps"
	"os"
	"path"
	"slices"
	"sync"
	"testing"
	"time"
)

func TestBackupAcceptance(t *testing.T) {
	const owner = ownermodel.Owner("ironman")

	rejectFolder, err := os.MkdirTemp(os.TempDir(), "dphoto-unit-testbackupacceptance")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(rejectFolder)

	analysedMedias := []*AnalysedMedia{
		{
			FoundMedia: NewInMemoryMedia("folder1/file_1.jpg", time.Now(), []byte("2022-06-18")),
			Type:       MediaTypeImage,
			Sha256Hash: "3e7574e8b640104d97597b200fd516c589f34be540e0a81a272fd488d12acaec",
			Details:    &MediaDetails{DateTime: time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC)},
		},
		{
			FoundMedia: NewInMemoryMedia("folder1/file_2.jpg", time.Now(), []byte("2022-06-19AB")),
			Type:       MediaTypeImage,
			Sha256Hash: "28f046d0ebae98f45512f98d581e7cdded28dd9cf50e7712615970dc15221cb3",
			Details:    &MediaDetails{DateTime: time.Date(2022, 6, 19, 0, 0, 0, 0, time.UTC)},
		},
		{
			FoundMedia: NewInMemoryMedia("folder1/file_3.mp4", time.Now(), []byte("2022-06-19ABC")),
			Type:       MediaTypeVideo,
			Sha256Hash: "e06a1b537d665585efa76f51a778bbb75473c984aa3ea3fbbcb8837db467d176",
			Details:    &MediaDetails{DateTime: time.Date(2022, 6, 19, 0, 0, 0, 0, time.UTC)},
		},
		{
			FoundMedia: NewInMemoryMedia("folder1/file_4.txt", time.Now(), []byte("2022-06-19ABCD")),
			Type:       MediaTypeOther,
			Sha256Hash: "ac94013eeb4c74b6f8c1ef2cea3bec85161732fe6161ba5d7e485efb28cb9e9e",
			Details:    &MediaDetails{DateTime: time.Date(2022, 6, 19, 0, 0, 0, 0, time.UTC)},
		},
	}
	doesNotExistReference1 := &CatalogReferenceStub{MediaIdValue: "media-id-1", AlbumFolderNameValue: "/album1"}
	doesExistReference1 := &CatalogReferenceStub{MediaIdValue: "media-id-1", AlbumFolderNameValue: "/album1", ExistsValue: true}
	doesNotExistReference2 := &CatalogReferenceStub{MediaIdValue: "media-id-2", AlbumFolderNameValue: "/album1"}
	doesNotExistReference3 := &CatalogReferenceStub{MediaIdValue: "media-id-3", AlbumFolderNameValue: "/album1"}
	doesNotExistReference4 := &CatalogReferenceStub{MediaIdValue: "media-id-4", AlbumFolderNameValue: "/album2", AlbumCreatedValue: true}

	readerFailingToParseFile := new(DetailsReaderFake)

	type fields struct {
		archive           ArchiveMediaPort
		cataloguerFactory CataloguerFactory
		insertMedia       InsertMediaPort
		detailsReaders    DetailsReader
	}
	type args struct {
		owner        ownermodel.Owner
		volume       SourceVolume
		optionsSlice []Options
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		want                    Report
		wantEvents              map[trackEvent]eventSummary // wantEvents won't be checked if nil
		wantErr                 assert.ErrorAssertionFunc
		wantRejectFolderContent []string
	}{
		{
			name: "it should upload a media going through the happy path",
			fields: fields{
				detailsReaders: new(DetailsReaderAdapterStub),
				archive: &AssertArchiveFake{
					ArchiveMediaPortFake: newArchiveMediaPortFake(),
					Want: map[ownermodel.Owner][]*BackingUpMediaRequest{
						owner: {
							{
								AnalysedMedia:    analysedMedias[0],
								CatalogReference: doesNotExistReference1,
							},
						},
					},
				},
				cataloguerFactory: &ReferencerFactoryFake{
					Cataloguer: &CatalogReferencerFake{
						analysedMedias[0]: doesNotExistReference1,
					},
				},
				insertMedia: &AssertInsertMediaPort{
					Want: []InsertMediaPortFakeEntry{
						{
							owner:           owner,
							ArchiveFilename: fakeArchiveFileName(analysedMedias[0]),
						},
					},
				},
			},
			args: args{
				owner: owner,
				volume: &InMemorySourceVolume{
					analysedMedias[0].FoundMedia,
				},
			},
			want: &backupReportBuilder{
				skipped: MediaCounterZero,
				countPerAlbum: map[string]*AlbumReport{
					doesNotExistReference1.AlbumFolderNameValue: countOfMedias(analysedMedias[0]),
				},
			},
			wantEvents: map[trackEvent]eventSummary{
				trackScanComplete: {SumCount: 1, SumSize: 10},
				trackCatalogued:   {SumCount: 1, SumSize: 10},
				trackUploaded:     {SumCount: 1, SumSize: 10, Albums: []string{"/album1"}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should upload a video and an unidentified medias, on two different albums, and get it represented on the report",
			fields: fields{
				detailsReaders: new(DetailsReaderAdapterStub),
				archive:        newArchiveMediaPortFake(),
				cataloguerFactory: &ReferencerFactoryFake{
					Cataloguer: &CatalogReferencerFake{
						analysedMedias[2]: doesNotExistReference3,
						analysedMedias[3]: doesNotExistReference4,
					},
				},
				insertMedia: newInsertMediaPortFake(),
			},
			args: args{
				owner: owner,
				volume: &InMemorySourceVolume{
					analysedMedias[2].FoundMedia,
					analysedMedias[3].FoundMedia,
				},
			},
			want: &backupReportBuilder{
				skipped: MediaCounterZero,
				countPerAlbum: map[string]*AlbumReport{
					doesNotExistReference3.AlbumFolderNameValue: {
						video: NewMediaCounter(1, analysedMedias[2].FoundMedia.Size()),
					},
					doesNotExistReference4.AlbumFolderNameValue: {
						isNew: true,
						other: NewMediaCounter(1, analysedMedias[3].FoundMedia.Size()),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should run uploader in 2 routines",
			fields: fields{
				detailsReaders: new(DetailsReaderAdapterStub),
				archive: &ArchiveGroupWaiter{
					Delegate: newArchiveMediaPortFake(),
					Waiter:   newWaitGroup(2),
				},
				cataloguerFactory: &ReferencerFactoryFake{
					Cataloguer: &CatalogReferencerFake{
						analysedMedias[0]: doesNotExistReference1,
						analysedMedias[1]: doesNotExistReference2,
					},
				},
				insertMedia: newInsertMediaPortFake(),
			},
			args: args{
				owner: owner,
				volume: &InMemorySourceVolume{
					analysedMedias[0].FoundMedia,
					analysedMedias[1].FoundMedia,
				},
				optionsSlice: []Options{
					OptionsConcurrentUploaderRoutines(2),
					OptionsBatchSize(1),
				},
			},
			want: &backupReportBuilder{
				skipped: MediaCounterZero,
				countPerAlbum: map[string]*AlbumReport{
					doesNotExistReference1.AlbumFolderNameValue: countOfMedias(analysedMedias[0], analysedMedias[1]),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should copy files into reject folder OptionsRejectFolder is used",
			fields: fields{
				detailsReaders: readerFailingToParseFile,
				archive:        newArchiveMediaPortFake(),
				cataloguerFactory: &ReferencerFactoryFake{
					Cataloguer: &CatalogReferencerFake{},
				},
				insertMedia: newInsertMediaPortFake(),
			},
			args: args{
				owner: owner,
				volume: &InMemorySourceVolume{
					analysedMedias[0].FoundMedia,
				},
				optionsSlice: []Options{
					OptionsWithRejectDir(rejectFolder),
				},
			},
			want: &backupReportBuilder{
				skipped:       NewMediaCounter(1, analysedMedias[0].FoundMedia.Size()),
				countPerAlbum: map[string]*AlbumReport{},
			},
			wantRejectFolderContent: []string{"folder1_file_1.jpg"},
			wantErr:                 assert.NoError,
		},
		{
			name: "it should fail the backup if a media is rejected during analysis",
			fields: fields{
				detailsReaders: readerFailingToParseFile,
				archive:        newArchiveMediaPortFake(),
				cataloguerFactory: &ReferencerFactoryFake{
					Cataloguer: &CatalogReferencerFake{},
				},
				insertMedia: newInsertMediaPortFake(),
			},
			args: args{
				owner: owner,
				volume: &InMemorySourceVolume{
					analysedMedias[0].FoundMedia,
				},
			},
			want: &backupReportBuilder{
				skipped:       NewMediaCounter(1, analysedMedias[0].FoundMedia.Size()),
				countPerAlbum: map[string]*AlbumReport{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrTestNoDetails, i...)
			},
		},
		{
			name: "it should report a media as filtered out because already exists",
			fields: fields{
				detailsReaders: new(DetailsReaderAdapterStub),
				archive:        newArchiveMediaPortFake(),
				cataloguerFactory: &ReferencerFactoryFake{
					Cataloguer: &CatalogReferencerFake{
						analysedMedias[0]: doesExistReference1,
					},
				},
				insertMedia: &AssertInsertMediaPort{
					Want: []InsertMediaPortFakeEntry{},
				},
			},
			args: args{
				owner: owner,
				volume: &InMemorySourceVolume{
					analysedMedias[0].FoundMedia,
				},
			},
			want: &backupReportBuilder{
				skipped:       NewMediaCounter(1, analysedMedias[0].FoundMedia.Size()),
				countPerAlbum: make(map[string]*AlbumReport),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should re-buffer after analysis to have upload batch full",
			fields: fields{
				detailsReaders: new(DetailsReaderAdapterStub),
				archive:        newArchiveMediaPortFake(),
				cataloguerFactory: &ReferencerFactoryFake{
					Cataloguer: &CatalogReferencerFake{
						analysedMedias[0]: doesExistReference1,
						analysedMedias[1]: doesNotExistReference2,
						analysedMedias[2]: doesNotExistReference3,
					},
				},
				insertMedia: &AssertInsertMediaPort{
					Want: []InsertMediaPortFakeEntry{
						{
							owner:           owner,
							ArchiveFilename: fakeArchiveFileName(analysedMedias[1]),
						},
						{
							owner:           owner,
							ArchiveFilename: fakeArchiveFileName(analysedMedias[2]),
						},
					},
					WantBatchCount: 1,
				},
			},
			args: args{
				owner: owner,
				volume: &InMemorySourceVolume{
					analysedMedias[0].FoundMedia,
					analysedMedias[1].FoundMedia,
					analysedMedias[2].FoundMedia,
				},
				optionsSlice: []Options{
					OptionsBatchSize(2),
				},
			},
			want: &backupReportBuilder{
				skipped: NewMediaCounter(1, analysedMedias[0].FoundMedia.Size()),
				countPerAlbum: map[string]*AlbumReport{
					doesNotExistReference2.AlbumFolderNameValue: countOfMedias(analysedMedias[1], analysedMedias[2]),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should use the Analyser decorator, and report accordingly",
			fields: fields{
				detailsReaders: &DetailsReaderFake{}, // will fail if called
				archive:        newArchiveMediaPortFake(),
				cataloguerFactory: &ReferencerFactoryFake{
					Cataloguer: &CatalogReferencerFake{
						analysedMedias[0]: doesNotExistReference1,
					},
				},
				insertMedia: newInsertMediaPortFake(),
			},
			args: args{
				owner: owner,
				volume: &InMemorySourceVolume{
					analysedMedias[0].FoundMedia,
				},
				optionsSlice: []Options{
					OptionsAnalyserDecorator(&AnalyserDecoratorFake{
						Cached: map[string]*AnalysedMedia{
							analysedMedias[0].FoundMedia.MediaPath().Filename: analysedMedias[0],
						},
					}),
				},
			},
			want: &backupReportBuilder{
				skipped: MediaCounterZero,
				countPerAlbum: map[string]*AlbumReport{
					doesNotExistReference1.AlbumFolderNameValue: countOfMedias(analysedMedias[0]),
				},
			},
			wantEvents: map[trackEvent]eventSummary{
				trackScanComplete:      {SumCount: 1, SumSize: 10},
				trackCatalogued:        {SumCount: 1, SumSize: 10},
				trackUploaded:          {SumCount: 1, SumSize: 10, Albums: []string{"/album1"}},
				trackAnalysedFromCache: {SumCount: 1, SumSize: 10},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventCatcher := newEventCapture()
			options := append([]Options{OptionsWithListener(eventCatcher)}, tt.args.optionsSlice...)

			backup := &BatchBackup{
				CataloguerFactory: tt.fields.cataloguerFactory,
				DetailsReaders:    []DetailsReader{tt.fields.detailsReaders},
				InsertMediaPort:   tt.fields.insertMedia,
				ArchivePort:       tt.fields.archive,
			}

			got, err := backup.Backup(context.Background(), tt.args.owner, tt.args.volume, options...)

			if !tt.wantErr(t, err, fmt.Sprintf("Backup(%v, %v, %v)", tt.args.owner, tt.args.volume, tt.args.optionsSlice)) {
				return
			}
			assert.Equalf(t, tt.want, convertToStaticCompletionReport(got), "Backup(%v, %v, %v)", tt.args.owner, tt.args.volume, tt.args.optionsSlice)
			assert.ElementsMatch(t, tt.wantRejectFolderContent, readAndClearFolder(t, rejectFolder))
			if tt.wantEvents != nil {
				assert.Equal(t, tt.wantEvents, eventCatcher.Captured)
			}

			if toBeSatisfied, ok := tt.fields.cataloguerFactory.(ToBeSatisfied); ok {
				toBeSatisfied.IsSatisfied(t)
			}
			if toBeSatisfied, ok := tt.fields.archive.(ToBeSatisfied); ok {
				toBeSatisfied.IsSatisfied(t)
			}
			if toBeSatisfied, ok := tt.fields.insertMedia.(ToBeSatisfied); ok {
				toBeSatisfied.IsSatisfied(t)
			}
		})
	}
}

func readAndClearFolder(t *testing.T, rejectFolder string) []string {
	var filenames []string

	dir, err := os.ReadDir(rejectFolder)
	if assert.NoError(t, err, "read dir", rejectFolder) {
		for _, entry := range dir {
			filenames = append(filenames, entry.Name())
			_ = os.Remove(path.Join(rejectFolder, entry.Name()))
		}
	}
	return filenames
}

func countOfMedias(medias ...*AnalysedMedia) *AlbumReport {
	report := &AlbumReport{}
	for _, media := range medias {
		switch media.Type {
		case MediaTypeImage:
			report.image = report.image.Add(1, media.FoundMedia.Size())
		case MediaTypeVideo:
			report.video = report.video.Add(1, media.FoundMedia.Size())
		default:
			report.other = report.other.Add(1, media.FoundMedia.Size())
		}
	}
	return report
}

func fakeArchiveFileName(media *AnalysedMedia) string {
	return media.FoundMedia.MediaPath().Filename
}

func newArchiveMediaPortFake() *ArchiveMediaPortFake {
	return &ArchiveMediaPortFake{
		lock: sync.Mutex{},
		got:  make(map[ownermodel.Owner][]*BackingUpMediaRequest),
	}
}

type ArchiveMediaPortFake struct {
	lock sync.Mutex
	got  map[ownermodel.Owner][]*BackingUpMediaRequest
}

func (a *ArchiveMediaPortFake) ArchiveMedia(ownerValue string, media *BackingUpMediaRequest) (string, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	owner := ownermodel.Owner(ownerValue)
	a.got[owner] = append(a.got[owner], media)
	return fakeArchiveFileName(media.AnalysedMedia), nil
}

type AssertArchiveFake struct {
	*ArchiveMediaPortFake
	Want map[ownermodel.Owner][]*BackingUpMediaRequest
}

func (a *AssertArchiveFake) IsSatisfied(t *testing.T) bool {
	if !assert.Equal(t, slices.Sorted(maps.Keys(a.Want)), slices.Sorted(maps.Keys(a.ArchiveMediaPortFake.got)), "different list of owners") {
		return false
	}

	passed := true
	for owner, expected := range a.Want {
		passed = passed && assert.ElementsMatch(t, expected, a.ArchiveMediaPortFake.got[owner])
	}
	return passed
}

type InsertMediaPortFakeEntry struct {
	owner           ownermodel.Owner
	ArchiveFilename string
}

func newInsertMediaPortFake() *InsertMediaPortFake {
	return &InsertMediaPortFake{
		lock: sync.Mutex{},
	}
}

type InsertMediaPortFake struct {
	lock       sync.Mutex
	Got        []InsertMediaPortFakeEntry
	BatchCount int
}

func (i *InsertMediaPortFake) IndexMedias(ctx context.Context, owner ownermodel.Owner, requests []*CatalogMediaRequest) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	if len(requests) == 0 {
		return errors.New("InsertMediaPortFake.IndexMedias(ctx, owner, request) request must not be empty")
	}

	i.BatchCount++
	for _, request := range requests {
		i.Got = append(i.Got, InsertMediaPortFakeEntry{
			owner:           owner,
			ArchiveFilename: request.ArchiveFilename,
		})
	}
	return nil
}

type AssertInsertMediaPort struct {
	InsertMediaPortFake
	Want           []InsertMediaPortFakeEntry
	WantBatchCount int
}

func (a *AssertInsertMediaPort) IsSatisfied(t *testing.T) bool {
	return assert.ElementsMatch(t, a.Want, a.InsertMediaPortFake.Got) &&
		(a.WantBatchCount == 0 || assert.Equal(t, a.WantBatchCount, a.BatchCount, "batch count"))
}

func convertToStaticCompletionReport(report Report) *backupReportBuilder {
	return report.(*backupReportBuilder)
}

type ArchiveGroupWaiter struct {
	Delegate ArchiveMediaPort
	Waiter   *sync.WaitGroup
}

func (a *ArchiveGroupWaiter) ArchiveMedia(owner string, media *BackingUpMediaRequest) (string, error) {
	a.Waiter.Done()
	a.Waiter.Wait()

	return a.Delegate.ArchiveMedia(owner, media)
}
