package analysers

import (
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers/avi"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers/exif"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers/m2ts"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers/mp4"
)

func init() {
	backup.RegisterDetailsReader(new(avi.Parser))
	backup.RegisterDetailsReader(new(exif.Parser))
	backup.RegisterDetailsReader(new(m2ts.Parser))
	backup.RegisterDetailsReader(new(mp4.Parser))
}
