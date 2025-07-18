package catalog

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"regexp"
	"strings"
	"time"
)

var (
	AlbumNotFoundErr = errors.New("album hasn't been found")

	EmptyFolderNameError = errors.New("folderName is mandatory and must be not empty")
)

// Album is a logical grouping of medias ; also used to physically store media next to each others.
type Album struct {
	AlbumId
	Name  string    // Name for displaying purpose, not unique
	Start time.Time // Start is datetime inclusive
	End   time.Time // End is the datetime exclusive
}

// IsEqual uses unique identifier to compare both albums
func (a Album) IsEqual(other *Album) bool {
	return a.AlbumId.IsEqual(other.AlbumId)
}

func (a Album) String() string {
	const layout = "2006-01-02T03"
	return fmt.Sprintf("%s '%s' [%s -> %s]", a.AlbumId, a.Name, a.Start.Format(layout), a.End.Format(layout))
}

type AlbumId struct {
	Owner      ownermodel.Owner
	FolderName FolderName
}

// IsEqual uses unique identifier to compare both albums
func (a AlbumId) IsEqual(other AlbumId) bool {
	return a.Owner == a.Owner && a.FolderName == other.FolderName
}

func (a AlbumId) IsValid() error {
	ownerErr := a.Owner.IsValid()
	folderNameErr := a.FolderName.IsValid()

	switch {
	case ownerErr != nil && folderNameErr != nil:
		return errors.Errorf("both owner [%s} and folderName [%s] must be valid", ownerErr.Error(), folderNameErr.Error())
	case ownerErr != nil:
		return ownerErr
	default:
		return folderNameErr
	}
}

func (a AlbumId) String() string {
	return fmt.Sprintf("%s%s", a.Owner, a.FolderName)
}

// NewAlbumIdFromStrings creates an AlbumId from 2 strings ; it doesn't guaranty its validity, use AlbumId.IsValid to check if any error.
func NewAlbumIdFromStrings(owner, folderName string) AlbumId {
	return AlbumId{Owner: ownermodel.Owner(owner), FolderName: NewFolderName(folderName)}
}

// FolderName is a normalised ID unique per Owner
type FolderName string

func (n FolderName) String() string {
	return string(n)
}

func (n FolderName) IsValid() error {
	if n == "" {
		return EmptyFolderNameError
	}

	return nil
}

// NewFolderName creates a FolderName with a normalised value ; it can still be invalid (empty)
func NewFolderName(name string) FolderName {
	nonAlphaNumeric := regexp.MustCompile("[^A-Za-z0-9-]+")
	return FolderName("/" + strings.Trim(nonAlphaNumeric.ReplaceAllString(name, "_"), "_"))
}

func generateFolderName(name string, start time.Time) FolderName {
	return NewFolderName(fmt.Sprintf("%s_%s", start.Format("2006-01"), name))
}

type PageRequest struct {
	Size     int64
	NextPage string
}
