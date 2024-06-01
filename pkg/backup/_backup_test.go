package backup_test

import (
	"github.com/stretchr/testify/assert"
	mocks "github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"testing"
	"time"
)

// TODO Backup should have Acceptance tests but not these ones. They are too implementation specific.

func TestShouldCreateAlbumsDuringBackup(t *testing.T) {
	// setup
	a := assert.New(t)
	const owner = "tony@stark.com"

	readerAdapter := mockDetailsReaderAdapter(t)
	backup.RegisterDetailsReader(readerAdapter)

	catalogMock := mocks.NewCatalogAdapter(t)
	timelineMock := mocks.NewTimelineAdapter(t)
	catalogMock.On("GetAlbumsTimeline", owner).Return(timelineMock, nil)
	archiveMock := mocks.NewBArchiveAdapter(t)
	backup.Init(catalogMock, archiveMock, nil)
	backup.BatchSize = 4

	eventCapture := newEventCapture()

	// given
	volume := SourceVolumeStub{
		backup.NewInMemoryMedia("file_1.jpg", time.Now(), []byte("2022-06-18")),
		backup.NewInMemoryMedia("file_2.mp4", time.Now(), []byte("2022-06-19A")),
		backup.NewInMemoryMedia("file_3.avi", time.Now(), []byte("2022-06-20AB")),
		backup.NewInMemoryMedia("file_4.jpg", time.Now(), []byte("2022-06-21ABC")),
		backup.NewInMemoryMedia("file_4.jpg", time.Now(), []byte("2022-06-22ABCD")),
	}
	expectedCatalogRequests := []*backup.CatalogMediaRequest{
		{
			BackingUpMediaRequest: &backup.BackingUpMediaRequest{
				AnalysedMedia: &backup.AnalysedMedia{
					FoundMedia: volume[0],
					Type:       backup.MediaTypeImage,
					Sha256Hash: "3e7574e8b640104d97597b200fd516c589f34be540e0a81a272fd488d12acaec",
					Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC)},
				},
				Id:         "id_file_1.jpg",
				FolderName: "/folder1",
			},
			ArchiveFilename: "new_file_1_name.jpg",
		},
		{
			BackingUpMediaRequest: &backup.BackingUpMediaRequest{
				AnalysedMedia: &backup.AnalysedMedia{
					FoundMedia: volume[1],
					Type:       backup.MediaTypeVideo,
					Sha256Hash: "794b2988415566b2d2f7d8f7d94bc188fba62ad6dfccef7c2446bde8cac86ec5",
					Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 19, 0, 0, 0, 0, time.UTC)},
				},
				Id:         "id_file_2.mp4",
				FolderName: "/folder2",
			},
			ArchiveFilename: "file_2.mp4",
		},
		{
			BackingUpMediaRequest: &backup.BackingUpMediaRequest{
				AnalysedMedia: &backup.AnalysedMedia{
					FoundMedia: volume[2],
					Type:       backup.MediaTypeVideo,
					Sha256Hash: "5853eb19ae52312e6c4750ee9409a7f378d6acfe897211b33b3b62b43694de91",
					Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 20, 0, 0, 0, 0, time.UTC)},
				},
				Id:         "OUT",
				FolderName: "",
			},
			ArchiveFilename: "",
		},
		{
			BackingUpMediaRequest: &backup.BackingUpMediaRequest{
				AnalysedMedia: &backup.AnalysedMedia{
					FoundMedia: volume[3],
					Type:       backup.MediaTypeImage,
					Sha256Hash: "9478d0ee02b23ad6d3e8c5051d23dec82526e42265c1aeec2f211294a835b562",
					Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 21, 0, 0, 0, 0, time.UTC)},
				},
				Id:         "id_file_4.jpg",
				FolderName: "/folder2",
			},
			ArchiveFilename: "",
		},
		{
			BackingUpMediaRequest: &backup.BackingUpMediaRequest{
				AnalysedMedia: &backup.AnalysedMedia{
					FoundMedia: volume[4],
					Type:       backup.MediaTypeImage,
					Sha256Hash: "a0bd7ffade864521e862647823c6cef0cdbc8a7695b9ab72267d4a576b48c998",
					Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 22, 0, 0, 0, 0, time.UTC)},
				},
				Id:         "id_file_4.jpg",
				FolderName: "/folder2",
			},
			ArchiveFilename: "",
		},
	}
	catalogMock.On("AssignIdsToNewMedias", owner, []*backup.AnalysedMedia{
		expectedCatalogRequests[0].BackingUpMediaRequest.AnalysedMedia,
		expectedCatalogRequests[1].BackingUpMediaRequest.AnalysedMedia,
		expectedCatalogRequests[2].BackingUpMediaRequest.AnalysedMedia,
		expectedCatalogRequests[3].BackingUpMediaRequest.AnalysedMedia,
	}).Once().Return(newIdGeneratorWithExclusion(func(filename string) bool {
		return filename != "file_3.avi"
	}), nil)
	catalogMock.On("AssignIdsToNewMedias", owner, []*backup.AnalysedMedia{
		expectedCatalogRequests[4].BackingUpMediaRequest.AnalysedMedia,
	}).Once().Return(newIdGeneratorWithExclusion(func(name string) bool {
		return true
	}), nil)

	timelineMock.On("FindOrCreateAlbum", time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC)).Return("/folder1", false, nil)
	timelineMock.On("FindOrCreateAlbum", time.Date(2022, 6, 19, 0, 0, 0, 0, time.UTC)).Return("/folder2", true, nil)
	timelineMock.On("FindOrCreateAlbum", time.Date(2022, 6, 21, 0, 0, 0, 0, time.UTC)).Return("/folder2", false, nil)
	timelineMock.On("FindOrCreateAlbum", time.Date(2022, 6, 22, 0, 0, 0, 0, time.UTC)).Return("/folder2", false, nil)

	archiveMock.On("ArchiveMedia", owner, expectedCatalogRequests[0].BackingUpMediaRequest).Return(expectedCatalogRequests[0].ArchiveFilename, nil)
	archiveMock.On("ArchiveMedia", owner, expectedCatalogRequests[1].BackingUpMediaRequest).Return(expectedCatalogRequests[1].ArchiveFilename, nil)
	archiveMock.On("ArchiveMedia", owner, expectedCatalogRequests[3].BackingUpMediaRequest).Return(expectedCatalogRequests[3].ArchiveFilename, nil)

	catalogMock.On("IndexMedias", owner, []*backup.CatalogMediaRequest{expectedCatalogRequests[0], expectedCatalogRequests[1], expectedCatalogRequests[3]}).Return(nil)

	// when
	report, err := backup.Backup(owner, &volume, backup.OptionWithListener(eventCapture))

	// then
	if a.NoError(err) {
		a.Equal(map[string]*backup.TypeCounter{
			"/folder1": backup.NewTypeCounter(backup.MediaTypeImage, 1, 10),
			"/folder2": backup.NewTypeCounter(backup.MediaTypeImage, 1, 13).IncrementFoundCounter(backup.MediaTypeVideo, 1, 11),
		}, report.CountPerAlbum())

		a.Equal(map[backup.ProgressEventType]eventSummary{
			backup.ProgressEventAlbumCreated:   {Number: 1, SumCount: 1, Albums: []string{"/folder2"}},
			backup.ProgressEventScanComplete:   {Number: 1, SumCount: 5, SumSize: 10 + 11 + 12 + 13 + 14},
			backup.ProgressEventAnalysed:       {Number: 5, SumCount: 5, SumSize: 10 + 11 + 12 + 13 + 14},
			backup.ProgressEventAlreadyExists:  {Number: 1, SumCount: 1, SumSize: 12},
			backup.ProgressEventCatalogued:     {Number: 2, SumCount: 4, SumSize: 10 + 11 + 13 + 14},
			backup.ProgressEventDuplicate:      {Number: 1, SumCount: 1, SumSize: 14},
			backup.ProgressEventReadyForUpload: {Number: 3, SumCount: 3, SumSize: 10 + 11 + 13},
			backup.ProgressEventUploaded:       {Number: 3, SumCount: 3, SumSize: 10 + 11 + 13, Albums: []string{"/folder1", "/folder2", "/folder2"}},
		}, eventCapture.Captured)
	}
}

