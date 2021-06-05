package interactive

func EditAlbumDates(operations CatalogOperations, record AlbumRecord) error {
	start, ok := scanDate("Start date", record.Start)
	if !ok {
		return nil
	}

	end, ok := scanDate("End date", record.End)
	if !ok {
		return nil
	}

	return operations.UpdateAlbum(record.FolderName, start, end)
}
