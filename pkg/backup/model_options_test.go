package backup

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOption_ReduceOptions(t *testing.T) {
	noRestrictedAlbumFolder := make(map[string]interface{})

	type args struct {
		option Options
	}
	tests := []struct {
		name string
		args args
		want Options
	}{
		{
			name: "it should retain the SkipRejects option",
			args: args{
				option: OptionsSkipRejects(true),
			},
			want: Options{
				RestrictedAlbumFolderName: noRestrictedAlbumFolder,
				SkipRejects:               true,
			},
		},
		{
			name: "it should retain the RejectDir option and enable the SkipRejects option implicitly",
			args: args{
				option: OptionsWithRejectDir("foobar"),
			},
			want: Options{
				RestrictedAlbumFolderName: noRestrictedAlbumFolder,
				SkipRejects:               true,
				RejectDir:                 "foobar",
			},
		},
		{
			name: "it should not set the SkipRejects if RejectDir option is empty",
			args: args{
				option: OptionsWithRejectDir(""),
			},
			want: Options{
				RestrictedAlbumFolderName: noRestrictedAlbumFolder,
			},
		},
		{
			name: "it should support concurrency parameters",
			args: args{
				option: ReduceOptions(OptionsConcurrentAnalyserRoutines(2),
					OptionsConcurrentCataloguerRoutines(3),
					OptionsConcurrentUploaderRoutines(4)),
			},
			want: Options{
				RestrictedAlbumFolderName: noRestrictedAlbumFolder,
				ConcurrencyParameters: ConcurrencyParameters{
					ConcurrentAnalyserRoutines:   2,
					ConcurrentCataloguerRoutines: 3,
					ConcurrentUploaderRoutines:   4,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReduceOptions(Options{}, tt.args.option, Options{})
			assert.Equal(t, tt.want, got)
		})
	}
}
