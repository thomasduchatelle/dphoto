package backup

type Options struct {
	RestrictedAlbumFolderName map[string]interface{} // RestrictedAlbumFolderName will restrict the media to only back up medias that are in one of these albums
	Listener                  interface{}            // Listener will receive progress events.
	SkipRejects               bool                   // SkipRejects mode will report any analysis error, or missing timestamp, and continue.
	AnalyserDecorator         AnalyserDecorator      // AnalyserDecorator is an optional decorator to add concept like caching (might be nil)
}

type AnalyserDecorator interface {
	Decorate(analyseFunc RunnerAnalyser) RunnerAnalyser
}

func readOptions(requestedOptions []Options) Options {
	aggregated := Options{
		RestrictedAlbumFolderName: make(map[string]interface{}),
	}
	for _, original := range requestedOptions {
		for folderName := range original.RestrictedAlbumFolderName {
			aggregated.RestrictedAlbumFolderName[folderName] = nil
		}

		if original.Listener != nil {
			aggregated.Listener = original.Listener
		}

		if original.AnalyserDecorator != nil {
			aggregated.AnalyserDecorator = original.AnalyserDecorator
		}

		aggregated.SkipRejects = aggregated.SkipRejects || original.SkipRejects
	}

	return aggregated
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

// WithCachedAnalysis adds a decorator on analysis function ; argument can be nil.
func (o Options) WithCachedAnalysis(analyserDecorator AnalyserDecorator) Options {
	o.AnalyserDecorator = analyserDecorator
	return o
}

// GetAnalyserDecorator is returning the AnalyserDecorator or NopeAnalyserDecorator, never nil.
func (o Options) GetAnalyserDecorator() AnalyserDecorator {
	if o.AnalyserDecorator != nil {
		return o.AnalyserDecorator
	}

	return new(NopeAnalyserDecorator)
}

// NopeAnalyserDecorator is a default implementation for AnalyserDecorator which doesn't decorate the AnalyseMediaFunc.
type NopeAnalyserDecorator struct {
}

func (n *NopeAnalyserDecorator) Decorate(analyseFunc RunnerAnalyser) RunnerAnalyser {
	return analyseFunc
}
