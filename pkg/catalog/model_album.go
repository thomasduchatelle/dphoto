package catalog

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

var (
	NotFoundError = errors.New("album hasn't been found")
	NotEmptyError = errors.New("album is not empty")
)

// Album defines how medias are physically re-grouped.
type Album struct {
	Owner      string    // Owner is a PK with FolderName
	Name       string    // Name for displaying purpose, not unique
	FolderName string    // FolderName is unique with Owner, and immutable
	Start      time.Time // Start is datetime inclusive
	End        time.Time // End is the datetime exclusive
	TotalCount int       // TotalCount is the number of media (of any type)
}

type CreateAlbum struct {
	Owner            string
	Name             string
	Start            time.Time
	End              time.Time
	ForcedFolderName string
}

type AlbumId struct {
	Owner      string
	FolderName string
}

func (a *Album) String() string {
	const layout = "2006-01-02T03"
	return fmt.Sprintf("[%s -> %s] %s (%s)", a.Start.Format(layout), a.End.Format(layout), a.FolderName, a.Name)
}

// IsEqual uses unique identifier to compare both albums
func (a *Album) IsEqual(other *Album) bool {
	return a.Owner == a.Owner && a.FolderName == other.FolderName
}

func (a *AlbumId) String() string {
	return fmt.Sprintf("%s/%s", a.Owner, a.FolderName)
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

func (c *CreateAlbum) String() string {
	const layout = "2006-01-02T03"
	return fmt.Sprintf("[%s -> %s] %s (%s/%s)", c.Start.Format(layout), c.End.Format(layout), c.Name, c.Owner, c.ForcedFolderName)
}

type PageRequest struct {
	Size     int64
	NextPage string
}
