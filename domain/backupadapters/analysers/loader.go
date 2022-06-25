package analysers

import (
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"github.com/thomasduchatelle/dphoto/domain/backupadapters/analysers/avi"
	"github.com/thomasduchatelle/dphoto/domain/backupadapters/analysers/exif"
	"github.com/thomasduchatelle/dphoto/domain/backupadapters/analysers/m2ts"
	"github.com/thomasduchatelle/dphoto/domain/backupadapters/analysers/mp4"
)

func init() {
	backup.RegisterDetailsReader(new(avi.Parser))
	backup.RegisterDetailsReader(new(exif.Parser))
	backup.RegisterDetailsReader(new(m2ts.Parser))
	backup.RegisterDetailsReader(new(mp4.Parser))
}
