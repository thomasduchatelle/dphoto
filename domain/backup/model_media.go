package backup

import (
	"fmt"
	"io"
	"time"
)

const (
	MediaTypeImage MediaType = "IMAGE"
	MediaTypeVideo MediaType = "VIDEO"
	MediaTypeOther MediaType = "OTHER"

	OrientationUpperLeft  ImageOrientation = "UPPER_LEFT"
	OrientationLowerRight ImageOrientation = "LOWER_RIGHT"
	OrientationUpperRight ImageOrientation = "UPPER_RIGHT"
	OrientationLowerLeft  ImageOrientation = "LOWER_LEFT"

	ProgressEventScanComplete   ProgressEventType = "scan-complete"
	ProgressEventAnalysed       ProgressEventType = "analysed"
	ProgressEventCatalogued     ProgressEventType = "catalogued"
	ProgressEventWrongAlbum     ProgressEventType = "wrong-album"
	ProgressEventAlreadyExists  ProgressEventType = "duplicate-catalog"
	ProgressEventDuplicate      ProgressEventType = "duplicate-unique"
	ProgressEventReadyForUpload ProgressEventType = "upload-ready"
	ProgressEventUploaded       ProgressEventType = "uploaded"
	ProgressEventAlbumCreated   ProgressEventType = "album-created"
)

// MediaType is photo or video
type MediaType string

// ImageOrientation is teh start point of stored data
type ImageOrientation string

// MediaPath is a breakdown of an absolute path, or URL, agnostic of its origin.
type MediaPath struct {
	ParentFullPath string // ParentFullPath is the absolute path of the media folder (URL = ParentFullPath + Filename)
	Root           string // Root is the path or URL representing the volume in which the media has been found. (URL = Root + Path + Filename)
	Path           string // Path is the directory path relative to Root: URL = Root + Path + Filename
	Filename       string // Filename does not contain any slash, and contains the extension
	ParentDir      string // ParentDir is the name of the directory: dirname(Root + Path)
}

// FoundMedia represents files found on the scanned volume
type FoundMedia interface {
	// MediaPath return breakdown of the absolute path of the media.
	MediaPath() MediaPath
	// ReadMedia reads content of the file ; it might not be optimised to call it several times (see VolumeToBackup)
	ReadMedia() (io.ReadCloser, error)
	// Size returns the size of the file
	Size() int

	String() string
}

// AnalysedMedia is a FoundMedia to which has been attached its type (photo / video) and other details usually found within the file.
type AnalysedMedia struct {
	FoundMedia FoundMedia    // FoundMedia is the reference of the file, implementation depends on the VolumeType
	Type       MediaType     // Type is 'photo' or 'video'
	Sha256Hash string        // Sha256Hash sha256 of the file
	Details    *MediaDetails // Details are data found within the file (location, date, ...)
}

// BackingUpMediaRequest is the requests that must be executed to back up the media
type BackingUpMediaRequest struct {
	AnalysedMedia *AnalysedMedia
	Id            string
	FolderName    string
}

// CatalogMediaRequest is the request passed to Archive domain
type CatalogMediaRequest struct {
	BackingUpMediaRequest *BackingUpMediaRequest
	ArchiveFilename       string // ArchiveFilename is a normalised named generated and used in archive.
}

// ClosableFoundMedia can be implemented alongside FoundMedia if the implementation requires to release resources once the media has been handled.
type ClosableFoundMedia interface {
	Close() error
}

type MediaDetails struct {
	Width, Height             int
	DateTime                  time.Time
	Orientation               ImageOrientation
	Make                      string
	Model                     string
	GPSLatitude, GPSLongitude float64
	Duration                  int64  // Duration is the length, in milliseconds, of a video
	VideoEncoding             string // VideoEncoding is the codec used to encode the video (ex: 'H264')
}

func (s *MediaDetails) String() string {
	return fmt.Sprintf("[Width=%d,Height=%d,DateTime=%s,Orientation=%s,Make=%s,Model=%s,GPSLatitude=%f,GPSLongitude=%f,Duration=%d,VideoEncoding=%s]", s.Width, s.Height, s.DateTime, s.Orientation, s.Make, s.Model, s.GPSLatitude, s.GPSLongitude, s.Duration, s.VideoEncoding)
}

// FullMediaSignature is the business key of the media, unique per user
type FullMediaSignature struct {
	Sha256 string
	Size   uint
}

func (s *FullMediaSignature) String() string {
	return fmt.Sprintf("%s (%s)", s.Sha256, byteCountIEC(int64(s.Size)))
}

type PostAnalyseFilter interface {
	// AcceptAnalysedMedia returns TRUE if the media should be backed-up.
	AcceptAnalysedMedia(media *AnalysedMedia, folderName string) bool
}

func byteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
