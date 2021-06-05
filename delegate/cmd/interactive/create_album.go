package interactive

func CreateAlbumForm(operations CatalogOperations, record AlbumRecord) error {
	creation := AlbumRecord{}
	ok := true

	creation.Name, ok = ReadString("Name of the album", record.Name)
	if !ok {
		return nil
	}

	creation.FolderName, _ = ReadString("Folder name (leave blank for automatically generated)", "")

	creation.Start, ok = ReadDate("Start date", record.Start)
	if !ok {
		return nil
	}

	creation.End, ok = ReadDate("End date", record.End)
	if !ok {
		return nil
	}

	return operations.Create(creation)
}
