package backup

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"maps"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestNewCataloguerAcceptance(t *testing.T) {
	ctx := context.Background()
	jan24 := time.Date(2020, time.January, 24, 0, 0, 0, 0, time.UTC)
	analysedMedia1 := newAnalysedMedia("file1.jpg", jan24, 12)
	reference1Exists := &CatalogReferenceStub{MediaIdValue: "media1", AlbumFolderNameValue: "album1", ExistsValue: true}
	reference1IsNew := &CatalogReferenceStub{MediaIdValue: "media1", AlbumFolderNameValue: "album1"}

	analysedMedia2 := newAnalysedMedia("file2.jpg", jan24, 13)
	reference2 := &CatalogReferenceStub{MediaIdValue: "media2", AlbumFolderNameValue: "album2"}

	analysedMedia3 := newAnalysedMedia("file3.jpg", jan24, 14)
	reference3 := &CatalogReferenceStub{MediaIdValue: "media3", AlbumFolderNameValue: "album3"}

	type fields struct {
		CatalogReferencer Cataloguer
		options           Options
	}
	type args struct {
		medias []*AnalysedMedia
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         []BackingUpMediaRequest
		wantFiltered map[string]assert.ErrorAssertionFunc
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "it shouldn't do anything if no media",
			fields: fields{
				CatalogReferencer: make(CatalogReferencerFake),
			},
			args: args{
				medias: nil,
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should add the reference to each media",
			fields: fields{
				CatalogReferencer: CatalogReferencerFake{
					analysedMedia1: reference1IsNew,
				},
			},
			args: args{
				medias: []*AnalysedMedia{analysedMedia1},
			},
			want: []BackingUpMediaRequest{
				{
					AnalysedMedia:    analysedMedia1,
					CatalogReference: reference1IsNew,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should filter out exiting media",
			fields: fields{
				CatalogReferencer: CatalogReferencerFake{
					analysedMedia1: reference1Exists,
				},
			},
			args: args{
				medias: []*AnalysedMedia{analysedMedia1},
			},
			want: nil,
			wantFiltered: map[string]assert.ErrorAssertionFunc{
				analysedMedia1.FoundMedia.MediaPath().Filename: errorIs(ErrCatalogerFilterMustNotAlreadyExists),
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should filter out medias that are not in the selected album",
			fields: fields{
				CatalogReferencer: CatalogReferencerFake{
					analysedMedia1: reference1Exists,
					analysedMedia2: reference2,
					analysedMedia3: reference3,
				},
				options: OptionsOnlyAlbums("album1", "album2"),
			},
			args: args{
				medias: []*AnalysedMedia{analysedMedia1, analysedMedia2, analysedMedia3},
			},
			want: []BackingUpMediaRequest{
				{
					AnalysedMedia:    analysedMedia2,
					CatalogReference: reference2,
				},
			},
			wantFiltered: map[string]assert.ErrorAssertionFunc{
				analysedMedia1.FoundMedia.MediaPath().Filename: errorIs(ErrCatalogerFilterMustNotAlreadyExists),
				analysedMedia3.FoundMedia.MediaPath().Filename: errorIs(ErrCatalogerFilterMustBeInAlbum),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observerOk := new(CataloguerObserverFake)
			observerFiltered := new(CataloguerFilterObserverFake)

			cataloger := &cataloguerAggregate{
				cataloguer: tt.fields.CatalogReferencer,
				observerWithFilters: applyFiltersOnCataloguer{
					CatalogReferencerObservers: []CatalogReferencerObserver{observerOk},
					CataloguerFilterObservers:  []CataloguerFilterObserver{observerFiltered},
					CataloguerFilters:          postCataloguerFiltersList(tt.fields.options),
				},
			}

			err := cataloger.OnBatchOfAnalysedMedia(ctx, tt.args.medias)

			if !tt.wantErr(t, err, fmt.Sprintf("Catalog(%v)", tt.args.medias)) {
				return
			}

			assert.Equalf(t, tt.want, observerOk.Got, "Catalog(%v)", tt.args.medias)
			if assert.Equal(t, slices.Sorted(maps.Keys(tt.wantFiltered)), slices.Sorted(maps.Keys(observerFiltered.Got))) {
				for filename, wantErr := range tt.wantFiltered {
					wantErr(t, observerFiltered.Got[filename], fmt.Sprintf("Catalog(%v)", tt.args.medias))
				}
			}
		})
	}
}

func errorIs(expectedErr error) func(t assert.TestingT, err error, i ...interface{}) bool {
	return func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, expectedErr, err, i...)
	}
}

func newAnalysedMedia(filename string, mediaDateTime time.Time, size int) *AnalysedMedia {
	return &AnalysedMedia{
		FoundMedia: NewInMemoryMedia(filename, mediaDateTime, []byte(filename+": "+strings.Repeat("a", size-len(filename)-2))),
		Details:    &MediaDetails{DateTime: mediaDateTime},
	}
}

type CatalogReferenceStub struct {
	MediaIdValue         string
	AlbumFolderNameValue string
	ExistsValue          bool
	AlbumCreatedValue    bool
}

func (c *CatalogReferenceStub) MediaId() string {
	return c.MediaIdValue
}

func (c *CatalogReferenceStub) AlbumCreated() bool {
	return c.AlbumCreatedValue
}

func (c *CatalogReferenceStub) Exists() bool {
	return c.ExistsValue
}

func (c *CatalogReferenceStub) UniqueIdentifier() string {
	return c.MediaIdValue
}

func (c *CatalogReferenceStub) AlbumFolderName() string {
	return c.AlbumFolderNameValue
}

type CatalogReferencerFake map[*AnalysedMedia]CatalogReference

func (c CatalogReferencerFake) Reference(ctx context.Context, medias []*AnalysedMedia, observer CatalogReferencerObserver) error {
	var result []BackingUpMediaRequest
	for _, media := range medias {
		for key, reference := range c {
			if key.FoundMedia.String() == media.FoundMedia.String() {
				result = append(result, BackingUpMediaRequest{
					AnalysedMedia:    media,
					CatalogReference: reference,
				})
				break
			}
		}
	}

	if len(medias) != len(result) {
		return fmt.Errorf("[CatalogReferencerFake] missing reference for some media")
	}

	return observer.OnMediaCatalogued(ctx, result)
}

type CataloguerObserverFake struct {
	Got []BackingUpMediaRequest
}

func (c *CataloguerObserverFake) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	c.Got = append(c.Got, requests...)
	return nil
}

type CataloguerFilterObserverFake struct {
	Got map[string]error
}

func (c *CataloguerFilterObserverFake) OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error {
	if c.Got == nil {
		c.Got = make(map[string]error)
	}

	c.Got[media.FoundMedia.MediaPath().Filename] = cause
	return nil
}
