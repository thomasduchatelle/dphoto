package common

import (
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"strings"
)

func ConvertFolderNameForREST(folderName catalog.FolderName) string {
	return strings.TrimPrefix(folderName.String(), "/")
}
