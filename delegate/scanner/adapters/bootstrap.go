package adapters

import (
	"duchatelle.io/dphoto/dphoto/scanner"
	"duchatelle.io/dphoto/dphoto/scanner/adapters/filesystem"
	"duchatelle.io/dphoto/dphoto/scanner/adapters/images"
)

func init() {
	scanner.ImageDetailsReader = new(images.ExifReader)
	scanner.SourceAdapters[scanner.VolumeTypeFileSystem] = new(filesystem.FsHandler)
}
