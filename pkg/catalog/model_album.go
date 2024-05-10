package catalog

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
	"time"
)

var (
	AlbumNotFoundError   = errors.New("album hasn't been found")
	EmptyOwnerError      = errors.New("owner is mandatory and must be not empty")
	EmptyFolderNameError = errors.New("folderName is mandatory and must be not empty")
)

// Album is a logical grouping of medias ; also used to physically store media next to each others.
type Album struct { // TODO is the total count appropriate on this object ??
	AlbumId
	Name       string    // Name for displaying purpose, not unique
	Start      time.Time // Start is datetime inclusive
	End        time.Time // End is the datetime exclusive
	TotalCount int       // TotalCount is the number of media (of any type)
}

// IsEqual uses unique identifier to compare both albums
func (a Album) IsEqual(other *Album) bool {
	return a.AlbumId.IsEqual(other.AlbumId)
}

func (a Album) String() string {
	const layout = "2006-01-02T03"
	return fmt.Sprintf("[%s -> %s] %s (%s)", a.Start.Format(layout), a.End.Format(layout), a.FolderName, a.Name)
}

type AlbumId struct {
	Owner      Owner
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
	return AlbumId{Owner: Owner(owner), FolderName: NewFolderName(folderName)}
}

// Owner is a non-empty ID
type Owner string

func (o Owner) IsValid() error {
	if o == "" {
		return EmptyOwnerError
	}

	return nil
}

func (o Owner) String() string {
	return string(o)
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

func NewTimeRangeFromAlbum(album Album) TimeRange {
	if album.Start.After(album.End) {
		panic("Album must end AFTER its start: " + album.String())
	}

	return TimeRange{
		Start: album.Start,
		End:   album.End,
	}
}

type MediaId string

type PageRequest struct {
	Size     int64
	NextPage string
}
