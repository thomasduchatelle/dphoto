package backup

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"strings"
	"testing"
	"time"
)

func TestNewCatalogerAcceptance(t *testing.T) {
	ctx := context.Background()
	owner := ownermodel.Owner("owner1")
	jan24 := time.Date(2020, time.January, 24, 0, 0, 0, 0, time.UTC)
	analysedMedia1 := newAnalysedMedia("file1.jpg", jan24, 12)
	reference1Exists := &CatalogReferenceStub{MediaIdValue: "media1", AlbumFolderNameValue: "album1", ExistsValue: true}

	analysedMedia2 := newAnalysedMedia("file2.jpg", jan24, 13)
	reference2 := &CatalogReferenceStub{MediaIdValue: "media2", AlbumFolderNameValue: "album2"}

	analysedMedia3 := newAnalysedMedia("file3.jpg", jan24, 14)
	reference3 := &CatalogReferenceStub{MediaIdValue: "media3", AlbumFolderNameValue: "album3"}

	type fields struct {
		ReferencerFactory ReferencerFactory
	}
	type newArgs struct {
		owner   ownermodel.Owner
		options Options
	}
	type args struct {
		medias []*AnalysedMedia
	}
	tests := []struct {
		name       string
		fields     fields
		newArgs    newArgs
		args       args
		want       []*BackingUpMediaRequest
		wantEvents []*ProgressEvent
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:    "it should filter out exiting media",
			newArgs: newArgs{owner: owner, options: Options{}},
			fields: fields{
				ReferencerFactory: &ReferencerFactoryFake{
					CreatorReferencer: CatalogReferencerStub{
						analysedMedia1: reference1Exists,
					},
				},
			},
			args: args{
				medias: []*AnalysedMedia{analysedMedia1},
			},
			want: nil,
			wantEvents: []*ProgressEvent{
				{Type: ProgressEventAlreadyExists, Count: 1, Size: analysedMedia1.FoundMedia.Size()},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should filter out medias that are not in the selected album",
			newArgs: newArgs{owner: owner, options: Options{RestrictedAlbumFolderName: map[string]interface{}{
				"album1": nil,
				"album2": nil,
			}}},
			fields: fields{
				ReferencerFactory: &ReferencerFactoryFake{
					CreatorReferencer: CatalogReferencerStub{
						analysedMedia1: reference1Exists,
						analysedMedia2: reference2,
						analysedMedia3: reference3,
					},
				},
			},
			args: args{
				medias: []*AnalysedMedia{analysedMedia1, analysedMedia2, analysedMedia3},
			},
			want: []*BackingUpMediaRequest{
				{
					AnalysedMedia: analysedMedia2,
					Id:            reference2.MediaIdValue,
					FolderName:    reference2.AlbumFolderNameValue,
				},
			},
			wantEvents: []*ProgressEvent{
				{Type: ProgressEventAlreadyExists, Count: 1, Size: analysedMedia1.FoundMedia.Size()},
				{Type: ProgressEventCatalogued, Count: 1, Size: analysedMedia2.FoundMedia.Size()},
				{Type: ProgressEventWrongAlbum, Count: 1, Size: analysedMedia3.FoundMedia.Size()},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "it should pick a dry-run referencer when in read-only mode",
			newArgs: newArgs{owner: owner, options: Options{DryRun: true}},
			fields: fields{
				ReferencerFactory: &ReferencerFactoryFake{
					DryRunReferencer: CatalogReferencerStub{
						analysedMedia1: reference1Exists,
					},
				},
			},
			args: args{
				medias: []*AnalysedMedia{analysedMedia1},
			},
			want: nil,
			wantEvents: []*ProgressEvent{
				{Type: ProgressEventAlreadyExists, Count: 1, Size: analysedMedia1.FoundMedia.Size()},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			referencerFactory = tt.fields.ReferencerFactory
			cataloger, err := NewCataloger(tt.newArgs.owner, tt.newArgs.options)
			if !tt.wantErr(t, err, fmt.Sprintf("NewCataloger(%v, %v)", tt.newArgs.owner, tt.newArgs.options)) {
				return
			}

			catcher := NewChanProgressEventCatcher()
			got, err := cataloger.Catalog(ctx, tt.args.medias, catcher)
			gotEvents := catcher.Catch()

			if !tt.wantErr(t, err, fmt.Sprintf("Catalog(%v)", tt.args.medias)) {
				return
			}

			assert.Equalf(t, tt.want, got, "Catalog(%v)", tt.args.medias)
			assert.ElementsMatchf(t, tt.wantEvents, gotEvents, "Catalog(%v)", tt.args.medias)
		})
	}
}

func TestCatalogerCreator_Catalog(t *testing.T) {
	ctx := context.Background()
	jan24 := time.Date(2020, time.January, 24, 0, 0, 0, 0, time.UTC)
	analysedMedia1 := newAnalysedMedia("file1.jpg", jan24, 12)
	reference1 := &CatalogReferenceStub{MediaIdValue: "media1", AlbumFolderNameValue: "album1"}
	//reference1Exists := &CatalogReferenceStub{MediaIdValue: "media1", AlbumFolderNameValue: "album1", ExistsValue: true}
	reference1AlbumCreated := &CatalogReferenceStub{MediaIdValue: "media1", AlbumFolderNameValue: "album1", AlbumCreatedValue: true}

	analysedMedia2 := newAnalysedMedia("file2.jpg", jan24, 13)
	reference2 := &CatalogReferenceStub{MediaIdValue: "media2", AlbumFolderNameValue: "album2"}

	analysedMedia3 := newAnalysedMedia("file3.jpg", jan24, 14)
	reference3 := &CatalogReferenceStub{MediaIdValue: "media3", AlbumFolderNameValue: "album3"}

	type fields struct {
		CatalogReferencer CatalogReferencer
		Filters           []CatalogerFilter
	}
	type args struct {
		medias []*AnalysedMedia
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       []*BackingUpMediaRequest
		wantEvents []*ProgressEvent
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "it shouldn't do anything if no media",
			fields: fields{
				CatalogReferencer: make(CatalogReferencerStub),
			},
			args: args{
				medias: nil,
			},
			want:       nil,
			wantEvents: nil,
			wantErr:    assert.NoError,
		},
		{
			name: "it should add the reference to each media",
			fields: fields{
				CatalogReferencer: CatalogReferencerStub{
					analysedMedia1: reference1,
				},
			},
			args: args{
				medias: []*AnalysedMedia{analysedMedia1},
			},
			want: []*BackingUpMediaRequest{
				{
					AnalysedMedia: analysedMedia1,
					Id:            reference1.MediaIdValue,
					FolderName:    reference1.AlbumFolderNameValue,
				},
			},
			wantEvents: []*ProgressEvent{
				{Type: ProgressEventCatalogued, Count: 1, Size: 12},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should skip filtered out medias with the right event",
			fields: fields{
				CatalogReferencer: CatalogReferencerStub{
					analysedMedia1: reference1,
					analysedMedia2: reference2,
					analysedMedia3: reference3,
				},
				Filters: []CatalogerFilter{
					CatalogerFilterFake{
						analysedMedia1.FoundMedia.String(): ProgressEventAlreadyExists,
						analysedMedia3.FoundMedia.String(): ProgressEventWrongAlbum,
					},
				},
			},
			args: args{
				medias: []*AnalysedMedia{analysedMedia1, analysedMedia2, analysedMedia3},
			},
			want: []*BackingUpMediaRequest{
				{
					AnalysedMedia: analysedMedia2,
					Id:            reference2.MediaIdValue,
					FolderName:    reference2.AlbumFolderNameValue,
				},
			},
			wantEvents: []*ProgressEvent{
				{Type: ProgressEventAlreadyExists, Count: 1, Size: analysedMedia1.FoundMedia.Size()},
				{Type: ProgressEventCatalogued, Count: 1, Size: analysedMedia2.FoundMedia.Size()},
				{Type: ProgressEventWrongAlbum, Count: 1, Size: analysedMedia3.FoundMedia.Size()},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should raise an event if an album has been created to store the media",
			fields: fields{
				CatalogReferencer: CatalogReferencerStub{
					analysedMedia1: reference1AlbumCreated,
				},
			},
			args: args{
				medias: []*AnalysedMedia{analysedMedia1},
			},
			want: []*BackingUpMediaRequest{
				{
					AnalysedMedia: analysedMedia1,
					Id:            reference1AlbumCreated.MediaIdValue,
					FolderName:    reference1AlbumCreated.AlbumFolderNameValue,
				},
			},
			wantEvents: []*ProgressEvent{
				{Type: ProgressEventAlbumCreated, Count: 1, Album: reference1AlbumCreated.AlbumFolderNameValue},
				{Type: ProgressEventCatalogued, Count: 1, Size: 12},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			catcher := NewChanProgressEventCatcher()
			c := Cataloger{
				CatalogReferencer: tt.fields.CatalogReferencer,
				CatalogerFilters:  tt.fields.Filters,
			}

			got, err := c.Catalog(ctx, tt.args.medias, catcher)
			gotEvents := catcher.Catch()

			if !tt.wantErr(t, err, fmt.Sprintf("Catalog(%v)", tt.args.medias)) {
				return
			}

			assert.Equalf(t, tt.want, got, "Catalog(%v)", tt.args.medias)
			assert.ElementsMatchf(t, tt.wantEvents, gotEvents, "Catalog(%v)", tt.args.medias)
		})
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

func (c *CatalogReferenceStub) AlbumCreated() bool {
	return c.AlbumCreatedValue
}

func (c *CatalogReferenceStub) Exists() bool {
	return c.ExistsValue
}

func (c *CatalogReferenceStub) MediaId() string {
	return c.MediaIdValue
}

func (c *CatalogReferenceStub) AlbumFolderName() string {
	return c.AlbumFolderNameValue
}

type ChanProgressEventCatcher chan *ProgressEvent

func NewChanProgressEventCatcher() ChanProgressEventCatcher {
	return make(chan *ProgressEvent, 255)
}

func (c ChanProgressEventCatcher) Catch() []*ProgressEvent {
	var events []*ProgressEvent
	for {
		select {
		case event := <-c:
			events = append(events, event)
		default:
			return events
		}
	}
}

type CatalogReferencerStub map[*AnalysedMedia]CatalogReference

func (c CatalogReferencerStub) Reference(ctx context.Context, medias []*AnalysedMedia) (map[*AnalysedMedia]CatalogReference, error) {
	result := make(map[*AnalysedMedia]CatalogReference)
	for _, media := range medias {
		for key, reference := range c {
			if key.FoundMedia.String() == media.FoundMedia.String() {
				result[media] = reference
				break
			}
		}
	}

	return result, nil
}

type CatalogerFilterFake map[string]ProgressEventType

func (c CatalogerFilterFake) FilterOut(media AnalysedMedia, reference CatalogReference) (ProgressEventType, bool) {
	if cause, ok := c[media.FoundMedia.String()]; ok {
		return cause, true
	}

	return ProgressEventCatalogued, false
}

type ReferencerFactoryFake struct {
	CreatorReferencer CatalogReferencer
	DryRunReferencer  CatalogReferencer
}

func (r *ReferencerFactoryFake) NewCreatorReferencer(ctx context.Context, owner ownermodel.Owner) (CatalogReferencer, error) {
	return r.CreatorReferencer, nil
}

func (r *ReferencerFactoryFake) NewDryRunReferencer(ctx context.Context, owner ownermodel.Owner) (CatalogReferencer, error) {
	return r.DryRunReferencer, nil
}
