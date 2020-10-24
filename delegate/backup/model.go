package backup

import (
	"path"
	"time"
)

type MediaType string

// start point of stored data is...
type ImageOrientation string

const (
	IMAGE MediaType = "IMAGE"
	VIDEO MediaType = "VIDEO"
	OTHER MediaType = "OTHER"

	UPPER_LEFT  ImageOrientation = "UPPER_LEFT"
	LOWER_RIGHT ImageOrientation = "LOWER_RIGHT"
	UPPER_RIGHT ImageOrientation = "UPPER_RIGHT"
	LOWER_LEFT  ImageOrientation = "LOWER_LEFT"
)

type RemovableVolume struct {
	UniqueId   string
	MountPaths []string
}

type VolumeMetadata struct {
	UniqueId   string
	Name       string
	AutoBackup bool
}

type MediaDetails struct {
	Width, Height             int
	DateTime                  time.Time
	Orientation               ImageOrientation
	Make                      string
	Model                     string
	GPSLatitude, GPSLongitude float64
}

type SimpleMediaSignature struct {
	RelativePath string
	Size         int
}

type FullMediaSignature struct {
	Sha256 string
	Size   int64
}

type FoundMedia struct {
	Type              MediaType
	LocalAbsolutePath string
	SimpleSignature   SimpleMediaSignature
}

type LocalMedia struct {
	Type              MediaType
	LocalAbsolutePath string
	Signature         FullMediaSignature
	Details           *MediaDetails
}

type Report struct {
	Success bool
	Errors  []error
	Counter *Counter
}

func NewFoundMedia(mediaType MediaType, volumeMount string, relativePath string, size int) FoundMedia {
	return FoundMedia{
		Type:              mediaType,
		LocalAbsolutePath: path.Clean(path.Join(volumeMount, relativePath)),
		SimpleSignature: SimpleMediaSignature{
			RelativePath: relativePath,
			Size:         size,
		},
	}
}
