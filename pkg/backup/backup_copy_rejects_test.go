package backup

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
	"time"
)

func TestCopyRejectsObserver_OnRejectedMedia(t *testing.T) {
	type args struct {
		found FoundMedia
		err   error
	}
	tests := []struct {
		name string
		args args
		want map[string][]byte
	}{
		{
			name: "it should copy the rejected files to the target directory",
			args: args{
				found: NewInMemoryMedia("somewhere/invalid.jpg", time.Now(), []byte("test")),
			},
			want: map[string][]byte{
				"somewhere_invalid.jpg": []byte("test"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp, err := os.MkdirTemp(os.TempDir(), "dphoto-unit-copyrejectsobserver")
			if !assert.NoError(t, err) {
				return
			}
			defer os.RemoveAll(temp)

			c := &copyRejectsObserver{
				RejectDir: temp,
			}

			err = c.OnRejectedMedia(context.Background(), tt.args.found, tt.args.err)
			if !assert.NoError(t, err) {
				return
			}

			files, err := os.ReadDir(temp)
			if assert.NoError(t, err) {
				got := make(map[string][]byte)
				for _, file := range files {
					content, err := os.ReadFile(path.Join(temp, file.Name()))
					if assert.NoError(t, err) {
						got[file.Name()] = content
					}
				}

				assert.Equal(t, tt.want, got)
			}
		})
	}
}
