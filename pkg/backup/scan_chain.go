package backup

func postCatalogFiltersList(options Options) []CataloguerFilter {
	filters := []CataloguerFilter{
		mustNotExists(),
		mustBeUniqueInVolume(),
	}

	if len(options.RestrictedAlbumFolderName) > 0 {
		var albumFolderNames []string
		for albumFolderName := range options.RestrictedAlbumFolderName {
			albumFolderNames = append(albumFolderNames, albumFolderName)
		}
		filters = append(filters, mustBeInAlbum(albumFolderNames...))
	}

	return filters
}
