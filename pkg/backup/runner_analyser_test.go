package backup

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAnalyserAcceptance_simple(t *testing.T) {
	// it should filter out files without dates
	// it should call the observer for each analysed media
	// it should update the progress event channel
	// it should add the media to the next step channel
	// (it should interrupt the backup if an error is encountered)

	validAnalysedMedia := &AnalysedMedia{
		Details: &MediaDetails{
			DateTime: time.Now(),
		},
	}
	aMedia := NewInMemoryMedia("a-media-without-content", time.Time{}, nil)

	type args struct {
		Options       Options
		AnalysedMedia *AnalysedMedia
	}
	tests := []struct {
		name         string
		args         args
		wantAnalysed []*AnalysedMedia
		wantRejected map[FoundMedia]error
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "should accept valid analysed media",
			args: args{
				AnalysedMedia: validAnalysedMedia,
			},
			wantAnalysed: []*AnalysedMedia{validAnalysedMedia},
			wantRejected: nil,
			wantErr:      assert.NoError,
		},
		{
			name: "should reject invalid analysed media without failing the process if skipRejects is TRUE",
			args: args{
				Options:       OptionsSkipRejects(true),
				AnalysedMedia: &AnalysedMedia{FoundMedia: aMedia, Details: &MediaDetails{}},
			},
			wantAnalysed: nil,
			wantRejected: map[FoundMedia]error{
				aMedia: ErrAnalyserNoDateTime,
			},
			wantErr: assert.NoError,
		},
		{
			name: "should interrupt the process if an error is encountered while analysing",
			args: args{
				AnalysedMedia: &AnalysedMedia{FoundMedia: aMedia, Details: &MediaDetails{}},
			},
			wantAnalysed: nil,
			wantRejected: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrAnalyserNoDateTime, i...)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeObserver := new(AnalyserObserverFake)
			handler := newAnalyserObserverChain(tt.args.Options, fakeObserver)

			err := handler.OnAnalysedMedia(context.Background(), tt.args.AnalysedMedia)
			if tt.wantErr(t, err, tt.name) {
				assert.Equal(t, fakeObserver.Analysed, tt.wantAnalysed)
				assert.Equal(t, fakeObserver.Rejected, tt.wantRejected)
			}

		})
	}
}

type AnalyserObserverFake struct {
	Analysed []*AnalysedMedia
	Rejected map[FoundMedia]error
}

func (a *AnalyserObserverFake) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	a.Analysed = append(a.Analysed, media)
	return nil
}

func (a *AnalyserObserverFake) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	if a.Rejected == nil {
		a.Rejected = make(map[FoundMedia]error)
	}
	a.Rejected[found] = cause
	return nil
}
