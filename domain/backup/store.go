package backup

import (
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"io"
)

// Note: new backup process: listing -> scan (hash + filter out duplicated) -> upload
// 'listing' is without any filter, media content can be cached in memory
// 'scan' compute the hash and filter out known medias
// 'upload' uploads the files, and inserts their metadata in the catalog

type MediaId struct {
	Owner string
	Hash  string
	Size  string
}

// MediaPath is a breakdown of an absolute path, or URL, agnostic of its origin.
type MediaPath struct {
	ParentFullPath string // ParentFullPath is the path that can be used to create a sub-volume that only contains sibling medias of this one
	Root           string // Root is the path or URL representing the volume in which the media has been found.
	Path           string // Path is the path between Root and Filename (ie: Root + Path + Filename would be the absolute URL)
	Filename       string // Filename does not contain any slash, and contains the extension.
	ParentDir      string // ParentDir is the name of the media folder ; it might be from the Path or from the Root
}

type FoundMedia interface {
	// MediaPath return breakdown of the absolute path of the media.
	MediaPath() MediaPath
	// ReadMedia reads content of the file ; it might not be optimised to call it several times (see VolumeToBackup)
	ReadMedia() (io.ReadCloser, error)
}

type Volume interface {
	// UniqueId represents a location unique for the computer
	UniqueId() string

	// Find lists the medias available in the Volume
	Find() ([]FoundMedia, error)
}

type ScanOptions struct {
	Listeners   []interface{}
	SkipRejects bool // SkipRejects mode will report any analysis error, or missing timestamp, and continue.
}

type Options struct {
	PostAnalyseFilter backupmodel.PostAnalyseFilter
	Listener          interface{}
}

func Backup(source Volume) {

}

func RecordAlbumChange(id MediaId, targetAlbum string) error {
	panic("not implemented")
}

func ScanVolume(volume backupmodel.VolumeToBackup, options ScanOptions) ([]*backupmodel.ScannedFolder, []backupmodel.FoundMedia, error) {
	panic("not implemented")
}

func StartBackupRunner(owner string, volume backupmodel.VolumeToBackup, options Options) (backupmodel.BackupReport, error) {
	panic("not implemented")
}
