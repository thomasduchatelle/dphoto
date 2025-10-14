package common

import (
	"strings"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func ConvertFolderNameForREST(folderName catalog.FolderName) string {
	return strings.TrimPrefix(folderName.String(), "/")
}
