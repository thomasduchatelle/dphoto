package analysers

import (
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers/avi"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers/exif"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers/m2ts"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers/mp4"
)

func ListDetailReaders() []backup.DetailsReaderAdapter {
	return []backup.DetailsReaderAdapter{
		new(avi.Parser),
		new(exif.Parser),
		new(m2ts.Parser),
		new(mp4.Parser),
	}
}
