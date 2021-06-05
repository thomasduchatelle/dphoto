package interactive

func EditAlbumDates(operations CatalogOperations, record AlbumRecord) error {
	start, ok := ReadDate("Start date", record.Start)
	if !ok {
		return nil
	}

	end, ok := ReadDate("End date", record.End)
	if !ok {
		return nil
	}

	return operations.UpdateAlbum(record.FolderName, start, end)
}
