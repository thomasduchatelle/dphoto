// Package model is mostly to break cyclic dependencies between backup package and runner sub-package (which is type-safe).
package model

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

	VolumeTypeFileSystem VolumeType = "filesystem" // Mounted folder
	VolumeTypeS3         VolumeType = "s3"         // Storage in S3
	VolumeTypeMtp        VolumeType = "mtp"        // MTP (Android drive)

	ProgressEventScanComplete        ProgressEventType = "scan-complete"
	ProgressEventDownloaded          ProgressEventType = "downloaded"
	ProgressEventSkipped             ProgressEventType = "skipped-before-download"
	ProgressEventSkippedAfterAnalyse ProgressEventType = "skipped-after-analyse"
	ProgressEventAnalysed            ProgressEventType = "analysed"
	ProgressEventUploaded            ProgressEventType = "uploaded"
	ProgressEventAlbumCreated        ProgressEventType = "album-created"
)

// MediaType is photo or video
type MediaType string

// ImageOrientation is teh start point of stored data
type ImageOrientation string

// VolumeType is used to choose the drive to scan and read medias from the volume
type VolumeType string

// VolumeToBackup represents a location to backup.
type VolumeToBackup struct {
	UniqueId string     // UniqueId represents a location unique for the computer
	Type     VolumeType // Type is used to determine the implementation to scan and read medias
	Path     string     // Path is the absolute path or URL of the media
	Local    bool       // Local is true when files can be analysed (hash, EXIF) directly ; otherwise, they need to be copied locally first
}

type VolumeMetadata struct {
	UniqueId   string
	Name       string
	AutoBackup bool
}

type FoundMedia interface {
	// Filename returns the original filename, used to determine the type of the media
	Filename() string
	// LastModificationDate returns the dte the physical file has been last updated
	LastModificationDate() time.Time
	// SimpleSignature gets a key, unique for the volume
	SimpleSignature() *SimpleMediaSignature
	// ReadMedia reads content of the file ; it might not be optimised to call it several times (see VolumeToBackup)
	ReadMedia() (io.Reader, error)
}

// FoundMediaWithHash can be implemented along side FoundMedia if the implementation can compute sha256 hash on the fly
type FoundMediaWithHash interface {
	Sha256Hash() string
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

// SimpleMediaSignature is unique only for a single volume, and only for a certain time (i.e.: filename)
type SimpleMediaSignature struct {
	RelativePath string
	Size         uint
}

// FullMediaSignature is the business key of the media, unique per user
type FullMediaSignature struct {
	Sha256 string
	Size   uint
}

type AnalysedMedia struct {
	FoundMedia FoundMedia          // FoundMedia is the reference of the file, implementation depends on the VolumeType
	Type       MediaType           // Photo or Video
	Signature  *FullMediaSignature // Business key of the media
	Details    *MediaDetails       // Details are data found within the file (location, date, ...)
}

type ProgressEventType string

type ProgressEvent struct {
	Type      ProgressEventType // Type defines what's count, and size are about ; some might not be used.
	Count     uint              // Count is the number of media
	Size      uint              // Size is the sum of the size of the media concerned by this event
	Album     string            // Album is the folder name of the medias concerned by this event
	MediaType MediaType         // MediaType is the type of media ; only mandatory with 'uploaded' event
}

type DetailsReaderOptions struct {
	Fast bool // Fast determine if the reader should extract all the details it can, or should only focus on the DateTime (and other it could get in the same time)
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

func (s *FullMediaSignature) String() string {
	return fmt.Sprintf("%s (%s)", s.Sha256, byteCountIEC(int64(s.Size)))
}

func (s *SimpleMediaSignature) String() string {
	return fmt.Sprintf("%s (%s)", s.RelativePath, byteCountIEC(int64(s.Size)))
}

func (s *MediaDetails) String() string {
	return fmt.Sprintf("[Width=%d,Height=%d,DateTime=%s,Orientation=%s,Make=%s,Model=%s,GPSLatitude=%f,GPSLongitude=%f,Duration=%d,VideoEncoding=%s]", s.Width, s.Height, s.DateTime, s.Orientation, s.Make, s.Model, s.GPSLatitude, s.GPSLongitude, s.Duration, s.VideoEncoding)
}
