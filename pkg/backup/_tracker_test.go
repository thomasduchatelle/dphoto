package backup_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"testing"
)

func TestTracker(t *testing.T) {
	noExtra := backup.ExtraCounts{
		Cached:   backup.MediaCounterZero,
		Rejected: backup.MediaCounterZero,
	}

	tests := []struct {
		name              string
		events            []*backup.progressEvent
		setMocks          func(*mocks.TrackScanComplete, *mocks.TrackAnalysed, *mocks.TrackUploaded)
		wantCountPerAlbum map[string]*backup.AlbumReport
		wantNewAlbums     []string
	}{
		{
			name: "it should call OnAnalysed and OnUploaded if everything follow the happy-order",
			events: []*backup.progressEvent{
				{Type: backup.trackScanComplete, Count: 1, Size: 42},
				{Type: backup.trackAlbumCreated, Album: "/avengers"},
				{Type: backup.ProgressEventAnalysed, Count: 1, Size: 42},
				{Type: backup.trackCatalogued, Count: 1, Size: 42},
				{Type: backup.trackUploaded, Count: 1, Size: 42, Album: "/avengers", MediaType: backup.MediaTypeImage},
			},
			setMocks: func(complete *mocks.TrackScanComplete, analysed *mocks.TrackAnalysed, uploaded *mocks.TrackUploaded) {
				complete.On("OnScanComplete", backup.NewMediaCounter(1, 42)).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(1, 42), backup.NewMediaCounter(1, 42), noExtra).Once().Return()
				uploaded.On("OnUploaded", backup.NewMediaCounter(1, 42), backup.NewMediaCounter(1, 42)).Once().Return()
			},
			wantCountPerAlbum: map[string]*backup.AlbumReport{
				"/avengers": backup.NewTypeCounter(backup.MediaTypeImage, 1, 42, true),
			},
			wantNewAlbums: []string{"/avengers"},
		},
		{
			name: "it should call OnAnalysed and OnUploaded without totals when ScanComplete haven't been received",
			events: []*backup.progressEvent{
				{Type: backup.ProgressEventAnalysed, Count: 1, Size: 1},
				{Type: backup.trackCatalogued, Count: 1, Size: 1},
				{Type: backup.trackUploaded, Count: 1, Size: 1, Album: "/avengers", MediaType: backup.MediaTypeImage},
				{Type: backup.trackScanComplete, Count: 1, Size: 1},
			},
			setMocks: func(complete *mocks.TrackScanComplete, analysed *mocks.TrackAnalysed, uploaded *mocks.TrackUploaded) {
				analysed.On("OnAnalysed", backup.NewMediaCounter(1, 1), backup.MediaCounterZero, noExtra).Once().Return()
				uploaded.On("OnUploaded", backup.NewMediaCounter(1, 1), backup.MediaCounterZero).Once().Return()
				complete.On("OnScanComplete", backup.NewMediaCounter(1, 1)).Once().Return()
			},
			wantCountPerAlbum: map[string]*backup.AlbumReport{
				"/avengers": backup.NewTypeCounter(backup.MediaTypeImage, 1, 1, false),
			},
			wantNewAlbums: nil,
		},
		{
			name: "it should call OnUploaded without totals when Analysed hasn't been received",
			events: []*backup.progressEvent{
				{Type: backup.trackScanComplete, Count: 1, Size: 1},
				{Type: backup.ProgressEventAnalysed, Count: 1, Size: 1},
				{Type: backup.trackUploaded, Count: 1, Size: 1, Album: "/avengers", MediaType: backup.MediaTypeImage},
				{Type: backup.trackCatalogued, Count: 1, Size: 1},
			},
			setMocks: func(complete *mocks.TrackScanComplete, analysed *mocks.TrackAnalysed, uploaded *mocks.TrackUploaded) {
				complete.On("OnScanComplete", backup.NewMediaCounter(1, 1)).Once().Return()
				uploaded.On("OnUploaded", backup.NewMediaCounter(1, 1), backup.MediaCounterZero).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(1, 1), backup.NewMediaCounter(1, 1), noExtra).Once().Return()
			},
			wantCountPerAlbum: map[string]*backup.AlbumReport{
				"/avengers": backup.NewTypeCounter(backup.MediaTypeImage, 1, 1, false),
			},
			wantNewAlbums: nil,
		},
		{
			name: "it should remove filtered out medias from upload totals",
			events: []*backup.progressEvent{
				{Type: backup.trackScanComplete, Count: 4, Size: 1111},
				{Type: backup.ProgressEventAnalysed, Count: 1, Size: 1},
				{Type: backup.trackCatalogued, Count: 1, Size: 1},
				{Type: backup.trackAlreadyExistsInCatalog, Count: 1, Size: 10},
				{Type: backup.trackDuplicatedInVolume, Count: 1, Size: 100},
				{Type: backup.trackWrongAlbum, Count: 1, Size: 1000},
				{Type: backup.trackUploaded, Count: 1, Size: 1, Album: "/avengers", MediaType: backup.MediaTypeImage},
			},
			setMocks: func(complete *mocks.TrackScanComplete, analysed *mocks.TrackAnalysed, uploaded *mocks.TrackUploaded) {
				complete.On("OnScanComplete", backup.NewMediaCounter(4, 1111)).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(1, 1), backup.NewMediaCounter(4, 1111), noExtra).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(2, 11), backup.NewMediaCounter(4, 1111), noExtra).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(3, 111), backup.NewMediaCounter(4, 1111), noExtra).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(4, 1111), backup.NewMediaCounter(4, 1111), noExtra).Once().Return()
				uploaded.On("OnUploaded", backup.NewMediaCounter(1, 1), backup.NewMediaCounter(1, 1)).Once().Return()
			},
			wantCountPerAlbum: map[string]*backup.AlbumReport{
				"/avengers": backup.NewTypeCounter(backup.MediaTypeImage, 1, 1, false),
			},
			wantNewAlbums: nil,
		},
		{
			name: "it should count the media rejected by the analysis",
			events: []*backup.progressEvent{
				{Type: backup.trackScanComplete, Count: 2, Size: 11},
				{Type: backup.ProgressEventAnalysed, Count: 1, Size: 1},
				{Type: backup.trackCatalogued, Count: 1, Size: 1},
				{Type: backup.trackAnalysisFailed, Count: 1, Size: 10},
				{Type: backup.trackUploaded, Count: 1, Size: 1, Album: "/avengers", MediaType: backup.MediaTypeImage},
			},
			setMocks: func(complete *mocks.TrackScanComplete, analysed *mocks.TrackAnalysed, uploaded *mocks.TrackUploaded) {
				complete.On("OnScanComplete", backup.NewMediaCounter(2, 11)).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(1, 1), backup.NewMediaCounter(2, 11), noExtra).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(2, 11), backup.NewMediaCounter(2, 11), backup.ExtraCounts{
					Cached:   backup.MediaCounterZero,
					Rejected: backup.NewMediaCounter(1, 10),
				}).Once().Return()
				uploaded.On("OnUploaded", backup.NewMediaCounter(1, 1), backup.NewMediaCounter(1, 1)).Once().Return()
			},
			wantCountPerAlbum: map[string]*backup.AlbumReport{
				"/avengers": backup.NewTypeCounter(backup.MediaTypeImage, 1, 1, false),
			},
			wantNewAlbums: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trackScanComplete := mocks.NewTrackScanComplete(t)
			trackAnalysed := mocks.NewTrackAnalysed(t)
			trackUploaded := mocks.NewTrackUploaded(t)

			tt.setMocks(trackScanComplete, trackAnalysed, trackUploaded)

			channel := make(chan *backup.progressEvent, len(tt.events))
			tracker := backup.NewTracker(channel, trackScanComplete, trackAnalysed, trackUploaded)

			for _, event := range tt.events {
				channel <- event
			}
			close(channel)

			<-tracker.Done

			assert.Equal(t, tt.wantCountPerAlbum, tracker.CountPerAlbum(), tt.name)
			assert.Equal(t, tt.wantNewAlbums, tracker.NewAlbums(), tt.name)
		})
	}
}
