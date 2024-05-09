package catalog_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/internal/mocks"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"testing"
	"time"
)

func TestMediaTransfer_Transfer(t *testing.T) {
	avenger1Id := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/avengers-1")}
	ironman1Id := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/ironman-1")}
	records := catalog.MediaTransferRecords{
		avenger1Id: {
			{
				FromAlbums: []catalog.AlbumId{ironman1Id},
				Start:      time.Time{},
				End:        time.Time{},
			},
		},
	}
	transfers := catalog.TransferredMedias{
		avenger1Id: []catalog.MediaId{"media-1", "media-2"},
	}
	emptyTransfers := catalog.TransferredMedias{
		avenger1Id: []catalog.MediaId{},
		ironman1Id: nil,
	}

	type fields struct {
		TransferMedias           func(t *testing.T) catalog.TransferMediasPort
		TimelineMutationObserver func(t *testing.T) catalog.TimelineMutationObserver
	}
	type args struct {
		deletedAlbum catalog.AlbumId
		transfers    catalog.MediaTransferRecords
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should perform the media transfer on the catalog DB",
			fields: fields{
				TransferMedias: expectTransferMediasPortCalled(records, transfers),
			},
			args: args{
				deletedAlbum: ironman1Id,
				transfers:    records,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should notify observers that medias should be transferred",
			fields: fields{
				TransferMedias:           expectTransferMediasPortCalled(records, transfers),
				TimelineMutationObserver: expectTimelineMutationObserverCalled(transfers),
			},
			args: args{
				deletedAlbum: ironman1Id,
				transfers:    records,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not notify observers when no medias should be transferred",
			fields: fields{
				TransferMedias:           expectTransferMediasPortCalled(records, emptyTransfers),
				TimelineMutationObserver: timelineMutationObserverNotCalled(),
			},
			args: args{
				deletedAlbum: ironman1Id,
				transfers:    records,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var observers []catalog.TimelineMutationObserver
			if tt.fields.TimelineMutationObserver != nil {
				observers = append(observers, tt.fields.TimelineMutationObserver(t))
			}
			d := &catalog.DeleteAlbumMediaTransfer{
				MediaTransfer: catalog.MediaTransfer{
					TransferMedias:            tt.fields.TransferMedias(t),
					TimelineMutationObservers: observers,
				},
			}
			err := d.OnDeleteAlbum(context.Background(), tt.args.deletedAlbum, tt.args.transfers)
			tt.wantErr(t, err, fmt.Sprintf("OnDeleteAlbum(%v, %v)", tt.args.deletedAlbum, tt.args.transfers))
		})
	}
}

func timelineMutationObserverNotCalled() func(t *testing.T) catalog.TimelineMutationObserver {
	return func(t *testing.T) catalog.TimelineMutationObserver {
		return mocks.NewTimelineMutationObserver(t)
	}
}

func expectTimelineMutationObserverCalled(transfers catalog.TransferredMedias) func(t *testing.T) catalog.TimelineMutationObserver {
	return func(t *testing.T) catalog.TimelineMutationObserver {
		observer := mocks.NewTimelineMutationObserver(t)
		observer.EXPECT().Observe(mock.Anything, transfers).Return(nil).Once()
		return observer
	}
}

// TODO use expectTransferMediasPortCalled on create album as well
func expectTransferMediasPortCalled(expectedRecords catalog.MediaTransferRecords, returnedTransfers catalog.TransferredMedias) func(t *testing.T) catalog.TransferMediasPort {
	return func(t *testing.T) catalog.TransferMediasPort {
		port := mocks.NewTransferMediasPort(t)
		port.EXPECT().
			TransferMediasFromRecords(mock.Anything, expectedRecords).
			Return(returnedTransfers, nil).
			Once()
		return port
	}
}
