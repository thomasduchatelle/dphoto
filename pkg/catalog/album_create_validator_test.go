package catalog_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"testing"
	"time"
)

func TestCreateAlbumValidator_Create(t *testing.T) {
	const owner = "tonystark"
	album := catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      owner,
			FolderName: catalog.FolderName("/2024-04_Ironman_1"),
		},
		Name:  "Ironman 1",
		Start: time.Date(2024, 04, 28, 8, 33, 42, 0, time.UTC),
		End:   time.Date(2024, 05, 1, 0, 0, 0, 0, time.UTC),
	}
	standardRequest := catalog.CreateAlbumRequest{
		Owner: owner,
		Name:  album.Name,
		Start: album.Start,
		End:   album.End,
	}

	type fields struct {
		Observer func(t *testing.T) catalog.CreateAlbumObserver
	}
	type args struct {
		request catalog.CreateAlbumRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should NOT create the album without owner",
			fields: fields{
				Observer: expectCreateAlbumObserveNotCalled(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: "",
					Name:  "foobar",
					Start: album.Start,
					End:   album.End,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ownermodel.EmptyOwnerError)
			},
		},
		{
			name: "it should NOT create the album without name",
			fields: fields{
				Observer: expectCreateAlbumObserveNotCalled(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "",
					Start: album.Start,
					End:   album.End,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumNameMandatoryErr)
			},
		},
		{
			name: "it should NOT create the album without start date",
			fields: fields{
				Observer: expectCreateAlbumObserveNotCalled(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "foobar",
					End:   album.End,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumStartAndEndDateMandatoryErr)
			},
		},
		{
			name: "it should NOT create the album without end date",
			fields: fields{
				Observer: expectCreateAlbumObserveNotCalled(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "foobar",
					Start: album.Start,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumStartAndEndDateMandatoryErr)
			},
		},
		{
			name: "it should NOT create the album with start and end reversed",
			fields: fields{
				Observer: expectCreateAlbumObserveNotCalled(),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner: owner,
					Name:  "foobar",
					Start: album.End,
					End:   album.Start,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, catalog.AlbumEndDateMustBeAfterStartErr)
			},
		},
		{
			name: "it should create the album with a generated name",
			fields: fields{
				Observer: expectCreateAlbumObserved(album),
			},
			args: args{
				request: standardRequest,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should create the album with a forced name",
			fields: fields{
				Observer: expectCreateAlbumObserved(catalog.Album{
					AlbumId: catalog.AlbumId{
						Owner:      owner,
						FolderName: "/Phase_1_Avenger",
					},
					Name:  "Avenger 1",
					Start: album.Start,
					End:   album.End,
				}),
			},
			args: args{
				request: catalog.CreateAlbumRequest{
					Owner:            owner,
					Name:             "Avenger 1",
					Start:            album.Start,
					End:              album.End,
					ForcedFolderName: "Phase_1_Avenger",
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var observers []catalog.CreateAlbumObserver
			if tt.fields.Observer != nil {
				observers = append(observers, tt.fields.Observer(t))
			}
			c := &catalog.CreateAlbum{
				CreateAlbumValidator: catalog.CreateAlbumValidator{
					Observers: observers,
				},
			}

			_, err := c.Create(context.TODO(), tt.args.request)
			tt.wantErr(t, err, fmt.Sprintf("Create(%v)", tt.args.request))
		})
	}
}