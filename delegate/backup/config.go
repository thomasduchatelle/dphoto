package backup

import (
	"github.com/pkg/errors"
	"os"
	"path"
)

const (
	ImageReaderThread = 4
)

var (
	ConfigurationDir string
	LocalMediaPath   string
)

func init() {
	ConfigurationDir = "/etc/dphoto"
	LocalMediaPath = "/var/dphoto"
}

func temporaryMediaPath(backupId string) (string, error) {
	temp := path.Join(LocalMediaPath, ".temp", backupId)
	err := os.MkdirAll(temp, 0744)

	return temp, errors.Wrapf(err, "Can't create temporary storage %s", temp)
}
