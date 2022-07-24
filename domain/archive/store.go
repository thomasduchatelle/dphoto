package archive

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"path"
	"strings"
)

// Store save the file content, register it, and generates miniature. Return the new filename.
func Store(request *StoreRequest) (string, error) {
	filename, stored, err := storeInArchive(request)
	if err != nil || !stored {
		return filename, err
	}

	reader, err := request.Open()
	if err != nil {
		return "", errors.Wrapf(err, "failed to read for the second time the reader %s/%s to generate a miniature", request.Owner, request.Id)
	}

	err = generateMiniature(request.Owner, request.Id, reader)
	return filename, err
}

func storeInArchive(request *StoreRequest) (string, bool, error) {
	key, err := repositoryPort.FindById(request.Owner, request.Id)
	if err == nil {
		return path.Base(key), false, nil
	}
	if err != nil && !errors.Is(err, NotFoundError) {
		return "", false, errors.Wrapf(err, "find existing location for the media")
	}

	content, err := request.Open()
	if err != nil {
		return "", false, errors.Wrapf(err, "couldn't archive the file")
	}
	defer content.Close()

	const dateFormatInFilename = "2006-01-02_15-04-05"
	cleanedFolderName := strings.Trim(request.FolderName, "/")
	key, err = storePort.Upload(DestructuredKey{
		Prefix: fmt.Sprintf("%s/%s/%s_%s", request.Owner, cleanedFolderName, request.DateTime.Format(dateFormatInFilename), request.SignatureSha256[:8]),
		Suffix: strings.ToLower(path.Ext(request.OriginalFilename)),
	}, content)
	if err != nil {
		return "", false, err
	}

	err = repositoryPort.AddLocation(request.Owner, request.Id, key)
	return path.Base(key), true, err
}

func generateMiniature(owner, mediaId string, reader io.Reader) error {
	content, mediaType, err := ResizerPort.ResizeImage(reader, MiniatureCachedWidth, false)
	if err != nil {
		return errors.Wrapf(err, "failed to generate minature")
	}

	cacheId := generateCacheId(owner, mediaId, MiniatureCachedWidth)
	return cachePort.Put(cacheId, mediaType, bytes.NewReader(content))
}
