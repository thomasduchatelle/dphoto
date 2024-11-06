package backup

type Options struct {
	RestrictedAlbumFolderName map[string]interface{} // RestrictedAlbumFolderName will restrict the media to only back up medias that are in one of these albums
	Listener                  interface{}            // Listener will receive progress events.
	SkipRejects               bool                   // SkipRejects mode will report any analysis error, or missing timestamp, and continue.
	AnalyserDecorator         AnalyserDecorator      // AnalyserDecorator is an optional decorator to add concept like caching (might be nil)
	DryRun                    bool                   // DryRun mode will not upload anything and do not create albums, still analyse
	ConcurrencyParameters     ConcurrencyParameters
	BatchSize                 int    // BatchSize is the number of items to read from the database at once (used by analyser) ; default to the maximum DynamoDB can handle
	RejectDir                 string // RejectDir is the directory where rejected files will be copied
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
		aggregated.ConcurrencyParameters.ConcurrentAnalyserRoutines = mergeIntOption(aggregated.ConcurrencyParameters.ConcurrentAnalyserRoutines, original.ConcurrencyParameters.ConcurrentAnalyserRoutines)
		aggregated.ConcurrencyParameters.ConcurrentCataloguerRoutines = mergeIntOption(aggregated.ConcurrencyParameters.ConcurrentCataloguerRoutines, original.ConcurrencyParameters.ConcurrentCataloguerRoutines)
		aggregated.ConcurrencyParameters.ConcurrentUploaderRoutines = mergeIntOption(aggregated.ConcurrencyParameters.ConcurrentUploaderRoutines, original.ConcurrencyParameters.ConcurrentUploaderRoutines)
		aggregated.BatchSize = mergeIntOption(aggregated.BatchSize, original.BatchSize)
	}

	return aggregated
}

type ConcurrencyParameters struct {
	ConcurrentAnalyserRoutines   int // ConcurrentAnalyserRoutines is the number of concurrent analyser (read files, compute hash, filter out duplicates, ...)
	ConcurrentCataloguerRoutines int // ConcurrentCataloguerRoutines is the number of concurrent cataloguer (find album, create new albums)
	ConcurrentUploaderRoutines   int // ConcurrentUploaderRoutines is the number of concurrent uploader (upload files)
}

func (c ConcurrencyParameters) NumberOfConcurrentAnalyserRoutines() int {
	return defaultValue(c.ConcurrentAnalyserRoutines, 1)
}

func (c ConcurrencyParameters) NumberOfConcurrentCataloguerRoutines() int {
	return defaultValue(c.ConcurrentCataloguerRoutines, 1)
}

func (c ConcurrencyParameters) NumberOfConcurrentUploaderRoutines() int {
	return defaultValue(c.ConcurrentUploaderRoutines, 1)
}

func defaultValue(value, fallback int) int {
	if value == 0 {
		return fallback
	}

	return value
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

// OptionsSkipRejects disables the strict mode and ignores invalid files (wrong / no date, ...)
func OptionsSkipRejects(skip bool) Options {
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

func (o Options) GetBatchSize() int {
	return defaultValue(o.BatchSize, 1)
}

func OptionsConcurrentAnalyserRoutines(concurrent int) Options {
	return Options{
		ConcurrencyParameters: ConcurrencyParameters{
			ConcurrentAnalyserRoutines: concurrent,
		},
	}
}

func OptionsConcurrentCataloguerRoutines(concurrent int) Options {
	return Options{
		ConcurrencyParameters: ConcurrencyParameters{
			ConcurrentCataloguerRoutines: concurrent,
		},
	}
}

func OptionsConcurrentUploaderRoutines(concurrent int) Options {
	return Options{
		ConcurrencyParameters: ConcurrencyParameters{
			ConcurrentUploaderRoutines: concurrent,
		},
	}
}

func OptionsBatchSize(batchSize int) Options {
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
