package backup_test

import (
	"github.com/stretchr/testify/assert"
	mocks2 "github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"testing"
)

func TestTracker(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name              string
		events            []*backup.ProgressEvent
		setMocks          func(*mocks2.TrackScanComplete, *mocks2.TrackAnalysed, *mocks2.TrackUploaded)
		wantCountPerAlbum map[string]*backup.TypeCounter
		wantNewAlbums     []string
	}{
		{
			name: "it should call OnAnalysed and OnUploaded if everything follow the happy-order",
			events: []*backup.ProgressEvent{
				{Type: backup.ProgressEventScanComplete, Count: 1, Size: 42},
				{Type: backup.ProgressEventAlbumCreated, Album: "/avengers"},
				{Type: backup.ProgressEventAnalysed, Count: 1, Size: 42},
				{Type: backup.ProgressEventReadyForUpload, Count: 1, Size: 42},
				{Type: backup.ProgressEventUploaded, Count: 1, Size: 42, Album: "/avengers", MediaType: backup.MediaTypeImage},
			},
			setMocks: func(complete *mocks2.TrackScanComplete, analysed *mocks2.TrackAnalysed, uploaded *mocks2.TrackUploaded) {
				complete.On("OnScanComplete", backup.NewMediaCounter(1, 42)).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(1, 42), backup.NewMediaCounter(1, 42)).Once().Return()
				uploaded.On("OnUploaded", backup.NewMediaCounter(1, 42), backup.NewMediaCounter(1, 42)).Once().Return()
			},
			wantCountPerAlbum: map[string]*backup.TypeCounter{
				"/avengers": backup.NewTypeCounter(backup.MediaTypeImage, 1, 42),
			},
			wantNewAlbums: []string{"/avengers"},
		},
		{
			name: "it should call OnAnalysed and OnUploaded without totals when ScanComplete haven't been received",
			events: []*backup.ProgressEvent{
				{Type: backup.ProgressEventAnalysed, Count: 1, Size: 1},
				{Type: backup.ProgressEventReadyForUpload, Count: 1, Size: 1},
				{Type: backup.ProgressEventUploaded, Count: 1, Size: 1, Album: "/avengers", MediaType: backup.MediaTypeImage},
				{Type: backup.ProgressEventScanComplete, Count: 1, Size: 1},
			},
			setMocks: func(complete *mocks2.TrackScanComplete, analysed *mocks2.TrackAnalysed, uploaded *mocks2.TrackUploaded) {
				analysed.On("OnAnalysed", backup.NewMediaCounter(1, 1), backup.MediaCounterZero).Once().Return()
				uploaded.On("OnUploaded", backup.NewMediaCounter(1, 1), backup.MediaCounterZero).Once().Return()
				complete.On("OnScanComplete", backup.NewMediaCounter(1, 1)).Once().Return()
			},
			wantCountPerAlbum: map[string]*backup.TypeCounter{
				"/avengers": backup.NewTypeCounter(backup.MediaTypeImage, 1, 1),
			},
			wantNewAlbums: nil,
		},
		{
			name: "it should call OnUploaded without totals when Analysed hasn't been received",
			events: []*backup.ProgressEvent{
				{Type: backup.ProgressEventScanComplete, Count: 1, Size: 1},
				{Type: backup.ProgressEventAnalysed, Count: 1, Size: 1},
				{Type: backup.ProgressEventUploaded, Count: 1, Size: 1, Album: "/avengers", MediaType: backup.MediaTypeImage},
				{Type: backup.ProgressEventReadyForUpload, Count: 1, Size: 1},
			},
			setMocks: func(complete *mocks2.TrackScanComplete, analysed *mocks2.TrackAnalysed, uploaded *mocks2.TrackUploaded) {
				complete.On("OnScanComplete", backup.NewMediaCounter(1, 1)).Once().Return()
				uploaded.On("OnUploaded", backup.NewMediaCounter(1, 1), backup.MediaCounterZero).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(1, 1), backup.NewMediaCounter(1, 1)).Once().Return()
			},
			wantCountPerAlbum: map[string]*backup.TypeCounter{
				"/avengers": backup.NewTypeCounter(backup.MediaTypeImage, 1, 1),
			},
			wantNewAlbums: nil,
		},
		{
			name: "it should remove filtered out medias from upload totals",
			events: []*backup.ProgressEvent{
				{Type: backup.ProgressEventScanComplete, Count: 4, Size: 1111},
				{Type: backup.ProgressEventAnalysed, Count: 1, Size: 1},
				{Type: backup.ProgressEventReadyForUpload, Count: 1, Size: 1},
				{Type: backup.ProgressEventAlreadyExists, Count: 1, Size: 10},
				{Type: backup.ProgressEventDuplicate, Count: 1, Size: 100},
				{Type: backup.ProgressEventWrongAlbum, Count: 1, Size: 1000},
				{Type: backup.ProgressEventUploaded, Count: 1, Size: 1, Album: "/avengers", MediaType: backup.MediaTypeImage},
			},
			setMocks: func(complete *mocks2.TrackScanComplete, analysed *mocks2.TrackAnalysed, uploaded *mocks2.TrackUploaded) {
				complete.On("OnScanComplete", backup.NewMediaCounter(4, 1111)).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(1, 1), backup.NewMediaCounter(4, 1111)).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(2, 11), backup.NewMediaCounter(4, 1111)).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(3, 111), backup.NewMediaCounter(4, 1111)).Once().Return()
				analysed.On("OnAnalysed", backup.NewMediaCounter(4, 1111), backup.NewMediaCounter(4, 1111)).Once().Return()
				uploaded.On("OnUploaded", backup.NewMediaCounter(1, 1), backup.NewMediaCounter(1, 1)).Once().Return()
			},
			wantCountPerAlbum: map[string]*backup.TypeCounter{
				"/avengers": backup.NewTypeCounter(backup.MediaTypeImage, 1, 1),
			},
			wantNewAlbums: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trackScanComplete := mocks2.NewTrackScanComplete(t)
			trackAnalysed := mocks2.NewTrackAnalysed(t)
			trackUploaded := mocks2.NewTrackUploaded(t)

			tt.setMocks(trackScanComplete, trackAnalysed, trackUploaded)

			channel := make(chan *backup.ProgressEvent, len(tt.events))
			tracker := backup.NewTracker(channel, trackScanComplete, trackAnalysed, trackUploaded)

			for _, event := range tt.events {
				channel <- event
			}
			close(channel)

			<-tracker.Done

			a.Equal(tt.wantCountPerAlbum, tracker.CountPerAlbum(), tt.name)
			a.Equal(tt.wantNewAlbums, tracker.NewAlbums(), tt.name)
		})
	}
}
