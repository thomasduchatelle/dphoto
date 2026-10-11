package catalog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func TestRenameAlbum_RenameAlbum(t *testing.T) {
	const owner = "ironman"
	may26 := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)
	jun26 := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	newName := "Avenger 1"

	existingAlbum := &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      ownermodel.Owner(owner),
			FolderName: catalog.NewFolderName("/may-26"),
		},
		Name:  "may 26",
		Start: may26,
		End:   jun26,
	}
	generatedRenamedId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/2026-05_Avenger_1")}
	forcedRenamedId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("Avengers_vs_Loki")}

	photo5may26 := photoAt("photo5may26", may26.AddDate(0, 0, 5))
	photo15may26 := photoAt("photo15may26", may26.AddDate(0, 0, 15))

	coversOnMay26 := func() *CoverRepositoryInMemory {
		return NewCoverRepositoryInMemory(coversFor(existingAlbum.AlbumId, pickedCover(photo5may26), randomCover(photo15may26)))
	}

	type fields struct {
		Catalog *CatalogInMemory
		Covers  *CoverRepositoryInMemory
	}
	type args struct {
		request catalog.RenameAlbumRequest
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		expectAlbumsByIds   map[catalog.AlbumId]string
		expectMediasByAlbum map[catalog.AlbumId][]*catalog.MediaMeta
		expectRenamedEvents []catalog.AlbumRenamed
		expectCoversByAlbum map[catalog.AlbumId][]catalog.Cover
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name: "it should get an error if the new name is empty",
			fields: fields{
				Catalog: NewCatalogInMemory(withAlbum(existingAlbum, photo5may26)),
				Covers:  coversOnMay26(),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:    existingAlbum.AlbumId,
					NewName:      "",
					RenameFolder: false,
				},
			},
			expectAlbumsByIds:   map[catalog.AlbumId]string{existingAlbum.AlbumId: existingAlbum.Name},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{existingAlbum.AlbumId: {photo5may26}},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{existingAlbum.AlbumId: {pickedCover(photo5may26), randomCover(photo15may26)}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNameMandatoryErr)
			},
		},
		{
			name: "it should get an error if the album doesn't exist",
			fields: fields{
				Catalog: NewCatalogInMemory(),
				Covers:  NewCoverRepositoryInMemory(),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:    existingAlbum.AlbumId,
					NewName:      newName,
					RenameFolder: true,
				},
			},
			expectAlbumsByIds:   map[catalog.AlbumId]string{},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNotFoundErr)
			},
		},
		{
			name: "it should update the name in place (folder unchanged, no transfer, covers untouched)",
			fields: fields{
				Catalog: NewCatalogInMemory(withAlbum(existingAlbum, photo5may26, photo15may26)),
				Covers:  coversOnMay26(),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:    existingAlbum.AlbumId,
					NewName:      newName,
					RenameFolder: false,
				},
			},
			expectAlbumsByIds:   map[catalog.AlbumId]string{existingAlbum.AlbumId: newName},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{existingAlbum.AlbumId: {photo5may26, photo15may26}},
			expectRenamedEvents: []catalog.AlbumRenamed{{
				ExistingAlbum: *existingAlbum,
				RenamedAlbum: catalog.Album{
					AlbumId: existingAlbum.AlbumId,
					Name:    newName,
					Start:   existingAlbum.Start,
					End:     existingAlbum.End,
				},
			}},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{existingAlbum.AlbumId: {pickedCover(photo5may26), randomCover(photo15may26)}},
			wantErr:             assert.NoError,
		},
		{
			name: "it should replace the album with a new folder name generated from the new name, moving medias and covers to the new identity",
			fields: fields{
				Catalog: NewCatalogInMemory(withAlbum(existingAlbum, photo5may26, photo15may26)),
				Covers:  coversOnMay26(),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:    existingAlbum.AlbumId,
					NewName:      newName,
					RenameFolder: true,
				},
			},
			expectAlbumsByIds:   map[catalog.AlbumId]string{generatedRenamedId: newName},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{generatedRenamedId: {photo5may26, photo15may26}},
			expectRenamedEvents: []catalog.AlbumRenamed{{
				ExistingAlbum: *existingAlbum,
				RenamedAlbum: catalog.Album{
					AlbumId: generatedRenamedId,
					Name:    newName,
					Start:   may26,
					End:     jun26,
				},
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{generatedRenamedId: {photo5may26.Id, photo15may26.Id}},
					FromAlbums: []catalog.AlbumId{existingAlbum.AlbumId},
				},
			}},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{generatedRenamedId: {pickedCover(photo5may26), randomCover(photo15may26)}},
			wantErr:             assert.NoError,
		},
		{
			name: "it should replace the album with a forced folder name, moving medias and covers to the forced identity",
			fields: fields{
				Catalog: NewCatalogInMemory(withAlbum(existingAlbum, photo5may26, photo15may26)),
				Covers:  coversOnMay26(),
			},
			args: args{
				request: catalog.RenameAlbumRequest{
					CurrentId:        existingAlbum.AlbumId,
					NewName:          newName,
					RenameFolder:     false,
					ForcedFolderName: "Avengers_vs_Loki",
				},
			},
			expectAlbumsByIds:   map[catalog.AlbumId]string{forcedRenamedId: newName},
			expectMediasByAlbum: map[catalog.AlbumId][]*catalog.MediaMeta{forcedRenamedId: {photo5may26, photo15may26}},
			expectRenamedEvents: []catalog.AlbumRenamed{{
				ExistingAlbum: *existingAlbum,
				RenamedAlbum: catalog.Album{
					AlbumId: forcedRenamedId,
					Name:    newName,
					Start:   may26,
					End:     jun26,
				},
				TransferredMedias: catalog.TransferredMedias{
					Transfers:  map[catalog.AlbumId][]catalog.MediaId{forcedRenamedId: {photo5may26.Id, photo15may26.Id}},
					FromAlbums: []catalog.AlbumId{existingAlbum.AlbumId},
				},
			}},
			expectCoversByAlbum: map[catalog.AlbumId][]catalog.Cover{forcedRenamedId: {pickedCover(photo5may26), randomCover(photo15may26)}},
			wantErr:             assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := &AlbumRenamedObserverInMemory{}
			renameAlbum := catalog.NewRenameAlbum(
				tt.fields.Catalog,
				&catalog.TransferMediasFromRepository{TransferMediasRepository: tt.fields.Catalog},
				tt.fields.Covers,
				observer,
			)

			err := renameAlbum.RenameAlbum(context.Background(), tt.args.request)
			if !tt.wantErr(t, err, fmt.Sprintf("RenameAlbum(%v)", tt.args.request)) {
				return
			}

			names := make(map[catalog.AlbumId]string, len(tt.fields.Catalog.AlbumsByIds()))
			for id, album := range tt.fields.Catalog.AlbumsByIds() {
				names[id] = album.Name
			}
			assert.Equal(t, tt.expectAlbumsByIds, names, "album names in the catalog")
			assert.Equal(t, tt.expectMediasByAlbum, tt.fields.Catalog.MediasByAlbum(), "medias per album")
			assert.Equal(t, tt.expectRenamedEvents, observer.Events, "AlbumRenamed events fired")
			assert.Equal(t, tt.expectCoversByAlbum, tt.fields.Covers.Covers, "covers in the repository")
		})
	}
}
