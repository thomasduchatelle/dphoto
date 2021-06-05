package interactive

import (
	"github.com/logrusorgru/aurora/v3"
)

func DeleteAlbum(operations CatalogOperations, record AlbumRecord) error {
	const pattern = "02/01/2006"
	proceed, ok := ReadBool(aurora.Sprintf("Are you sure you want to delete %s (%s) [%s -> %s] with %d medias in it?", aurora.Cyan(record.Name), record.FolderName, record.Start.Format(pattern), record.End.Format(pattern), record.Count), "y/N")
	if ok && proceed {
		return operations.DeleteAlbum(record.FolderName)
	}

	return nil
}
