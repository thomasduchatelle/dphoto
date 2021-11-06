package analyser

import (
	"crypto/sha256"
	"github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"
	"github.com/thomasduchatelle/dphoto/delegate/backup/interactors"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"hash"
	"io"
	"io/ioutil"
	"path"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var SupportedExtensions = map[string]backupmodel.MediaType{
	"jpg":  backupmodel.MediaTypeImage,
	"jpeg": backupmodel.MediaTypeImage,
	"png":  backupmodel.MediaTypeImage,
	"gif":  backupmodel.MediaTypeImage,
	"webp": backupmodel.MediaTypeImage,
	"raw":  backupmodel.MediaTypeImage,
	//"bmp":  backupmodel.MediaTypeImage,
	"svg": backupmodel.MediaTypeImage,
	"eps": backupmodel.MediaTypeImage,

	"mkv":  backupmodel.MediaTypeVideo,
	"mts":  backupmodel.MediaTypeVideo,
	"avi":  backupmodel.MediaTypeVideo,
	"mp4":  backupmodel.MediaTypeVideo,
	"mpeg": backupmodel.MediaTypeVideo,
	"mov":  backupmodel.MediaTypeVideo,
	"wmv":  backupmodel.MediaTypeVideo,
	"webm": backupmodel.MediaTypeVideo,
}

var DefaultMediaTimestamp = true // DefaultMediaTimestamp set to TRUE will use file last modification date as its timestamp when it can't be found within the file.

func AnalyseMedia(found backupmodel.FoundMedia) (*backupmodel.AnalysedMedia, error) {

	reader, hasher, err := readerSpyingForHash(found)
	defer reader.Close()

	mediaType, details, err := extractDetails(found, backupmodel.DetailsReaderOptions{}, func() (io.ReadCloser, error) {
		return reader, nil
	})
	if err != nil {
		return nil, err
	}

	fileHash, err := hasher.computeHash()

	return &backupmodel.AnalysedMedia{
		FoundMedia: found,
		Type:       mediaType,
		Signature: &backupmodel.FullMediaSignature{
			Sha256: fileHash,
			Size:   found.SimpleSignature().Size,
		},
		Details: details,
	}, errors.Wrapf(err, "failed to compute HASH of media %s", found)
}

func ExtractTypeAndDetails(found backupmodel.FoundMedia) (backupmodel.MediaType, *backupmodel.MediaDetails, error) {
	media, err := found.ReadMedia()
	if err != nil {
		return "", nil, err
	}
	defer media.Close()

	return extractDetails(found, backupmodel.DetailsReaderOptions{Fast: true}, func() (io.ReadCloser, error) {
		return media, err
	})
}

func extractDetails(found backupmodel.FoundMedia, options backupmodel.DetailsReaderOptions, readMedia func() (io.ReadCloser, error)) (backupmodel.MediaType, *backupmodel.MediaDetails, error) {
	mediaType := getMediaType(found)

	details := &backupmodel.MediaDetails{}

	var matchingReaders []string
	for _, detailsReader := range interactors.DetailsReaders {
		if detailsReader.Supports(found, mediaType) {
			matchingReaders = append(matchingReaders, getType(detailsReader))
			reader, err := readMedia()
			if err != nil {
				return mediaType, nil, errors.Wrapf(err, "failed to open media %s for analyse", found)
			}

			details, err = detailsReader.ReadDetails(reader, options)
			if err != nil {
				return mediaType, nil, errors.Wrapf(err, "failed to analyse %s", found)
			}
		}
	}

	if details.DateTime.IsZero() {
		log.WithField("Media", found).Warnf("Modification date not found with readers: %s", strings.Join(matchingReaders, ", "))
		if DefaultMediaTimestamp {
			details.DateTime = extractFromFileName(found.MediaPath().Filename, found.LastModificationDate())
		}
	}

	return mediaType, details, nil
}

func extractFromFileName(filename string, defaultDate time.Time) time.Time {
	formats := []struct {
		Regex            string
		Layout           string
		RemoveNonNumeric bool
	}{
		{"(199[0-9]|20[012][0-9])[^0-9A-Za-z]?[01][^0-9A-Za-z]?[0-9][^0-9A-Za-z]?[0-3][0-9][^0-9]?[0-2][0-9][^0-9A-Za-z]?[0-6][0-9][^0-9A-Za-z]?[0-6][0-9]", "20060102150405", true},
		{"(199[0-9]|20[012][0-9])[^0-9A-Za-z]?[01][^0-9A-Za-z]?[0-9][^0-9A-Za-z]?[0-3][0-9]", "20060102", true},
	}

	for _, f := range formats {
		date := regexp.MustCompile(f.Regex).FindString(filename)
		if f.RemoveNonNumeric {
			date = regexp.MustCompile("[^0-9]").ReplaceAllLiteralString(date, "")
		}
		if date != "" {
			parsedDate, err := time.Parse(f.Layout, date)
			if err == nil && !parsedDate.IsZero() && parsedDate.Before(time.Now()) && parsedDate.After(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)){
				return parsedDate
			}
		}
	}

	return defaultDate
}

type hashSpy struct {
	reader     io.Reader
	shaWriter  hash.Hash
	cachedHash string
}

type teeCloser struct {
	io.Reader
	io.Closer
}

func readerSpyingForHash(found backupmodel.FoundMedia) (io.ReadCloser, *hashSpy, error) {
	reader, err := found.ReadMedia()
	if mediaWithHash, ok := found.(backupmodel.FoundMediaWithHash); ok {
		return reader, &hashSpy{cachedHash: mediaWithHash.Sha256Hash()}, err
	}

	shaWriter := sha256.New()
	teeReader := io.TeeReader(reader, shaWriter)

	return teeCloser{Reader: teeReader, Closer: reader}, &hashSpy{
		reader:    teeReader,
		shaWriter: shaWriter,
	}, err
}

func (r *hashSpy) computeHash() (string, error) {
	if r.cachedHash != "" {
		return r.cachedHash, nil
	}

	_, err := io.Copy(ioutil.Discard, r.reader)
	return hex.EncodeToString(r.shaWriter.Sum(nil)), err
}

func getMediaType(media backupmodel.FoundMedia) backupmodel.MediaType {
	extension := strings.TrimPrefix(strings.ToLower(path.Ext(media.MediaPath().Filename)), ".")
	if t, ok := SupportedExtensions[extension]; ok {
		return t
	}

	return backupmodel.MediaTypeOther
}

func getType(myvar interface{}) (res string) {
	t := reflect.TypeOf(myvar)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		res += "*"
	}
	return fmt.Sprintf("%s%s/%s", res, t.PkgPath(), t.Name())
}
