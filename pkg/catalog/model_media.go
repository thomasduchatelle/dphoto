package catalog

import (
	"fmt"
	"time"
)

var (
	MediaNotFoundError = fmt.Errorf("media not found")
)

type MediaType string
type MediaOrientation string

type MediaSignature struct {
	SignatureSha256 string
	SignatureSize   int
}

func (s MediaSignature) String() string {
	return fmt.Sprintf("Signature[%s - %d]", s.SignatureSha256, s.SignatureSize)
}

// CreateMediaRequest is the request to add a new media to the  It's within an Owner context.
type CreateMediaRequest struct {
	Id         string         // Id is generated from its signature with GenerateMediaId(MediaSignature)
	Signature  MediaSignature // Signature is the business key of a media
	FolderName string         // FolderName is the name of the album the media is in
	Filename   string         // Filename is a user-friendly name that have the right extension.
	Type       MediaType
	Details    MediaDetails
}

// MediaMeta is an entry (read) of a media in the catalog
type MediaMeta struct {
	Id        string         // Id is the unique identifier to use across all domains
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

// MediaPage is the current page MediaMeta, and the token of the next page
type MediaPage struct {
	NextPage string // NextPage is empty if no other pages
	Content  []*MediaMeta
}
