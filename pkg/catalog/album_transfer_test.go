package catalog

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransferMediasFromRepository_TransferMedias(t *testing.T) {
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
	transfersToAvenger1 := map[AlbumId][]MediaId{
		avenger1Id: {"media-1", "media-2"},
	}
	emptyTransfers := map[AlbumId][]MediaId{
		avenger1Id: {},
		ironman1Id: nil,
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
	transfersSwapped := map[AlbumId][]MediaId{
		avenger1Id: {"media-1"},
		ironman1Id: {"media-2"},
	}

	type fields struct {
		TransferMediasRepository TransferMediasRepositoryPort
	}
	type args struct {
		records MediaTransferRecords
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    TransferredMedias
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return an empty TransferredMedias when the repository moved nothing",
			fields: fields{
				TransferMediasRepository: &TransferMediasRepositoryPortFake{
					Transferred: map[AlbumId][]MediaId{},
				},
			},
			args:    args{records: nil},
			want:    TransferredMedias{Transfers: map[AlbumId][]MediaId{}},
			wantErr: assert.NoError,
		},
		{
			name: "it should derive FromAlbums from the input records when medias were moved",
			fields: fields{
				TransferMediasRepository: &TransferMediasRepositoryPortFake{
					Transferred: transfersToAvenger1,
				},
			},
			args: args{records: records},
			want: TransferredMedias{
				Transfers:  transfersToAvenger1,
				FromAlbums: []AlbumId{ironman1Id},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should not add a source album that is also a destination (swap case)",
			fields: fields{
				TransferMediasRepository: &TransferMediasRepositoryPortFake{
					Transferred: transfersSwapped,
				},
			},
			args: args{records: recordsSwapped},
			want: TransferredMedias{
				Transfers:  transfersSwapped,
				FromAlbums: nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return an empty TransferredMedias without deriving FromAlbums when no media was actually moved",
			fields: fields{
				TransferMediasRepository: &TransferMediasRepositoryPortFake{
					Transferred: emptyTransfers,
				},
			},
			args: args{records: records},
			want: TransferredMedias{
				Transfers: map[AlbumId][]MediaId{
					avenger1Id: {},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &TransferMediasFromRepository{
				TransferMediasRepository: tt.fields.TransferMediasRepository,
			}

			got, err := service.TransferMedias(context.Background(), tt.args.records)
			if tt.wantErr(t, err, fmt.Sprintf("TransferMedias(%v)", tt.args.records)) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

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
	transfersToAvenger1 := map[AlbumId][]MediaId{
		avenger1Id: {"media-1", "media-2"},
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
					Transferred: map[AlbumId][]MediaId{},
				},
			},
			args:    args{records: nil},
			wantErr: assert.NoError,
		},
		{
			name: "it should notify observers that medias should be transferred",
			fields: fields{
				TransferMedias: &TransferMediasRepositoryPortFake{
					Transferred: transfersToAvenger1,
				},
			},
			args: args{records: records},
			wantObserved: []TransferredMedias{{
				Transfers:  transfersToAvenger1,
				FromAlbums: []AlbumId{ironman1Id},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should not notify observers when no medias should be transferred",
			fields: fields{
				TransferMedias: &TransferMediasRepositoryPortFake{
					Transferred: map[AlbumId][]MediaId{avenger1Id: {}, ironman1Id: nil},
				},
			},
			args:    args{records: records},
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

// TransferMediasRepositoryPortFake implements TransferMediasRepositoryPort: it captures the
// records passed in and returns the intersection of the requested destinations with the
// pre-set Transferred map (and only the origins that also appear in that map).
type TransferMediasRepositoryPortFake struct {
	GotSelectors MediaTransferRecords
	Transferred  map[AlbumId][]MediaId
}

func (t *TransferMediasRepositoryPortFake) TransferMediasFromRecords(_ context.Context, records MediaTransferRecords) (map[AlbumId][]MediaId, error) {
	t.GotSelectors = records

	transfer := make(map[AlbumId][]MediaId)
	for albumId := range records {
		if ids, ok := t.Transferred[albumId]; ok {
			transfer[albumId] = ids
		}
	}

	return transfer, nil
}

type TimelineMutationObserverFake struct {
	Observed []TransferredMedias
}

func (t *TimelineMutationObserverFake) OnTransferredMedias(_ context.Context, transfers TransferredMedias) error {
	t.Observed = append(t.Observed, transfers)
	return nil
}
