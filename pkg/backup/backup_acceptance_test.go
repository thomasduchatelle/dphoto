package backup

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"maps"
	"os"
	"path"
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
	}
	doesNotExistReference1 := &CatalogReferenceStub{MediaIdValue: "media-id-1", AlbumFolderNameValue: "/album1"}
	doesNotExistReference2 := &CatalogReferenceStub{MediaIdValue: "media-id-2", AlbumFolderNameValue: "/album1"}

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
		wantErr                 assert.ErrorAssertionFunc
		wantRejectFolderContent []string
	}{
		{
			name: "it should upload a media going through the happy path",
			fields: fields{
				detailsReaders: new(DetailsReaderAdapterStub),
				archive: &AssertArchiveFake{
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
					CreatorReferencer: &CatalogReferencerFake{
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
			wantErr: assert.NoError,
		},
		{
			name: "it should run uploader in 2 routines",
			fields: fields{
				detailsReaders: new(DetailsReaderAdapterStub),
				archive: &ArchiveGroupWaiter{
					Delegate: new(ArchiveMediaPortFake),
					Waiter:   newWaitGroup(2),
				},
				cataloguerFactory: &ReferencerFactoryFake{
					CreatorReferencer: &CatalogReferencerFake{
						analysedMedias[0]: doesNotExistReference1,
						analysedMedias[1]: doesNotExistReference2,
					},
				},
				insertMedia: new(InsertMediaPortFake),
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
				archive:        new(ArchiveMediaPortFake),
				cataloguerFactory: &ReferencerFactoryFake{
					CreatorReferencer: &CatalogReferencerFake{},
				},
				insertMedia: new(InsertMediaPortFake),
			},
			args: args{
				owner: owner,
				volume: &InMemorySourceVolume{
					analysedMedias[0].FoundMedia,
				},
				optionsSlice: []Options{
					OptionWithRejectDir(rejectFolder),
				},
			},
			want: &backupReportBuilder{
				skipped:       NewMediaCounter(1, analysedMedias[0].FoundMedia.Size()),
				countPerAlbum: map[string]*AlbumReport{},
			},
			wantRejectFolderContent: []string{"folder1_file_1.jpg"},
			wantErr:                 assert.NoError,
		},
		// TODO dry run
		// TODO interrupt if errors
		// TODO report OnFilteredOut
		// TODO re-buffer before uploading
		// TODO video in report
		// TODO other in report
		// TODO options WithCachedAnalysis
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backup := &BatchBackup{
				CataloguerFactory: tt.fields.cataloguerFactory,
				DetailsReaders:    []DetailsReader{tt.fields.detailsReaders},
				InsertMediaPort:   tt.fields.insertMedia,
				ArchivePort:       tt.fields.archive,
			}

			got, err := backup.Backup(context.Background(), tt.args.owner, tt.args.volume, tt.args.optionsSlice...)

			if !tt.wantErr(t, err, fmt.Sprintf("Backup(%v, %v, %v)", tt.args.owner, tt.args.volume, tt.args.optionsSlice)) {
				return
			}
			assert.Equalf(t, tt.want, convertToStaticCompletionReport(got), "Backup(%v, %v, %v)", tt.args.owner, tt.args.volume, tt.args.optionsSlice)

			assert.ElementsMatch(t, tt.wantRejectFolderContent, readAndClearFolder(t, rejectFolder))
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

type ArchiveMediaPortFake struct {
	got map[ownermodel.Owner][]*BackingUpMediaRequest
}

func (a *ArchiveMediaPortFake) ArchiveMedia(ownerValue string, media *BackingUpMediaRequest) (string, error) {
	if a.got == nil {
		a.got = make(map[ownermodel.Owner][]*BackingUpMediaRequest)
	}

	owner := ownermodel.Owner(ownerValue)
	a.got[owner] = append(a.got[owner], media)
	return fakeArchiveFileName(media.AnalysedMedia), nil
}

type AssertArchiveFake struct {
	ArchiveMediaPortFake
	Want map[ownermodel.Owner][]*BackingUpMediaRequest
}

func (a *AssertArchiveFake) IsSatisfied(t *testing.T) bool {
	if !assert.Equal(t, maps.Keys(a.Want), maps.Keys(a.ArchiveMediaPortFake.got), "different list of owners") {
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

type InsertMediaPortFake struct {
	got []InsertMediaPortFakeEntry
}

func (i *InsertMediaPortFake) IndexMedias(ctx context.Context, owner ownermodel.Owner, requests []*CatalogMediaRequest) error {
	for _, request := range requests {
		i.got = append(i.got, InsertMediaPortFakeEntry{
			owner:           owner,
			ArchiveFilename: request.ArchiveFilename,
		})
	}
	return nil
}

type AssertInsertMediaPort struct {
	InsertMediaPortFake
	Want []InsertMediaPortFakeEntry
}

func (a *AssertInsertMediaPort) IsSatisfied(t *testing.T) bool {
	return assert.ElementsMatch(t, a.Want, a.InsertMediaPortFake.got)
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
