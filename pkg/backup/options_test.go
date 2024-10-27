package backup_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"testing"
)

func TestOption_ReduceOptions(t *testing.T) {
	noRestrictedAlbumFolder := make(map[string]interface{})

	type args struct {
		option backup.Options
	}
	tests := []struct {
		name string
		args args
		want backup.Options
	}{
		{
			name: "it should retain the SkipRejects option",
			args: args{
				option: backup.OptionsSkipRejects(true),
			},
			want: backup.Options{
				RestrictedAlbumFolderName: noRestrictedAlbumFolder,
				SkipRejects:               true,
			},
		},
		{
			name: "it should retain the RejectDir option and enable the SkipRejects option implicitly",
			args: args{
				option: backup.OptionWithRejectDir("foobar"),
			},
			want: backup.Options{
				RestrictedAlbumFolderName: noRestrictedAlbumFolder,
				SkipRejects:               true,
				RejectDir:                 "foobar",
			},
		},
		{
			name: "it should not set the SkipRejects if RejectDir option is empty",
			args: args{
				option: backup.OptionWithRejectDir(""),
			},
			want: backup.Options{
				RestrictedAlbumFolderName: noRestrictedAlbumFolder,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := backup.ReduceOptions(backup.Options{}, tt.args.option, backup.Options{})
			assert.Equal(t, tt.want, got)
		})
	}
}