func TestShouldFilterMediasBasedOnAlbumDuringBackup(t *testing.T) {
	// setup
	a := assert.New(t)
	const owner = "tony@stark.com"

	readerAdapter := mockDetailsReaderAdapter(t)
	backup.RegisterDetailsReader(readerAdapter)

	catalogMock := mocks.NewCatalogAdapter(t)
	timelineMock := mocks.NewTimelineAdapter(t)
	catalogMock.On("GetAlbumsTimeline", owner).Return(timelineMock, nil)
	archiveMock := mocks.NewBArchiveAdapter(t)
	backup.Init(archiveMock, nil)
	backup.BatchSize = 4

	eventCapture := newEventCapture()

	// given
	volume := SourceVolumeStub{
		backup.NewInMemoryMedia("file_1.jpg", time.Now(), []byte("2022-06-18")),
		backup.NewInMemoryMedia("file_2.jpg", time.Now(), []byte("2022-06-19A")),
		backup.NewInMemoryMedia("file_3.jpg", time.Now(), []byte("2022-06-20AB")),
		backup.NewInMemoryMedia("file_4.jpg", time.Now(), []byte("2022-06-21ABC")),
	}
	expectedCatalogRequests := []*backup.CatalogMediaRequest{
		{
			BackingUpMediaRequest: &backup.BackingUpMediaRequest{
				AnalysedMedia: &backup.AnalysedMedia{
					FoundMedia: volume[0],
					Type:       backup.MediaTypeImage,
					Sha256Hash: "3e7574e8b640104d97597b200fd516c589f34be540e0a81a272fd488d12acaec",
					Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC)},
				},
				Id:         "id_file_1.jpg",
				FolderName: "/folder1",
			},
			ArchiveFilename: "file_1.jpg",
		},
		{
			BackingUpMediaRequest: &backup.BackingUpMediaRequest{
				AnalysedMedia: &backup.AnalysedMedia{
					FoundMedia: volume[1],
					Type:       backup.MediaTypeImage,
					Sha256Hash: "794b2988415566b2d2f7d8f7d94bc188fba62ad6dfccef7c2446bde8cac86ec5",
					Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 19, 0, 0, 0, 0, time.UTC)},
				},
				Id:         "FILTERED_OUT",
				FolderName: "/folder2",
			},
			ArchiveFilename: "",
		},
		{
			BackingUpMediaRequest: &backup.BackingUpMediaRequest{
				AnalysedMedia: &backup.AnalysedMedia{
					FoundMedia: volume[2],
					Type:       backup.MediaTypeImage,
					Sha256Hash: "5853eb19ae52312e6c4750ee9409a7f378d6acfe897211b33b3b62b43694de91",
					Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 20, 0, 0, 0, 0, time.UTC)},
				},
				Id:         "FILTERED_OUT",
				FolderName: "/folder3",
			},
			ArchiveFilename: "",
		},
		{
			BackingUpMediaRequest: &backup.BackingUpMediaRequest{
				AnalysedMedia: &backup.AnalysedMedia{
					FoundMedia: volume[3],
					Type:       backup.MediaTypeImage,
					Sha256Hash: "9478d0ee02b23ad6d3e8c5051d23dec82526e42265c1aeec2f211294a835b562",
					Details:    &backup.MediaDetails{DateTime: time.Date(2022, 6, 21, 0, 0, 0, 0, time.UTC)},
				},
				Id:         "FILTERED_OUT",
				FolderName: "/folder3",
			},
			ArchiveFilename: "",
		},
	}
	catalogMock.On("AssignIdsToNewMedias", owner, []*backup.AnalysedMedia{
		expectedCatalogRequests[0].BackingUpMediaRequest.AnalysedMedia,
		expectedCatalogRequests[2].BackingUpMediaRequest.AnalysedMedia,
	}).Once().Return(newIdGeneratorWithExclusion(func(filename string) bool {
		return filename != "file_3.jpg"
	}), nil)

	timelineMock.On("FindAlbum", time.Date(2022, 6, 18, 0, 0, 0, 0, time.UTC)).Once().Return("/folder1", true, nil)
	timelineMock.On("FindAlbum", time.Date(2022, 6, 19, 0, 0, 0, 0, time.UTC)).Once().Return("/folder2", true, nil)
	timelineMock.On("FindAlbum", time.Date(2022, 6, 20, 0, 0, 0, 0, time.UTC)).Once().Return("/folder1", true, nil)
	timelineMock.On("FindAlbum", time.Date(2022, 6, 21, 0, 0, 0, 0, time.UTC)).Once().Return("/folder1", false, nil)

	archiveMock.On("ArchiveMedia", owner, expectedCatalogRequests[0].BackingUpMediaRequest).Once().Return(expectedCatalogRequests[0].ArchiveFilename, nil)

	catalogMock.On("IndexMedias", owner, expectedCatalogRequests[0:1]).Once().Return(nil)

	// when
	report, err := backup.Backup(owner, &volume, backup.OptionWithListener(eventCapture), backup.OptionOnlyAlbums("/folder1"))

	// then
	if a.NoError(err) {
		a.Equal(map[string]*backup.TypeCounter{
			"/folder1": backup.NewTypeCounter(backup.MediaTypeImage, 1, 10),
		}, report.CountPerAlbum())

		a.Equal(map[backup.ProgressEventType]eventSummary{
			backup.ProgressEventScanComplete:   {Number: 1, SumCount: 4, SumSize: 10 + 11 + 12 + 13},
			backup.ProgressEventAnalysed:       {Number: 4, SumCount: 4, SumSize: 10 + 11 + 12 + 13},
			backup.ProgressEventWrongAlbum:     {Number: 2, SumCount: 2, SumSize: 11 + 13},
			backup.ProgressEventAlreadyExists:  {Number: 1, SumCount: 1, SumSize: 12},
			backup.ProgressEventCatalogued:     {Number: 1, SumCount: 1, SumSize: 10},
			backup.ProgressEventReadyForUpload: {Number: 1, SumCount: 1, SumSize: 10},
			backup.ProgressEventUploaded:       {Number: 1, SumCount: 1, SumSize: 10, Albums: []string{"/folder1"}},
		}, eventCapture.Captured)
	}
}
