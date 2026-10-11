package catalog_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func TestTimelineAggregate_RemoveAlbum(t *testing.T) {
	const owner = "ironman"
	jan24 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mar24 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	apr24 := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	may24 := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	jul24 := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	nov24 := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
	dec24 := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	toDeleteId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers-1")}
	toDeleteAlbum := catalog.Album{AlbumId: toDeleteId, Name: "Avenger 1", Start: mar24, End: may24}

	isolatedAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/isolated")},
		Name:    "Isolated", Start: nov24, End: dec24,
	}
	allYearAlbum := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/lifetime")},
		Name:    "lifetime", Start: jan24, End: jan25,
	}
	q1Album := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/q1")},
		Name:    "q1", Start: jan24, End: apr24,
	}
	q2Album := catalog.Album{
		AlbumId: catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/q2")},
		Name:    "q2", Start: apr24, End: jul24,
	}

	type args struct {
		albums  []catalog.Album
		albumId catalog.AlbumId
	}
	tests := []struct {
		name          string
		args          args
		wantRecords   catalog.MediaTransferRecords
		wantOrphaned  []catalog.MediaSelector
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "it should hand every media over to the single surrounding album when it covers the deleted range",
			args: args{
				albums:  []catalog.Album{allYearAlbum, toDeleteAlbum},
				albumId: toDeleteId,
			},
			wantRecords: catalog.MediaTransferRecords{
				allYearAlbum.AlbumId: {{FromAlbums: []catalog.AlbumId{toDeleteId}, Start: mar24, End: may24}},
			},
			wantOrphaned: []catalog.MediaSelector{},
			wantErr:      assert.NoError,
		},
		{
			name: "it should split medias between consecutive surrounding albums",
			args: args{
				albums:  []catalog.Album{q1Album, q2Album, toDeleteAlbum},
				albumId: toDeleteId,
			},
			wantRecords: catalog.MediaTransferRecords{
				q1Album.AlbumId: {{FromAlbums: []catalog.AlbumId{toDeleteId}, Start: mar24, End: apr24}},
				q2Album.AlbumId: {{FromAlbums: []catalog.AlbumId{toDeleteId}, Start: apr24, End: may24}},
			},
			wantOrphaned: []catalog.MediaSelector{},
			wantErr:      assert.NoError,
		},
		{
			name: "it should report the uncovered part as orphaned when no surrounding album covers it",
			args: args{
				albums:  []catalog.Album{q1Album, toDeleteAlbum},
				albumId: toDeleteId,
			},
			wantRecords: catalog.MediaTransferRecords{
				q1Album.AlbumId: {{FromAlbums: []catalog.AlbumId{toDeleteId}, Start: mar24, End: apr24}},
			},
			wantOrphaned: []catalog.MediaSelector{
				{FromAlbums: []catalog.AlbumId{toDeleteId}, Start: apr24, End: may24},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should produce no transfer and one orphaned selector when the album stands alone",
			args: args{
				albums:  []catalog.Album{isolatedAlbum, toDeleteAlbum},
				albumId: toDeleteId,
			},
			wantRecords: catalog.MediaTransferRecords{},
			wantOrphaned: []catalog.MediaSelector{
				{FromAlbums: []catalog.AlbumId{toDeleteId}, Start: mar24, End: may24},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return AlbumNotFoundErr when the album does not exist in the timeline",
			args: args{
				albums:  []catalog.Album{allYearAlbum},
				albumId: toDeleteId,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr, i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptrs := make([]*catalog.Album, 0, len(tt.args.albums))
			for i := range tt.args.albums {
				ptrs = append(ptrs, &tt.args.albums[i])
			}
			aggregate, err := catalog.NewTimelineAggregate(ptrs)
			if !assert.NoError(t, err) {
				return
			}

			records, orphaned, err := aggregate.RemoveAlbum(tt.args.albumId)
			if !tt.wantErr(t, err, fmt.Sprintf("RemoveAlbum(%v)", tt.args.albumId)) {
				return
			}
			assert.Equal(t, tt.wantRecords, records, "transfer records")
			assert.Equal(t, tt.wantOrphaned, orphaned, "orphaned selectors")
		})
	}
}
