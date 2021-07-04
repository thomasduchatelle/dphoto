package analyser

import (
	"crypto/sha256"
	"duchatelle.io/dphoto/dphoto/backup/interactors"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"encoding/hex"
	"github.com/pkg/errors"
	"io"
	"path"
	"path/filepath"
	"strings"
)

var SupportedExtensions = map[string]model.MediaType{
	"jpg":  model.MediaTypeImage,
	"jpeg": model.MediaTypeImage,
	"png":  model.MediaTypeImage,
	"gif":  model.MediaTypeImage,
	"webp": model.MediaTypeImage,
	"raw":  model.MediaTypeImage,
	"bmp":  model.MediaTypeImage,
	"svg":  model.MediaTypeImage,
	"eps":  model.MediaTypeImage,

	"mkv":  model.MediaTypeVideo,
	"mts":  model.MediaTypeVideo,
	"avi":  model.MediaTypeVideo,
	"mp4":  model.MediaTypeVideo,
	"mpeg": model.MediaTypeVideo,
	"mov":  model.MediaTypeVideo,
	"wmv":  model.MediaTypeVideo,
	"webm": model.MediaTypeVideo,
}

func AnalyseMedia(found model.FoundMedia) (*model.AnalysedMedia, error) {
	mediaType, details, err := ExtractTypeAndDetails(found)
	if err != nil {
		return nil, err
	}

	fileHash, err := computeMediaHash(found) // todo - do it while analysing files
	return &model.AnalysedMedia{
		FoundMedia: found,
		Type:       mediaType,
		Signature: &model.FullMediaSignature{
			Sha256: fileHash,
			Size:   found.SimpleSignature().Size,
		},
		Details: details,
	}, errors.Wrapf(err, "failed to compute HASH of media %s", found)
}

func ExtractTypeAndDetails(found model.FoundMedia) (model.MediaType, *model.MediaDetails, error) {
	mediaType := getMediaType(found)

	details := &model.MediaDetails{}

	var detailsReaderType interactors.DetailsReaderType

	switch {
	case mediaType == model.MediaTypeImage:
		detailsReaderType = interactors.DetailsReaderTypeImage

	case strings.ToUpper(filepath.Ext(found.Filename())) == ".MTS":
		detailsReaderType = interactors.DetailsReaderTypeM2TS
	}

	if detailsReader, ok := interactors.DetailsReaders[detailsReaderType]; ok {
		content, err := found.ReadMedia()
		if err != nil {
			return mediaType, nil, errors.Wrapf(err, "failed to open media %s for analyse", found)
		}

		details, err = detailsReader.ReadDetails(content, model.DetailsReaderOptions{Fast: true})
		if err != nil {
			return mediaType, nil, errors.Wrapf(err, "failed to analyse %s", found)
		}
	}

	if details.DateTime.IsZero() {
		details.DateTime = found.LastModificationDate()
	}

	return mediaType, details, nil
}

func computeMediaHash(found model.FoundMedia) (string, error) {
	if mediaWithHash, ok := found.(model.FoundMediaWithHash); ok {
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

func getMediaType(media model.FoundMedia) model.MediaType {
	extension := strings.TrimPrefix(strings.ToLower(path.Ext(media.Filename())), ".")
	if t, ok := SupportedExtensions[extension]; ok {
		return t
	}

	return model.MediaTypeOther
}
