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

type analyseMedia struct {
	options DetailsReaderOptions
}

func newBackupAnalyseMedia() runnerAnalyser {
	analyser := analyseMedia{
		options: DetailsReaderOptions{Fast: false},
	}
	return analyser.analyseMedia
}

func (a *analyseMedia) analyseMedia(found FoundMedia, eventChannel chan *ProgressEvent) (*AnalysedMedia, error) {
	reader, hasher, err := readerSpyingForHash(found)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s for analyse", found)
	}
	defer func() {
		if reader != nil {
			_ = reader.Close()
		}
	}()

	mediaType, details, err := extractDetails(found, reader, a.options)
	if err != nil {
		return nil, err
	}

	filehash := ""
	if !a.options.Fast {
		filehash, err = hasher.computeHash()
	}

	eventChannel <- &ProgressEvent{Type: ProgressEventAnalysed, Count: 1, Size: found.Size()}

	return &AnalysedMedia{
		FoundMedia: found,
		Type:       mediaType,
		Sha256Hash: filehash,
		Details:    details,
	}, errors.Wrapf(err, "failed to compute file %s SHA256", found)
}

func extractDetails(found FoundMedia, reader io.Reader, options DetailsReaderOptions) (MediaType, *MediaDetails, error) {
	mediaType := getMediaType(found)

	loadedReaders := make([]string, len(detailsReaders), len(detailsReaders))
	for i, detailsReader := range detailsReaders {
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
