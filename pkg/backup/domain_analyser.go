package backup

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"hash"
	"io"
	"path"
	"strings"
)

// SupportedExtensions is used by SourceVolume adapters to filter files they find
var SupportedExtensions = map[string]MediaType{
	"jpg":  MediaTypeImage,
	"jpeg": MediaTypeImage,
	"png":  MediaTypeImage,
	"gif":  MediaTypeImage,
	"webp": MediaTypeImage,
	"raw":  MediaTypeImage,
	//"bmp":  backupmodel.MediaTypeImage,
	"svg": MediaTypeImage,
	"eps": MediaTypeImage,

	"mkv":  MediaTypeVideo,
	"mts":  MediaTypeVideo,
	"avi":  MediaTypeVideo,
	"mp4":  MediaTypeVideo,
	"mpeg": MediaTypeVideo,
	"mov":  MediaTypeVideo,
	"wmv":  MediaTypeVideo,
	"webm": MediaTypeVideo,
}

// Analyser is a method to get the details of the media from the content of the file.
type Analyser interface {
	Analyse(ctx context.Context, found FoundMedia) (*AnalysedMedia, error)
}

// AnalyserDecorator allows to customize the Analyser behaviour, like adding a cache.
type AnalyserDecorator interface {
	Decorate(analyse Analyser, observers ...AnalyserDecoratorObserver) Analyser
}

// AnalyserDecoratorObserver is used to observe the decorator (if the cache hits, it will call this observer).
type AnalyserDecoratorObserver interface {
	OnSkipDelegateAnalyser(ctx context.Context, found FoundMedia) error
}

func newDefaultAnalyser(readers ...DetailsReader) Analyser {
	return &AnalyserFromMediaDetails{
		DetailsReaders: readers,
	}
}

// AnalyserFromMediaDetails is using DetailsReader to extract data from the file (EXIF, MP4, ...).
type AnalyserFromMediaDetails struct {
	options        DetailsReaderOptions
	DetailsReaders []DetailsReader
}

func (a *AnalyserFromMediaDetails) Analyse(ctx context.Context, found FoundMedia) (*AnalysedMedia, error) {
	reader, hasher, err := readerSpyingForHash(found)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s for analyse", found)
	}
	defer func() {
		if reader != nil {
			_ = reader.Close()
		}
	}()

	mediaType, details, err := a.extractDetails(found, reader, a.options)
	if err != nil {
		return nil, err
	}

	filehash := ""
	if !a.options.Fast {
		filehash, err = hasher.computeHash()
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compute file %s SHA256", found)
	}

	media := &AnalysedMedia{
		FoundMedia: found,
		Type:       mediaType,
		Sha256Hash: filehash,
		Details:    details,
	}
	return media, nil
}

func (a *AnalyserFromMediaDetails) extractDetails(found FoundMedia, reader io.Reader, options DetailsReaderOptions) (MediaType, *MediaDetails, error) {
	mediaType := getMediaType(found)

	loadedReaders := make([]string, len(a.DetailsReaders), len(a.DetailsReaders))
	for i, detailsReader := range a.DetailsReaders {
		loadedReaders[i] = fmt.Sprint(detailsReader)

		if detailsReader.Supports(found, mediaType) {
			details, err := detailsReader.ReadDetails(reader, options)
			return mediaType, details, errors.Wrapf(err, "failed to analyse %s file", found)
		}
	}

	return MediaTypeOther, nil, errors.Wrapf(ErrAnalyserNotSupported, "none of the details readers [%s] can parse %s", strings.Join(loadedReaders, ", "), found)
}

type hashSpy struct {
	reader    io.Reader
	shaWriter hash.Hash
}

type teeCloser struct {
	io.Reader
	io.Closer
}

func readerSpyingForHash(found FoundMedia) (io.ReadCloser, *hashSpy, error) {
	reader, err := found.ReadMedia()

	shaWriter := sha256.New()
	teeReader := io.TeeReader(reader, shaWriter)

	return teeCloser{Reader: teeReader, Closer: reader}, &hashSpy{
		reader:    teeReader,
		shaWriter: shaWriter,
	}, err
}

func (r *hashSpy) computeHash() (string, error) {
	_, err := io.Copy(io.Discard, r.reader)
	return hex.EncodeToString(r.shaWriter.Sum(nil)), err
}

func getMediaType(media FoundMedia) MediaType {
	extension := strings.TrimPrefix(strings.ToLower(path.Ext(media.MediaPath().Filename)), ".")
	if t, ok := SupportedExtensions[extension]; ok {
		return t
	}

	return MediaTypeOther
}
