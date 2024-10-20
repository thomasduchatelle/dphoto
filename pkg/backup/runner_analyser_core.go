package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"hash"
	"io"
	"io/ioutil"
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

type CoreAnalyser struct {
	options        DetailsReaderOptions
	detailsReaders []DetailsReaderAdapter // DetailsReaders is a list of specific details extractor can auto-register
}

func (a *CoreAnalyser) Analyse(found FoundMedia, analysedMediaObserver AnalysedMediaObserver, rejectedMediaObserver RejectedMediaObserver) {
	reader, hasher, err := readerSpyingForHash(found)
	if err != nil {
		rejectedMediaObserver.OnRejectedMedia(found, errors.Wrapf(err, "failed to open file %s for analyse", found))
		return
	}
	defer func() {
		if reader != nil {
			_ = reader.Close()
		}
	}()

	mediaType, details, err := a.extractDetails(found, reader, a.options)
	if err != nil {
		rejectedMediaObserver.OnRejectedMedia(found, err)
		return
	}

	filehash := ""
	if !a.options.Fast {
		filehash, err = hasher.computeHash()
	}
	if err != nil {
		rejectedMediaObserver.OnRejectedMedia(found, errors.Wrapf(err, "failed to compute file %s SHA256", found))
		return
	}

	analysedMediaObserver.OnAnalysedMedia(&AnalysedMedia{
		FoundMedia: found,
		Type:       mediaType,
		Sha256Hash: filehash,
		Details:    details,
	})
}

func (a *CoreAnalyser) extractDetails(found FoundMedia, reader io.Reader, options DetailsReaderOptions) (MediaType, *MediaDetails, error) {
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
	_, err := io.Copy(ioutil.Discard, r.reader)
	return hex.EncodeToString(r.shaWriter.Sum(nil)), err
}

func getMediaType(media FoundMedia) MediaType {
	extension := strings.TrimPrefix(strings.ToLower(path.Ext(media.MediaPath().Filename)), ".")
	if t, ok := SupportedExtensions[extension]; ok {
		return t
	}

	return MediaTypeOther
}
