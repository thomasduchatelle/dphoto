package backup

type Options struct {
	RestrictedAlbumFolderName map[string]interface{} // RestrictedAlbumFolderName will restrict the media to only back up medias that are in one of these albums
	Listener                  interface{}            // Listener will receive progress events.
	SkipRejects               bool                   // SkipRejects mode will report any analysis error, or missing timestamp, and continue.
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

		opt.SkipRejects = opt.SkipRejects || o.SkipRejects
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

// OptionSkipRejects disables the strict mode and ignores invalid files (wrong / no date, ...)
func OptionSkipRejects(skip bool) Options {
	return Options{
		SkipRejects: skip,
	}
}
