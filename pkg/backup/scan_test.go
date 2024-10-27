package backup

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"io"
	"testing"
	"time"
)

func TestScanAcceptance(t *testing.T) {
	owner := ownermodel.Owner("tony@stark.com")

	analysedMedias := []*AnalysedMedia{
		{
			FoundMedia: NewInMemoryMedia("folder1/file_1.jpg", time.Now(), []byte("2022-06-18")),
			Type:       MediaTypeImage,
			Sha256Hash: "3e7574e8b640104d97597b200fd516c589f34be540e0a81a272fd488d12acaec",
			Details:    &MediaDetails{DateTime: time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC)},
		},
		{
			FoundMedia: NewInMemoryMedia("folder1/file_2.jpg", time.Now(), []byte("2022-06-18A")),
			Type:       MediaTypeImage,
			Sha256Hash: "43e41e253022d4e2e4bf3d8388d5cb0e7553b2da3e8495c5e8617c961aa0a0bd",
			Details:    &MediaDetails{DateTime: time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC)},
		},
		{
			FoundMedia: NewInMemoryMedia("folder1/file_3.jpg", time.Now(), []byte("2022-06-19AB")),
			Type:       MediaTypeImage,
			Sha256Hash: "28f046d0ebae98f45512f98d581e7cdded28dd9cf50e7712615970dc15221cb3",
			Details:    &MediaDetails{DateTime: time.Date(2022, 6, 19, 0, 0, 0, 0, time.UTC)},
		},
		{
			FoundMedia: NewInMemoryMedia("folder1/file_4.jpg", time.Now(), []byte("2022-06-20ABC")),
			Type:       MediaTypeImage,
			Sha256Hash: "b9506fc17d9a648b448efa042a76bcae587e7e2afe02c00c539e5905b9dbb5b3",
			Details:    &MediaDetails{DateTime: time.Date(2022, 6, 20, 0, 0, 0, 0, time.UTC)},
		},
		{
			FoundMedia: NewInMemoryMedia("folder1/folder1a/file_5.jpg", time.Now(), []byte("2022-06-21ABCD")),
			Type:       MediaTypeImage,
			Sha256Hash: "ce2b4c6e0f8cf6c2be15d85925f8e6c79cef5c9fbbe5578e6dd0ae419c222d53",
			Details:    &MediaDetails{DateTime: time.Date(2022, 6, 21, 0, 0, 0, 0, time.UTC)},
		},
		{
			FoundMedia: NewInMemoryMedia("folder2/file_6.jpg", time.Now(), []byte("2022-06-22ABCDE")),
			Type:       MediaTypeImage,
			Sha256Hash: "248960db17bc3e685260f28c0af7fb3b1b3b8659d476c42ccc2a5871c53ab438",
			Details:    &MediaDetails{DateTime: time.Date(2022, 6, 22, 0, 0, 0, 0, time.UTC)},
		},
		{
			FoundMedia: NewInMemoryMedia("folder1/file_7_no_date.jpg", time.Now(), []byte{0}),
			Type:       MediaTypeImage,
			Sha256Hash: "28f046d0ebae98f45512f98d581e7cdded28dd9cf50e7712615970dc15221cb3",
			Details:    &MediaDetails{},
		},
	}

	volumeStub := make(SourceVolumeStub, 0)
	for _, analysedMedia := range analysedMedias {
		volumeStub = append(volumeStub, analysedMedia.FoundMedia)
	}

	type fields struct {
		detailsReaders    DetailsReaderAdapter
		referencerFactory ReferencerFactory
	}
	type args struct {
		owner       string
		volume      SourceVolume
		optionSlice []Options
	}
	tests := []struct {
		name                   string
		fields                 fields
		args                   args
		wantFolders            []*ScannedFolder
		wantSkippedMediasCount int
		wantEvents             map[trackEvent]eventSummary
		wantErr                assert.ErrorAssertionFunc
	}{
		{
			name: "it should scan files per folder, using read-only referencer, ignoring files that are already catalogued",
			fields: fields{
				detailsReaders: new(DetailsReaderAdapterStub),
				referencerFactory: &ReferencerFactoryFake{
					DryRunReferencer: &CatalogReferencerFake{
						analysedMedias[0]: &CatalogReferenceStub{MediaIdValue: "media-id-1", AlbumFolderNameValue: "/album1"},
						analysedMedias[1]: &CatalogReferenceStub{MediaIdValue: "media-id-2", AlbumFolderNameValue: "/album1"},
						analysedMedias[2]: &CatalogReferenceStub{MediaIdValue: "media-id-3", AlbumFolderNameValue: "/album1"},
						analysedMedias[3]: &CatalogReferenceStub{MediaIdValue: "media-id-4", AlbumFolderNameValue: "/album3", ExistsValue: true},
						analysedMedias[4]: &CatalogReferenceStub{MediaIdValue: "media-id-5", AlbumFolderNameValue: "/album2", AlbumCreatedValue: true},
						analysedMedias[5]: &CatalogReferenceStub{MediaIdValue: "media-id-6", AlbumFolderNameValue: "/album3", ExistsValue: true},
					},
				},
			},
			args: args{
				owner:  owner.Value(),
				volume: &volumeStub,
				optionSlice: []Options{
					OptionsSkipRejects(true),
					OptionsBatchSize(3),
				},
			},
			wantFolders: []*ScannedFolder{
				{
					Name:         "folder1",
					RelativePath: "folder1",
					FolderName:   "folder1",
					AbsolutePath: "/ram/folder1",
					Start:        time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC),
					End:          time.Date(2022, 6, 20, 0, 0, 0, 0, time.UTC),
					Distribution: map[string]MediaCounter{
						"2022-06-18": NewMediaCounter(2, 10+11),
						"2022-06-19": NewMediaCounter(1, 12),
					},
					RejectsCount: 1,
				},
				{
					Name:         "folder1a",
					RelativePath: "folder1/folder1a",
					FolderName:   "folder1a",
					AbsolutePath: "/ram/folder1/folder1a",
					Start:        time.Date(2022, 6, 21, 0, 0, 0, 0, time.UTC),
					End:          time.Date(2022, 6, 22, 0, 0, 0, 0, time.UTC),
					Distribution: map[string]MediaCounter{
						"2022-06-21": NewMediaCounter(1, 14),
					},
				},
			},
			wantEvents: map[trackEvent]eventSummary{
				trackAlbumCreated:           {SumCount: 1, Albums: []string{"/album2"}},
				trackScanComplete:           {SumCount: 7, SumSize: 10 + 11 + 12 + 13 + 14 + 15 + 1},
				trackAnalysisFailed:         {SumCount: 1, SumSize: 1},
				trackAlreadyExistsInCatalog: {SumCount: 2, SumSize: 13 + 15},
				trackCatalogued:             {SumCount: 4, SumSize: 10 + 11 + 12 + 14},
			},
			wantSkippedMediasCount: 1,
			wantErr:                assert.NoError,
		},
		{
			name: "it should ignore the files that are not readable by the analyser",
			fields: fields{
				detailsReaders: new(DetailsReaderAdapterStub),
				referencerFactory: &ReferencerFactoryFake{
					DryRunReferencer: &CatalogReferencerFake{},
				},
			},
			args: args{
				owner: owner.Value(),
				volume: &SourceVolumeStub{
					&UnreadableMedia{FoundMedia: NewInMemoryMedia("folder66/file_unreadable.jpg", time.Now(), []byte("will not be readable"))},
				},
				optionSlice: []Options{
					OptionsSkipRejects(true),
				},
			},
			wantFolders: []*ScannedFolder{
				{
					Name:         "folder66",
					RelativePath: "folder66",
					FolderName:   "folder66",
					AbsolutePath: "/ram/folder66",
					Distribution: make(map[string]MediaCounter),
					RejectsCount: 1,
				},
			},
			wantEvents: map[trackEvent]eventSummary{
				trackScanComplete:   {SumCount: 1, SumSize: 20},
				trackAnalysisFailed: {SumCount: 1, SumSize: 20},
			},
			wantSkippedMediasCount: 1,
			wantErr:                assert.NoError,
		},
		{
			name: "it should fail the scan if one of the files is unreadable",
			fields: fields{
				detailsReaders: new(DetailsReaderAdapterStub),
				referencerFactory: &ReferencerFactoryFake{
					DryRunReferencer: &CatalogReferencerFake{},
				},
			},
			args: args{
				owner: owner.Value(),
				volume: &SourceVolumeStub{
					&UnreadableMedia{FoundMedia: NewInMemoryMedia("folder66/file_unreadable.jpg", time.Now(), []byte("will not be readable"))},
				},
			},
			wantFolders: []*ScannedFolder{},
			wantEvents: map[trackEvent]eventSummary{
				trackScanComplete:   {SumCount: 1, SumSize: 20},
				trackAnalysisFailed: {SumCount: 1, SumSize: 20},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "[UnreadableMedia]", i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(nil, tt.fields.referencerFactory, nil)
			ClearDetailsReader()
			RegisterDetailsReader(tt.fields.detailsReaders)

			eventCatcher := newEventCapture()
			options := append([]Options{OptionWithListener(eventCatcher)}, tt.args.optionSlice...)

			scanner := new(BatchScanner)
			gotFolder, err := scanner.Scan(context.Background(), ownermodel.Owner(tt.args.owner), tt.args.volume, options...)
			if !tt.wantErr(t, err, fmt.Sprintf("Scan(%v, %v, %v)", tt.args.owner, tt.args.volume, options)) {
				return
			}
			assert.Equalf(t, tt.wantFolders, gotFolder, "Scan(%v, %v, %v)", tt.args.owner, tt.args.volume, options)
			assert.Equalf(t, tt.wantEvents, eventCatcher.Captured, "Scan(%v, %v, %v)", tt.args.owner, tt.args.volume, options)

			rejectsCount := 0
			for _, folder := range gotFolder {
				rejectsCount += folder.RejectsCount
			}
			assert.Equalf(t, tt.wantSkippedMediasCount, rejectsCount, "Scan(%v, %v, %v)", tt.args.owner, tt.args.volume, options)
		})
	}
}

type DetailsReaderAdapterStub struct {
}

func (d *DetailsReaderAdapterStub) Supports(media FoundMedia, mediaType MediaType) bool {
	return true
}

func (d *DetailsReaderAdapterStub) ReadDetails(reader io.Reader, options DetailsReaderOptions) (*MediaDetails, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if len(content) < 10 {
		return &MediaDetails{}, nil // no date
	}

	datetime, err := time.Parse("2006-01-02", string(content)[:10])
	if err != nil {
		return nil, err
	}

	return &MediaDetails{
		DateTime: datetime,
	}, nil
}

type SourceVolumeStub []FoundMedia

func (m *SourceVolumeStub) String() string {
	return "Mocked Volume"
}

func (m *SourceVolumeStub) FindMedias() ([]FoundMedia, error) {
	return *m, nil
}

func (m *SourceVolumeStub) Children(MediaPath) (SourceVolume, error) {
	return m, errors.New("SourceVolumeStub cannot generate procreate")
}

type UnreadableMedia struct {
	FoundMedia
}

func (u *UnreadableMedia) ReadMedia() (io.ReadCloser, error) {
	return nil, errors.New("[UnreadableMedia] stubbed error")
}
