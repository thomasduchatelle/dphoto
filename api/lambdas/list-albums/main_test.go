package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func TestConvertCoversForREST(t *testing.T) {
	type args struct {
		covers []catalog.Cover
	}
	tests := []struct {
		name string
		args args
		want []CoverDTO
	}{
		{
			name: "it should return an empty slice when there is no cover",
			args: args{covers: nil},
			want: []CoverDTO{},
		},
		{
			name: "it should preserve the order and map each cover's identity and origin",
			args: args{covers: []catalog.Cover{
				{MediaId: catalog.MediaId("media-1"), Filename: "beach.jpg", Origin: catalog.CoverOriginRandom},
				{MediaId: catalog.MediaId("media-2"), Filename: "sunset.jpg", Origin: catalog.CoverOriginCherryPicked},
			}},
			want: []CoverDTO{
				{MediaId: "media-1", Filename: "beach.jpg", Origin: "RANDOM"},
				{MediaId: "media-2", Filename: "sunset.jpg", Origin: "CHERRY_PICKED"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, convertCoversForREST(tt.args.covers))
		})
	}
}
