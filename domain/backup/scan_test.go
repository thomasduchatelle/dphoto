package backup_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"github.com/thomasduchatelle/dphoto/mocks"
	"testing"
	"time"
)

func TestShouldReportScannedItems(t *testing.T) {
	// setup
	a := assert.New(t)
	const owner = "tony@stark.com"

	readerAdapter := mockDetailsReaderAdapter(t)
	backup.RegisterDetailsReader(readerAdapter)

	catalogMock := mocks.NewCatalogAdapter(t)
	archiveMock := mocks.NewArchiveAdapter(t)
	backup.Init(catalogMock, archiveMock)
	backup.BatchSize = 128

	eventCapture := newEventCapture()

	// given
	volume := mockVolume{
		backup.NewInmemoryMedia("folder1/file_1.jpg", time.Now(), []byte("2022-06-18")),
		backup.NewInmemoryMedia("folder1/file_2.jpg", time.Now(), []byte("2022-06-18A")),
		backup.NewInmemoryMedia("folder1/file_3.jpg", time.Now(), []byte("2022-06-19AB")),
		backup.NewInmemoryMedia("folder1/file_4.jpg", time.Now(), []byte("2022-06-20ABC")),
		backup.NewInmemoryMedia("folder1/folder1a/file_5.jpg", time.Now(), []byte("2022-06-21ABCD")),
		backup.NewInmemoryMedia("folder2/file_6.jpg", time.Now(), []byte("2022-06-22ABCDE")),
	}
	expectedCatalogRequests := []*backup.BackingUpMediaRequest{
		{
			AnalysedMedia: &backup.AnalysedMedia{
				FoundMedia: volume[0],
				Type:       backup.MediaTypeImage,
				Sha256Hash: "3e7574e8b640104d97597b200fd516c589f34be540e0a81a272fd488d12acaec",
				Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC)},
			},
			Id:         "id_file_1.jpg",
			FolderName: "/album1",
		},
		{
			AnalysedMedia: &backup.AnalysedMedia{
				FoundMedia: volume[1],
				Type:       backup.MediaTypeImage,
				Sha256Hash: "43e41e253022d4e2e4bf3d8388d5cb0e7553b2da3e8495c5e8617c961aa0a0bd",
				Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC)},
			},
			Id:         "id_file_2.jpg",
			FolderName: "/album1",
		},
		{
			AnalysedMedia: &backup.AnalysedMedia{
				FoundMedia: volume[2],
				Type:       backup.MediaTypeImage,
				Sha256Hash: "28f046d0ebae98f45512f98d581e7cdded28dd9cf50e7712615970dc15221cb3",
				Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 19, 0, 0, 0, 0, time.UTC)},
			},
			Id:         "id_file_3.jpg",
			FolderName: "/album1",
		},
		{
			AnalysedMedia: &backup.AnalysedMedia{
				FoundMedia: volume[3],
				Type:       backup.MediaTypeImage,
				Sha256Hash: "b9506fc17d9a648b448efa042a76bcae587e7e2afe02c00c539e5905b9dbb5b3",
				Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 20, 0, 0, 0, 0, time.UTC)},
			},
			Id:         "OUT",
			FolderName: "",
		},
		{
			AnalysedMedia: &backup.AnalysedMedia{
				FoundMedia: volume[4],
				Type:       backup.MediaTypeImage,
				Sha256Hash: "ce2b4c6e0f8cf6c2be15d85925f8e6c79cef5c9fbbe5578e6dd0ae419c222d53",
				Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 21, 0, 0, 0, 0, time.UTC)},
			},
			Id:         "id_file_5.jpg",
			FolderName: "",
		},
		{
			AnalysedMedia: &backup.AnalysedMedia{
				FoundMedia: volume[5],
				Type:       backup.MediaTypeImage,
				Sha256Hash: "248960db17bc3e685260f28c0af7fb3b1b3b8659d476c42ccc2a5871c53ab438",
				Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 22, 0, 0, 0, 0, time.UTC)},
			},
			Id:         "OUT",
			FolderName: "",
		},
	}
	catalogMock.On("AssignIdsToNewMedias", owner, []*backup.AnalysedMedia{
		expectedCatalogRequests[0].AnalysedMedia,
		expectedCatalogRequests[1].AnalysedMedia,
		expectedCatalogRequests[2].AnalysedMedia,
		expectedCatalogRequests[3].AnalysedMedia,
		expectedCatalogRequests[4].AnalysedMedia,
		expectedCatalogRequests[5].AnalysedMedia,
	}).Once().Return(newIdGeneratorWithExclusion(func(name string) bool {
		return name != "file_4.jpg" && name != "file_6.jpg"
	}), nil)

	// when
	folders, skippedMedias, err := backup.Scan(owner, &volume, backup.OptionWithListener(eventCapture))

	// then
	name := "it should find 2 folders"
	if a.NoError(err) {
		for _, folder := range folders {
			folder.Volume = nil
		}

		a.Equal([]*backup.ScannedFolder{
			{
				Name:         "folder1",
				RelativePath: "folder1",
				FolderName:   "folder1",
				Start:        time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC),
				End:          time.Date(2022, 6, 20, 0, 0, 0, 0, time.UTC),
				Distribution: map[string]backup.MediaCounter{
					"2022-06-18": backup.NewMediaCounter(2, 10+11),
					"2022-06-19": backup.NewMediaCounter(1, 12),
				},
			},
			{
				Name:         "folder1a",
				RelativePath: "folder1/folder1a",
				FolderName:   "folder1a",
				Start:        time.Date(2022, 6, 21, 0, 0, 0, 0, time.UTC),
				End:          time.Date(2022, 6, 22, 0, 0, 0, 0, time.UTC),
				Distribution: map[string]backup.MediaCounter{
					"2022-06-21": backup.NewMediaCounter(1, 14),
				},
			},
		}, folders, name)

		name = "it should not have any skipped files because it's not implemented"
		a.Empty(skippedMedias, name)

		name = "it should have generated necessary events to track the progress"
		a.Equal(map[backup.ProgressEventType]eventSummary{
			backup.ProgressEventScanComplete:   {Number: 1, SumCount: 6, SumSize: 10 + 11 + 12 + 13 + 14 + 15},
			backup.ProgressEventAnalysed:       {Number: 6, SumCount: 6, SumSize: 10 + 11 + 12 + 13 + 14 + 15},
			backup.ProgressEventAlreadyExists:  {Number: 1, SumCount: 2},
			backup.ProgressEventCatalogued:     {Number: 1, SumCount: 4, SumSize: 10 + 11 + 12 + 14},
			backup.ProgressEventReadyForUpload: {Number: 4, SumCount: 4, SumSize: 10 + 11 + 12 + 14},
			backup.ProgressEventUploaded:       {Number: 4, SumCount: 4, SumSize: 10 + 11 + 12 + 14, Albums: nil},
		}, eventCapture.Captured, name)
	}
}
