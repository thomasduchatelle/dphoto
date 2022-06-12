package catalog

import (
	"fmt"
	"time"
)

var (
	MediaNotFoundError = fmt.Errorf("media not found")
)

// Album defines how medias are physically stored.
type Album struct {
	Owner      string // Owner is a PK with FolderName
	Name       string
	FolderName string // FolderName is unique with Owner, and immutable
	Start      time.Time
	End        time.Time
}

type CreateAlbum struct {
	Owner            string
	Name             string
	Start            time.Time
	End              time.Time
	ForcedFolderName string
}

func (a *Album) String() string {
	const layout = "2006-01-02T03"
	return fmt.Sprintf("[%s -> %s] %s (%s)", a.Start.Format(layout), a.End.Format(layout), a.FolderName, a.Name)
}

// IsEqual uses unique identifier to compare both albums
func (a *Album) IsEqual(other *Album) bool {
	return a.Owner == a.Owner && a.FolderName == other.FolderName
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

// CreateMediaRequest is the request to add a new media to the  It's within an Owner context.
type CreateMediaRequest struct {
	Location  MediaLocation
	Type      MediaType
	Details   MediaDetails
	Signature MediaSignature
}

// MediaMeta is an entry (read) of a media in the catalog
type MediaMeta struct {
	Signature MediaSignature // Signature is the key used to get the image (or its location)
	Filename  string         // Filename original filename when image was uploaded
	Type      MediaType
	Details   MediaDetails
}

// MediaDetails are extracted from the metadata within photos and videos and stored as it.
type MediaDetails struct {
	Width, Height             int
	DateTime                  time.Time
	Orientation               MediaOrientation
	Make                      string
	Model                     string
	GPSLatitude, GPSLongitude float64
	Duration                  int64  // Duration is the length, in milliseconds, of a video
	VideoEncoding             string // VideoEncoding is the codec used to encode the video (ex: 'H264')
}

// MovedMedia is a record of a media that will be, or have been, physically moved
type MovedMedia struct {
	Signature        MediaSignature
	SourceFolderName string
	SourceFilename   string
	TargetFolderName string
	TargetFilename   string
}

// MediaPage is the current page MediaMeta, and the token of the next page
type MediaPage struct {
	NextPage string // NextPage is empty if no other pages
	Content  []*MediaMeta
}

type MediaSignatureAndLocation struct {
	Location  MediaLocation
	Signature MediaSignature
}

func (s MediaSignature) String() string {
	return fmt.Sprintf("Signature[%s - %d]", s.SignatureSha256, s.SignatureSize)
}

// AlbumStat has the counts of media on the album ; it's currently limited to total number because of the database.
type AlbumStat struct {
	Album      Album
	TotalCount int
}

type PageRequest struct {
	Size     int64  // defaulted to 50 if not defined
	NextPage string // empty for the first page
}

type MoveTransaction struct {
	TransactionId string
	Count         int // Count is the number of medias to be moved as part of this transaction
}

type MoveMediaOperator interface {
	// Move must perform the physical move of the file to a different directory ; return the final name if it has been changed
	Move(source, dest MediaLocation) (string, error)

	// UpdateStatus informs of the global status of the move operation
	UpdateStatus(done, total int) error
	// Continue requests if the operation should continue or be interrupted
	Continue() bool
}
