package backup

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"maps"
	"slices"
	"testing"
	"time"
)

var (
	ErrTestNoDetails = errors.New("TEST no details have been set")
)

func Test_analyserAggregate_OnFoundMedia(t *testing.T) {
	now := time.Now()
	noDate := time.Time{}
	image1 := NewInMemoryMedia("image-1.jpg", time.Time{}, []byte("nice picture"))
	video2 := NewInMemoryMedia("video-2.mp4", time.Time{}, []byte("nice video"))

	type fields struct {
		detailReaders []DetailsReader
	}
	type args struct {
		media FoundMedia
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantObserved []AnalysedMedia
		wantRejected map[string]assert.ErrorAssertionFunc
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "it should accept valid analysed IMAGE",
			fields: fields{
				detailReaders: detailsReadersAlwaysReturning(detailsWithDate(now)),
			},
			args: args{
				media: image1,
			},
			wantObserved: []AnalysedMedia{{
				FoundMedia: image1,
				Type:       MediaTypeImage,
				Sha256Hash: "88fce3dd1fbbba642c86785206be78fb9a3004f115744fecd27c79668b5219e6",
				Details:    detailsWithDate(now),
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should accept valid analysed VIDEO",
			fields: fields{
				detailReaders: detailsReadersAlwaysReturning(detailsWithDate(now)),
			},
			args: args{
				media: video2,
			},
			wantObserved: []AnalysedMedia{{
				FoundMedia: video2,
				Type:       MediaTypeVideo,
				Sha256Hash: "29528cf72ae8b60255642189a9e1813d8fa0f38e2b2f8a5e6f183a44f02182d4",
				Details:    detailsWithDate(now),
			}},
			wantErr: assert.NoError,
		},
		{
			name: "it should reject medias not supported",
			fields: fields{
				detailReaders: nil,
			},
			args: args{
				media: video2,
			},
			wantRejected: map[string]assert.ErrorAssertionFunc{
				video2.MediaPath().Filename: func(t assert.TestingT, err error, i ...interface{}) bool {
					return assert.ErrorIs(t, err, ErrAnalyserNotSupported, i...)
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should reject medias without valid date",
			fields: fields{
				detailReaders: detailsReadersAlwaysReturning(detailsWithDate(noDate)),
			},
			args: args{
				media: video2,
			},
			wantRejected: map[string]assert.ErrorAssertionFunc{
				video2.MediaPath().Filename: func(t assert.TestingT, err error, i ...interface{}) bool {
					return assert.ErrorIs(t, err, ErrAnalyserNoDateTime, i...)
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should reject media an error if the details reader fails",
			fields: fields{
				detailReaders: []DetailsReader{&DetailsReaderFake{}},
			},
			args: args{
				media: video2,
			},
			wantRejected: map[string]assert.ErrorAssertionFunc{
				video2.MediaPath().Filename: func(t assert.TestingT, err error, i ...interface{}) bool {
					return assert.ErrorIs(t, err, ErrTestNoDetails, i...)
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observer := new(AnalyserObserverFake)

			a := &analyserAggregate{
				analyser: &AnalyserFromMediaDetails{
					detailsReaders: tt.fields.detailReaders,
				},
				analysedMediaObservers: []AnalysedMediaObserver{observer},
				rejectedMediaObservers: []RejectedMediaObserver{observer},
			}

			err := a.OnFoundMedia(context.Background(), tt.args.media)
			if tt.wantErr(t, err, fmt.Sprintf("OnFoundMedia(%v)", tt.args.media)) {
				assert.Equal(t, tt.wantObserved, observer.Analysed)
				if assert.Equal(t, slices.Sorted(maps.Keys(tt.wantRejected)), slices.Sorted(maps.Keys(observer.Rejected)), "errors are not matching") {
					for filename, wantErr := range tt.wantRejected {
						wantErr(t, observer.Rejected[filename])
					}
				}
			}
		})
	}
}

func detailsWithDate(dateTime time.Time) *MediaDetails {
	return &MediaDetails{
		DateTime: dateTime,
	}
}

func detailsReadersAlwaysReturning(details *MediaDetails) []DetailsReader {
	return []DetailsReader{&DetailsReaderFake{Details: details}}
}

type DetailsReaderFake struct {
	Details *MediaDetails
}

func (d *DetailsReaderFake) Supports(media FoundMedia, mediaType MediaType) bool {
	return true
}

func (d *DetailsReaderFake) ReadDetails(reader io.Reader, options DetailsReaderOptions) (*MediaDetails, error) {
	if d.Details == nil {
		return nil, ErrTestNoDetails
	}
	return d.Details, nil
}

type AnalyserObserverFake struct {
	Analysed []AnalysedMedia
	Rejected map[string]error
}

func (a *AnalyserObserverFake) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	a.Analysed = append(a.Analysed, *media)
	return nil
}

func (a *AnalyserObserverFake) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	if a.Rejected == nil {
		a.Rejected = make(map[string]error)
	}
	a.Rejected[found.MediaPath().Filename] = cause
	return nil
}
