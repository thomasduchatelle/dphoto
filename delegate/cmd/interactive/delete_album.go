package interactive

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
)

func DeleteAlbum(operations CatalogOperations, record AlbumRecord) error {
	const pattern = "02/01/2006"
	proceed, ok := scanBool(aurora.Sprintf("Are you sure you want to delete %s (%s) [%s -> %s] with %d medias in it?", aurora.Cyan(record.Name), record.FolderName, record.Start.Format(pattern), record.End.Format(pattern), record.Count), "y/N")
	fmt.Println("scan bool", proceed, ok)
	if ok && proceed {
		return operations.DeleteAlbum(record.FolderName)
	}

	return nil
}
