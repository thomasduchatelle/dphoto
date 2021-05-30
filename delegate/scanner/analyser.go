package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/pkg/errors"
	"io"
	"path"
	"strings"
)

var SupportedExtensions = map[string]MediaType{
	"jpg":  MediaTypeImage,
	"jpeg": MediaTypeImage,
	"png":  MediaTypeImage,
	"gif":  MediaTypeImage,
	"webp": MediaTypeImage,
	"raw":  MediaTypeImage,
	"bmp":  MediaTypeImage,
	"svg":  MediaTypeImage,
	"eps":  MediaTypeImage,

	"mkv":  MediaTypeVideo,
	"mts":  MediaTypeVideo,
	"avi":  MediaTypeVideo,
	"mp4":  MediaTypeVideo,
	"mpeg": MediaTypeVideo,
	"mov":  MediaTypeVideo,
	"wmv":  MediaTypeVideo,
	"webm": MediaTypeVideo,
}

func AnalyseMedia(found FoundMedia) (*AnalysedMedia, error) {
	mediaType := getMediaType(found)

	details := &MediaDetails{
		DateTime: found.LastModificationDate(),
	}

	if mediaType == MediaTypeImage {
		content, err := found.ReadMedia()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open media %s for analyse", found)
		}

		details, err = ImageDetailsReader.ReadImageDetails(content, found.LastModificationDate())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to analyse image %s", found)
		}
	}

	fileHash, err := computeMediaHash(found)
	return &AnalysedMedia{
		FoundMedia: found,
		Type:       mediaType,
		Signature: &FullMediaSignature{
			Sha256: fileHash,
			Size:   found.SimpleSignature().Size,
		},
		Details: details,
	}, errors.Wrapf(err, "failed to compute HASH of media %s", found)
}

func computeMediaHash(found FoundMedia) (string, error) {
	if mediaWithHash, ok := found.(FoundMediaWithHash); ok {
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

func getMediaType(media FoundMedia) MediaType {
	extension := strings.TrimPrefix(strings.ToLower(path.Ext(media.Filename())), ".")
	if t, ok := SupportedExtensions[extension]; ok {
		return t
	}

	return MediaTypeOther
}
