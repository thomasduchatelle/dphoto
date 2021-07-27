package analyser

import (
	"crypto/sha256"
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/backup/interactors"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"hash"
	"io"
	"io/ioutil"
	"path"
	"reflect"
	"strings"
)

var SupportedExtensions = map[string]backupmodel.MediaType{
	"jpg":  backupmodel.MediaTypeImage,
	"jpeg": backupmodel.MediaTypeImage,
	"png":  backupmodel.MediaTypeImage,
	"gif":  backupmodel.MediaTypeImage,
	"webp": backupmodel.MediaTypeImage,
	"raw":  backupmodel.MediaTypeImage,
	"bmp":  backupmodel.MediaTypeImage,
	"svg":  backupmodel.MediaTypeImage,
	"eps":  backupmodel.MediaTypeImage,

	"mkv":  backupmodel.MediaTypeVideo,
	"mts":  backupmodel.MediaTypeVideo,
	"avi":  backupmodel.MediaTypeVideo,
	"mp4":  backupmodel.MediaTypeVideo,
	"mpeg": backupmodel.MediaTypeVideo,
	"mov":  backupmodel.MediaTypeVideo,
	"wmv":  backupmodel.MediaTypeVideo,
	"webm": backupmodel.MediaTypeVideo,
}

func AnalyseMedia(found backupmodel.FoundMedia) (*backupmodel.AnalysedMedia, error) {

	reader, hasher, err := readerSpyingForHash(found)

	mediaType, details, err := extractDetails(found, backupmodel.DetailsReaderOptions{}, func() (io.Reader, error) {
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
	return extractDetails(found, backupmodel.DetailsReaderOptions{Fast: true}, func() (io.Reader, error) {
		return found.ReadMedia()
	})
}

func extractDetails(found backupmodel.FoundMedia, options backupmodel.DetailsReaderOptions, readMedia func() (io.Reader, error)) (backupmodel.MediaType, *backupmodel.MediaDetails, error) {
	mediaType := getMediaType(found)

	details := &backupmodel.MediaDetails{}

	var matchingReaders []string
	for _, detailsReader := range interactors.DetailsReaders {
		if detailsReader.Supports(found, mediaType) {
			matchingReaders = append(matchingReaders, getType(detailsReader))
			content, err := readMedia()
			if err != nil {
				return mediaType, nil, errors.Wrapf(err, "failed to open media %s for analyse", found)
			}

			details, err = detailsReader.ReadDetails(content, options)
			if err != nil {
				return mediaType, nil, errors.Wrapf(err, "failed to analyse %s", found)
			}
		}
	}

	if details.DateTime.IsZero() {
		log.WithField("Media", found).Warnf("Modification date not found with readers: %s", strings.Join(matchingReaders, ", "))
		details.DateTime = found.LastModificationDate()
	}

	return mediaType, details, nil
}

type hashSpy struct {
	reader     io.Reader
	shaWriter  hash.Hash
	cachedHash string
}

func readerSpyingForHash(found backupmodel.FoundMedia) (io.Reader, *hashSpy, error) {
	reader, err := found.ReadMedia()
	if mediaWithHash, ok := found.(backupmodel.FoundMediaWithHash); ok {
		return reader, &hashSpy{cachedHash: mediaWithHash.Sha256Hash()}, err
	}

	shaWriter := sha256.New()
	teeReader := io.TeeReader(reader, shaWriter)

	return teeReader, &hashSpy{
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
	extension := strings.TrimPrefix(strings.ToLower(path.Ext(media.Filename())), ".")
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
