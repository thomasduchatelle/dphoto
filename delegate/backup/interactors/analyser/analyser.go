package analyser

import (
	"crypto/sha256"
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/backup/interactors"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
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
	mediaType, details, err := ExtractTypeAndDetails(found)
	if err != nil {
		return nil, err
	}

	fileHash, err := computeMediaHash(found) // todo - do it while analysing files
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
	mediaType := getMediaType(found)

	details := &backupmodel.MediaDetails{}

	var matchingReaders []string
	for _, detailsReader := range interactors.DetailsReaders {
		if detailsReader.Supports(found, mediaType) {
			matchingReaders = append(matchingReaders, getType(detailsReader))
			content, err := found.ReadMedia()
			if err != nil {
				return mediaType, nil, errors.Wrapf(err, "failed to open media %s for analyse", found)
			}

			details, err = detailsReader.ReadDetails(content, backupmodel.DetailsReaderOptions{Fast: true})
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

func computeMediaHash(found backupmodel.FoundMedia) (string, error) {
	if mediaWithHash, ok := found.(backupmodel.FoundMediaWithHash); ok {
		return mediaWithHash.Sha256Hash(), nil
	}

	shaWriter := sha256.New()
	reader, err := found.ReadMedia()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(shaWriter, reader)
	return hex.EncodeToString(shaWriter.Sum(nil)), err
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
