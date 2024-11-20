package backup

import (
	"context"
	"io"
	"os"
	"path"
	"strings"
)

func newCopyRejectsObserver(rejectDir string) (RejectedMediaObserver, error) {
	err := os.MkdirAll(rejectDir, 0755)
	return &copyRejectsObserver{
		RejectDir: rejectDir,
	}, err
}

type copyRejectsObserver struct {
	RejectDir string
}

func (c *copyRejectsObserver) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	return c.copyFile(found)
}

func (c *copyRejectsObserver) copyFile(found FoundMedia) error {
	filename := strings.ReplaceAll(path.Join(found.MediaPath().Path, found.MediaPath().Filename), "/", "_")
	fullPath := path.Join(c.RejectDir, filename)

	media, err := found.ReadMedia()
	if err != nil {
		return err
	}
	defer media.Close()

	all, err := io.ReadAll(media)
	if err != nil {
		return err
	}

	return os.WriteFile(fullPath, all, 0644)
}
