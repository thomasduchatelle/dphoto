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

type CreateAlbumPort interface {
	Create(createRequest RecordCreation) error
}

type RenameAlbumPort interface {
	RenameAlbum(folderName, newName string, renameFolder bool) error
}

type UpdateAlbumPort interface {
	UpdateAlbum(folderName string, start, end time.Time) error
}

type DeleteAlbumPort interface {
	DeleteAlbum(folderName string) error
}

// InteractiveActionsPort are actions on 'Record.Suggestion = false' records
type InteractiveActionsPort interface {
	CreateAlbumPort
	RenameAlbumPort
	UpdateAlbumPort
	DeleteAlbumPort
}

// UserInputPort listens user input (keyboard) to interact with the session
type UserInputPort interface {
	startListening()
}

// PrintReadTerminalPort is a port to print questions (simple strings), and read answers (strings as well)
type PrintReadTerminalPort interface {
	Print(question string)
	ReadAnswer() (string, error)
}

type recordsState struct {
	Records      []*Record
	Selected     int // Selected can be -1 to not highlight any line
	PageSize     int
	FirstElement int
}

type interactiveViewState struct {
	recordsState
	Actions []string
}
