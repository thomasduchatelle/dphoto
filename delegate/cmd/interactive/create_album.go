package interactive

func CreateAlbumForm(operations CatalogOperations, record AlbumRecord) error {
	creation := AlbumRecord{}
	ok := true

	creation.Name, ok = scanString("Name of the album", record.Name)
	if !ok {
		return nil
	}

	creation.FolderName, _ = scanString("Folder name (leave blank for automatically generated)", "")

	creation.Start, ok = scanDate("Start date", record.Start)
	if !ok {
		return nil
	}

	creation.End, ok = scanDate("End date", record.End)
	if !ok {
		return nil
	}

	return operations.Create(creation)
}
