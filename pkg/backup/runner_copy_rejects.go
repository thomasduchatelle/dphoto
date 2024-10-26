package backup

import (
	"context"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"strings"
)

func NewCopyRejectsObserver(rejectDir string) (RejectedMediaObserver, error) {
	err := os.MkdirAll(rejectDir, 0755)
	return &CopyRejectsObserver{
		RejectDir: rejectDir,
	}, err
}

type CopyRejectsObserver struct {
	RejectDir string
}

func (c *CopyRejectsObserver) OnRejectedMedia(ctx context.Context, found FoundMedia, cause error) error {
	err := c.copyFile(found)
	if err != nil {
		// note - 'OnRejectedMedia' is called within error handling (hence not returning errors itself) ; suggested improvement is to pass the error in a channel where the consumer has its own error handling
		log.WithError(err).Errorf("failed to copy rejected file %s to %s", found.MediaPath().Path, c.RejectDir)
	}
	return nil
}

func (c *CopyRejectsObserver) copyFile(found FoundMedia) error {
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
