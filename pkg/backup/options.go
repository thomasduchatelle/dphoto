package backup

type Options struct {
	RestrictedAlbumFolderName map[string]interface{} // RestrictedAlbumFolderName will restrict the media to only back up medias that are in one of these albums
	Listener                  interface{}            // Listener will receive progress events.
	SkipRejects               bool                   // SkipRejects mode will report any analysis error, or missing timestamp, and continue.
	AnalyserDecorator         AnalyserDecorator      // AnalyserDecorator is an optional decorator to add concept like caching (might be nil)
	DryRun                    bool                   // DryRun mode will not upload anything and do not create albums, still analyse
	ConcurrentAnalyser        int                    // ConcurrentAnalyser is the number of concurrent analyser (read files, compute hash, filter out duplicates, ...)
	ConcurrentCataloguer      int                    // ConcurrentCataloguer is the number of concurrent cataloguer (find album, create new albums)
	ConcurrentUploader        int                    // ConcurrentUploader is the number of concurrent uploader (upload files)
	BatchSize                 int                    // BatchSize is the number of items to read from the database at once (used by analyser) ; default to the maximum DynamoDB can handle
	RejectDir                 string                 // RejectDir is the directory where rejected files will be copied
}

func ReduceOptions(requestedOptions ...Options) Options {
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

		aggregated.RejectDir = mergeStringOption(aggregated.RejectDir, original.RejectDir)
		aggregated.ConcurrentAnalyser = mergeIntOption(aggregated.ConcurrentAnalyser, original.ConcurrentAnalyser)
		aggregated.ConcurrentCataloguer = mergeIntOption(aggregated.ConcurrentCataloguer, original.ConcurrentCataloguer)
		aggregated.ConcurrentUploader = mergeIntOption(aggregated.ConcurrentUploader, original.ConcurrentUploader)
		aggregated.BatchSize = mergeIntOption(aggregated.BatchSize, original.BatchSize)
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

func WithConcurrentAnalyser(concurrent int) Options {
	return Options{
		ConcurrentAnalyser: concurrent,
	}
}

func WithConcurrentCataloguer(concurrent int) Options {
	return Options{
		ConcurrentCataloguer: concurrent,
	}
}

func WithConcurrentUploader(concurrent int) Options {
	return Options{
		ConcurrentUploader: concurrent,
	}
}

func WithBatchSize(batchSize int) Options {
	return Options{
		BatchSize: batchSize,
	}
}

func OptionWithRejectDir(rejectDir string) Options {
	skip := false
	if rejectDir != "" {
		skip = true
	}
	return Options{
		SkipRejects: skip,
		RejectDir:   rejectDir,
	}
}

// NopeAnalyserDecorator is a default implementation for AnalyserDecorator which doesn't decorate the AnalyseMediaFunc.
type NopeAnalyserDecorator struct {
}

func (n *NopeAnalyserDecorator) Decorate(analyseFunc Analyser, observers ...AnalyserDecoratorObserver) Analyser {
	return analyseFunc
}

func mergeIntOption(current, value int) int {
	if current > 0 {
		return current
	}

	return value
}

func mergeStringOption(current, value string) string {
	if value != "" {
		return value
	}

	return current
}
