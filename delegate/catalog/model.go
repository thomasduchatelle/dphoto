package catalog

import (
	"fmt"
	"time"
)

// Album defines how medias are physically stored.
type Album struct {
	Name       string
	FolderName string // unique and immutable
	Start      time.Time
	End        time.Time
}

type CreateAlbum struct {
	Name             string
	Start            time.Time
	End              time.Time
	ForcedFolderName string
}

func (a *Album) String() string {
	const layout = "2006-01-02T03"
	return fmt.Sprintf("[%s-%s] %s (%s)", a.Start.Format(layout), a.End.Format(layout), a.FolderName, a.Name)
}

// IsEqual uses unique identifier to compare both albums
func (a *Album) IsEqual(other *Album) bool {
	return a.FolderName == other.FolderName
}

func newTimeRangeFromAlbum(album Album) timeRange {
	if album.Start.After(album.End) {
		panic("Album must end AFTER its start: " + album.String())
	}

	return timeRange{
		Start: album.Start,
		End:   album.End,
	}
}

type MediaType string
type MediaOrientation string

type MediaLocation struct {
	FolderName string
	Filename   string
}

type MediaSignature struct {
	SignatureSha256 string
	SignatureSize   int
}

type CreateMediaRequest struct {
	Location  MediaLocation
	Type      MediaType
	Details   MediaDetails
	Signature MediaSignature
}

type MediaMeta struct {
	Signature MediaSignature // Signature is the key used to get the image (or its location)
	Filename  string         // Filename original filename when image was uploaded
	Type      MediaType
	Details   MediaDetails
}

type MediaDetails struct {
	Width, Height             int
	DateTime                  time.Time
	Orientation               MediaOrientation
	Make                      string
	Model                     string
	GPSLatitude, GPSLongitude float64
}

type MovedMedia struct {
	Signature        MediaSignature
	SourceFolderName string
	TargetFolderName string
	Filename         string
}

type MediaPage struct {
	NextPage string // empty if no other pages
	Content  []*MediaMeta
}

type MediaSignatureAndLocation struct {
	Location  MediaLocation
	Signature MediaSignature
}

func (s MediaSignature) String() string {
	return fmt.Sprintf("Signature[%s - %d]", s.SignatureSha256, s.SignatureSize)
}

// AlbumStat has the counts of media in the album ; it's currently limited to total number because of the database.
type AlbumStat struct {
	Album      Album
	totalCount int
}

// TotalCount return the number of medias, no matter their type, in the album.
func (s *AlbumStat) TotalCount() int {
	return s.totalCount
}
