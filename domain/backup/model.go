package backup

// MediaKey is the unique signature for a media
type MediaKey struct {
	Owner      string
	FolderName string
	Filename   string
}
