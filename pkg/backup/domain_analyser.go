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
	"slices"
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

type mediaReader struct {
	options        DetailsReaderOptions
	detailsReaders []DetailsReader
}

func (a *mediaReader) analyseMedia(found FoundMedia) (*AnalysedMedia, error) {
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

func (a *mediaReader) extractDetails(found FoundMedia, reader io.Reader, options DetailsReaderOptions) (MediaType, *MediaDetails, error) {
	mediaType := getMediaType(found)

	loadedReaders := make([]string, len(a.detailsReaders), len(a.detailsReaders))
	for i, detailsReader := range a.detailsReaders {
		loadedReaders[i] = fmt.Sprint(detailsReader)

		if detailsReader.Supports(found, mediaType) {
			details, err := detailsReader.ReadDetails(reader, options)
			return mediaType, details, errors.Wrapf(err, "failed to analyse %s file", found)
		}
	}

	return MediaTypeOther, nil, errors.Errorf("%s not supported. Parser loaded: %s", found, strings.Join(loadedReaders, ", "))
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

type analyserAdapter struct {
	analyser     Analyser
	analysed     AnalysedMediaObservers
	beforeFilter AnalysedMediaObservers
	filteredOut  RejectedMediaObservers
	rejected     RejectedMediaObservers
}

func (a *analyserAdapter) OnFoundMedia(ctx context.Context, media FoundMedia) error {
	return a.analyser.Analyse(
		ctx,
		media,
		slices.Concat(a.beforeFilter, []AnalysedMediaObserver{&analyserNoDateTimeFilter{
			analysedMediaObserver: a.analysed,
			rejectedMediaObserver: a.filteredOut,
		}}),
		&a.rejected,
	)
}

type analyserNoDateTimeFilter struct {
	analysedMediaObserver AnalysedMediaObserver
	rejectedMediaObserver RejectedMediaObserver
}

func (a *analyserNoDateTimeFilter) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	if media.Details.DateTime.IsZero() {
		return a.rejectedMediaObserver.OnRejectedMedia(ctx, media.FoundMedia, ErrAnalyserNoDateTime)
	}

	return a.analysedMediaObserver.OnAnalysedMedia(ctx, media)
}
