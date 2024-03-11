package analysiscache

import (
	"fmt"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"time"
)

type Key struct {
	AbsolutePath string
	Size         int
}

func (k Key) SerialisedKey() []byte {
	return []byte(fmt.Sprintf("%s##%d", k.AbsolutePath, k.Size))
}

func (k Key) String() string {
	return string(k.SerialisedKey())
}

type Payload struct {
	LastModification time.Time           `json:"lastModification,omitempty"`
	Type             string              `json:"type,omitempty"`
	Sha256Hash       string              `json:"sha256Hash,omitempty"`
	Details          backup.MediaDetails `json:"details,omitempty"`
}
