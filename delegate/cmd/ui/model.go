package ui

import "time"

// Record is a record that will be displayed and handled in the UI. It can be an existing catalog.Album or a backup.FoundAlbum
type Record struct {
	Suggestion bool   // Suggestion true means the album does not exists
	FolderName string // FolderName is a suggested name when Suggestion is true, not a unique key
	Name       string
	Start, End time.Time
	Count      uint
}

// RecordCreation contains parameter to create a new album.
type RecordCreation struct {
	FolderName string // FolderName might be empty to be generated
	Name       string
	Start, End time.Time
}

// RecordRepositoryPort is the port providing data to the UI
type RecordRepositoryPort interface {
	FindRecords() ([]*Record, error)
}

// InteractiveActionsPort are actions on 'Record.Suggestion = false' records
type InteractiveActionsPort interface {
	Create(createRequest RecordCreation) error
	RenameAlbum(folderName, newName string, renameFolder bool) error
	UpdateAlbum(folderName string, start, end time.Time) error
	DeleteAlbum(folderName string) error
}

type recordsState struct {
	Records  []*Record
	Selected int // Selected can be -1 to not highlight any line
}

type interactiveViewState struct {
	recordsState
	Actions []string
}
