package catalog

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
	"time"
)

func TestMediaTransferExecutor_Transfer(t *testing.T) {
	avenger1Id := AlbumId{Owner: "ironman", FolderName: NewFolderName("/avengers-1")}
	ironman1Id := AlbumId{Owner: "ironman", FolderName: NewFolderName("/ironman-1")}
	records := MediaTransferRecords{
		avenger1Id: {
			{
				FromAlbums: []AlbumId{ironman1Id},
				Start:      time.Time{},
				End:        time.Time{},
			},
		},
	}
	transfersToAvenger1 := TransferredMedias{
		Transfers: map[AlbumId][]MediaId{
			avenger1Id: {"media-1", "media-2"},
		},
	}
	emptyTransfers := TransferredMedias{
		Transfers: map[AlbumId][]MediaId{
			avenger1Id: {},
			ironman1Id: nil,
		},
	}
	recordsSwapped := MediaTransferRecords{
		avenger1Id: {
			{
				FromAlbums: []AlbumId{ironman1Id},
				Start:      time.Time{},
				End:        time.Time{},
			},
		},
		ironman1Id: {
			{
				FromAlbums: []AlbumId{avenger1Id},
				Start:      time.Time{},
				End:        time.Time{},
			},
		},
	}
	transfersSwapped := TransferredMedias{
		Transfers: map[AlbumId][]MediaId{
			avenger1Id: {"media-1"},
			ironman1Id: {"media-2"},
		},
	}

	type fields struct {
		TransferMedias TransferMediasRepositoryPort
	}
	type args struct {
		records MediaTransferRecords
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantObserved []TransferredMedias
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "it should accept an empty records",
			fields: fields{
				TransferMedias: &TransferMediasRepositoryPortFake{
					TransferredMedias: NewTransferredMedias(),
				},
			},
			args: args{
				records: nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should notify observers that medias should be transferred",
			fields: fields{
				TransferMedias: &TransferMediasRepositoryPortFake{
					TransferredMedias: transfersToAvenger1,
				},
			},
			args: args{
				records: records,
			},
			wantObserved: []TransferredMedias{{
				Transfers:  transfersToAvenger1.Transfers,
				FromAlbums: []AlbumId{ironman1Id},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should notify observers that media have been swapped (without any source albums)",
			fields: fields{
				TransferMedias: &TransferMediasRepositoryPortFake{
					TransferredMedias: transfersSwapped,
				},
			},
			args: args{
				records: recordsSwapped,
			},
			wantObserved: []TransferredMedias{{
				Transfers:  transfersSwapped.Transfers,
				FromAlbums: nil,
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should not notify observers when no medias should be transferred",
			fields: fields{
				TransferMedias: &TransferMediasRepositoryPortFake{
					TransferredMedias: emptyTransfers,
				},
			},
			args: args{
				records: records,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := new(TimelineMutationObserverFake)
			d := &MediaTransferExecutor{
				TransferMediasRepository:  tt.fields.TransferMedias,
				TimelineMutationObservers: []TimelineMutationObserver{observer},
			}

			err := d.Transfer(context.Background(), tt.args.records)
			if tt.wantErr(t, err, fmt.Sprintf("Transfer(%v)", tt.args.records)) {
				assert.Equal(t, tt.wantObserved, observer.Observed)
			}
		})
	}
}

type TransferMediasRepositoryPortFake struct {
	GotSelectors      MediaTransferRecords
	TransferredMedias TransferredMedias
}

func (t *TransferMediasRepositoryPortFake) TransferMediasFromRecords(ctx context.Context, records MediaTransferRecords) (TransferredMedias, error) {
	t.GotSelectors = records

	transfer := NewTransferredMedias()
	for albumId := range records {
		if ids, ok := t.TransferredMedias.Transfers[albumId]; ok {
			transfer.Transfers[albumId] = ids
		}
		if slices.Contains(t.TransferredMedias.FromAlbums, albumId) {
			transfer.FromAlbums = append(transfer.FromAlbums, albumId)
		}
	}

	return transfer, nil
}

type TimelineMutationObserverFake struct {
	Observed []TransferredMedias
}

func (t *TimelineMutationObserverFake) OnTransferredMedias(ctx context.Context, transfers TransferredMedias) error {
	t.Observed = append(t.Observed, transfers)
	return nil
}
