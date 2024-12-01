package backup

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"maps"
	"testing"
	"time"
)

func TestBackupAcceptance(t *testing.T) {
	const owner = ownermodel.Owner("ironman")

	analysedMedias := []*AnalysedMedia{
		{
			FoundMedia: NewInMemoryMedia("folder1/file_1.jpg", time.Now(), []byte("2022-06-18")),
			Type:       MediaTypeImage,
			Sha256Hash: "3e7574e8b640104d97597b200fd516c589f34be540e0a81a272fd488d12acaec",
			Details:    &MediaDetails{DateTime: time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC)},
		},
	}
	doesNotExistReference := &CatalogReferenceStub{MediaIdValue: "media-id-1", AlbumFolderNameValue: "/album1"}

	type fields struct {
		archive           BArchiveAdapter
		cataloguerFactory CataloguerFactory
		insertMedia       InsertMediaPort
		detailsReaders    DetailsReaderAdapter
	}
	type args struct {
		owner        ownermodel.Owner
		volume       SourceVolume
		optionsSlice []Options
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    CompletionReport
		wantErr assert.ErrorAssertionFunc
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
								CatalogReference: doesNotExistReference,
							},
						},
					},
				},
				cataloguerFactory: &ReferencerFactoryFake{
					CreatorReferencer: &CatalogReferencerFake{
						analysedMedias[0]: doesNotExistReference,
					},
				},
				insertMedia: &AssertInsertMediaPort{
					Want: []InsertMediaPortFakeEntry{
						{
							owner:           owner,
							reference:       doesNotExistReference,
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
			want: &StaticCompletionReport{
				skipped: MediaCounterZero,
				countPerAlbum: map[string]IAlbumReport{
					doesNotExistReference.AlbumFolderNameValue: countOfMedias(analysedMedias[0]),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backup := &BatchBackup{
				CataloguerFactory: tt.fields.cataloguerFactory,
				DetailsReaders:    []DetailsReaderAdapter{tt.fields.detailsReaders},
				InsertMediaPort:   tt.fields.insertMedia,
				ArchivePort:       tt.fields.archive,
			}

			got, err := backup.Backup(context.Background(), tt.args.owner, tt.args.volume, tt.args.optionsSlice...)

			if !tt.wantErr(t, err, fmt.Sprintf("Backup(%v, %v, %v)", tt.args.owner, tt.args.volume, tt.args.optionsSlice)) {
				return
			}
			assert.Equalf(t, tt.want, convertToStaticCompletionReport(got), "Backup(%v, %v, %v)", tt.args.owner, tt.args.volume, tt.args.optionsSlice)
		})
	}
}

func countOfMedias(medias ...*AnalysedMedia) IAlbumReport {
	report := &StaticAlbumReport{}
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

func staticCountWithSingleReport(album string, created bool) map[string]*AlbumReport {
	return map[string]*AlbumReport{
		"/album1": {
			New:    false,
			counts: [3]int{},
			sizes:  [3]int{},
		},
	}
}

func fakeArchiveFileName(media *AnalysedMedia) string {
	return media.FoundMedia.MediaPath().Filename
}

type ArchiveFake struct {
	got map[ownermodel.Owner][]*BackingUpMediaRequest
}

func (a *ArchiveFake) ArchiveMedia(ownerValue string, media *BackingUpMediaRequest) (string, error) {
	if a.got == nil {
		a.got = make(map[ownermodel.Owner][]*BackingUpMediaRequest)
	}

	owner := ownermodel.Owner(ownerValue)
	a.got[owner] = append(a.got[owner], media)
	return fakeArchiveFileName(media.AnalysedMedia), nil
}

type AssertArchiveFake struct {
	ArchiveFake
	Want map[ownermodel.Owner][]*BackingUpMediaRequest
}

func (a *AssertArchiveFake) IsSatisfied(t *testing.T) bool {
	if !assert.Equal(t, maps.Keys(a.Want), maps.Keys(a.ArchiveFake.got), "different list of owners") {
		return false
	}

	passed := true
	for owner, expected := range a.Want {
		passed = passed && assert.ElementsMatch(t, expected, a.ArchiveFake.got[owner])
	}
	return passed
}

type InsertMediaPortFakeEntry struct {
	owner           ownermodel.Owner
	reference       CatalogReference
	ArchiveFilename string
}

type InsertMediaPortFake struct {
	got []InsertMediaPortFakeEntry
}

func (i *InsertMediaPortFake) IndexMedias(ctx context.Context, owner ownermodel.Owner, requests []*CatalogMediaRequest) error {
	for _, request := range requests {
		i.got = append(i.got, InsertMediaPortFakeEntry{
			owner:           owner,
			reference:       request.BackingUpMediaRequest.CatalogReference,
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

func convertToStaticCompletionReport(report CompletionReport) *StaticCompletionReport {
	countPerAlbum := make(map[string]IAlbumReport)
	for album, albumReport := range report.CountPerAlbum() {
		countPerAlbum[album] = convertToStaticAlbumReport(albumReport)
	}

	return &StaticCompletionReport{
		skipped:       report.Skipped(),
		countPerAlbum: countPerAlbum,
	}
}

type StaticCompletionReport struct {
	skipped       MediaCounter
	countPerAlbum map[string]IAlbumReport
}

func (e *StaticCompletionReport) Skipped() MediaCounter {
	return e.skipped
}

func (e *StaticCompletionReport) CountPerAlbum() map[string]IAlbumReport {
	return e.countPerAlbum
}

func convertToStaticAlbumReport(report IAlbumReport) *StaticAlbumReport {
	return &StaticAlbumReport{
		isNew: report.IsNew(),
		image: report.OfType(MediaTypeImage),
		video: report.OfType(MediaTypeVideo),
		other: report.OfType(MediaTypeOther),
	}
}

type StaticAlbumReport struct {
	isNew bool
	image MediaCounter
	video MediaCounter
	other MediaCounter
}

func (c *StaticAlbumReport) IsNew() bool {
	return c.isNew
}

func (c *StaticAlbumReport) Total() MediaCounter {
	return c.image.AddCounter(c.video).AddCounter(c.other)
}

func (c *StaticAlbumReport) OfType(mediaType MediaType) MediaCounter {
	switch mediaType {
	case MediaTypeImage:
		return c.image
	case MediaTypeVideo:
		return c.video
	default:
		return c.other
	}
}
