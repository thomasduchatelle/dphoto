package interactive

import "github.com/logrusorgru/aurora/v3"

func EditAlbumName(operations CatalogOperations, record AlbumRecord) error {
	newName, ok := ReadString("Name of the album", record.Name)
	if !ok {
		return nil
	}

	if newName != record.Name {
		proceed, ok := ReadBool(aurora.Sprintf("Re-generate folder name /%s ?", aurora.Cyan(record.FolderName)), "[Y/n]")
		return operations.RenameAlbum(record.FolderName, newName, !ok || proceed)
	}

	return nil
}
