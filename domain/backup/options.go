package backup

type Options struct {
	RestrictedAlbumFolderName map[string]interface{} // RestrictedAlbumFolderName will restrict the media to only back up medias that are in one of these albums
	Listener                  interface{}            // Listener will receive progress events.
}

func readOptions(optionSlice []Options) Options {
	opt := Options{
		RestrictedAlbumFolderName: make(map[string]interface{}),
	}
	for _, o := range optionSlice {
		for folderName := range o.RestrictedAlbumFolderName {
			opt.RestrictedAlbumFolderName[folderName] = nil
		}

		if o.Listener != nil {
			opt.Listener = o.Listener
		}
	}

	return opt
}

// OptionWithListener creates an option with a listener
func OptionWithListener(listener interface{}) Options {
	return Options{
		Listener: listener,
	}
}

// OptionOnlyAlbums restricts backed up medias to those in these albums
func OptionOnlyAlbums(albums ...string) Options {
	options := Options{
		RestrictedAlbumFolderName: make(map[string]interface{}),
	}

	for _, album := range albums {
		options.RestrictedAlbumFolderName[album] = nil
	}

	return options
}
