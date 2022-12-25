package ui

import (
	"time"
)

type Period struct {
	Start, End time.Time
}

// SuggestionRecord is a record that will be displayed and handled in the UI. It can be an existing catalog.Album or a backup.FoundAlbum
type SuggestionRecord struct {
	FolderName   string // FolderName is a suggested name when Suggestion is true, not a unique key
	Name         string
	Start, End   time.Time
	Distribution map[string]int // Distribution is the number of media found for each day (format YYYY-MM-DD)
	Original     interface{}    // Original is used by adapter for targeted backup
}

// ExistingRecord is an album already existing
type ExistingRecord struct {
	FolderName    string // FolderName is a suggested name when Suggestion is true, not a unique key
	Name          string
	Start, End    time.Time
	Count         int
	ActivePeriods []Period
}

type Record struct {
	Indent               int    // Indent is to represent the list as a tree
	Suggestion           bool   // Suggestion is TRUE when it's a suggestion, not an existing album.
	FolderName           string // FolderName is a suggested name when Suggestion is true, not a unique key
	Name                 string
	Start, End           time.Time
	Count                int               // Count is the number of files relevant to the context (if in a tree branch)
	TotalCount           int               // TotalCount is the total number of file
	ParentExistingRecord *ExistingRecord   // ParentExistingRecord is the album if the suggestion is a child (used to limit the backup to a single album)
	SuggestionRecord     *SuggestionRecord // SuggestionRecord is the original when Suggestion is true (used for backup)
}

type RecordsState struct {
	Records      []*Record
	Rejected     int // Rejected is the number of file that has been rejected (ignored)
	Selected     int // Selected can be -1 to not highlight any line
	PageSize     int // PageSize can be 0 to disable pagination
	FirstElement int // FirstElement is the index of the first shown record ; default (or pagination disabled): 0
}

type InteractiveViewState struct {
	RecordsState
	Actions []string
}

// RecordCreation contains parameter to create a new album.
type RecordCreation struct {
	Owner      string
	FolderName string // FolderName might be empty to be generated
	Name       string
	Start, End time.Time
}

// SuggestionRecordRepositoryPort is the port providing data to the UI
type SuggestionRecordRepositoryPort interface {
	FindSuggestionRecords() []*SuggestionRecord
	Count() int
	Rejects() int
}

// ExistingRecordRepositoryPort is the port providing data to the UI
type ExistingRecordRepositoryPort interface {
	FindExistingRecords() ([]*ExistingRecord, error)
}

// NoopRepository is only used for Noop version
type NoopRepository struct{}

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

type BackupSuggestionPort interface {
	BackupSuggestion(record *SuggestionRecord, existing *ExistingRecord, listener InteractiveRendererPort) error
}

// InteractiveActionsPort are actions on 'SuggestionRecord.Suggestion = false' records
type InteractiveActionsPort interface {
	CreateAlbumPort
	RenameAlbumPort
	UpdateAlbumPort
	DeleteAlbumPort
	BackupSuggestionPort
}

// UserInputPort listens user input (keyboard) to interact with the session
type UserInputPort interface {
	StartListening()
}

// PrintReadTerminalPort is a port to print questions (simple strings), and read answers (strings as well)
type PrintReadTerminalPort interface {
	Print(question string)
	ReadAnswer() (string, error)
}

// InteractiveRendererPort is handling the rendering of an interactive session
type InteractiveRendererPort interface {
	PrintReadTerminalPort
	Render(state *InteractiveViewState) error
	Height() int
	// TakeOverScreen is clearing the screen and let another object handling the rendering
	TakeOverScreen()
}

// NewNoopRepository implements both SuggestionRecordRepositoryPort and ExistingRecordRepositoryPort but won't returns anything.
func NewNoopRepository() SuggestionRecordRepositoryPort {
	return new(NoopRepository)
}

func (r NoopRepository) FindSuggestionRecords() []*SuggestionRecord {
	return nil
}

func (r NoopRepository) Count() int {
	return 0
}

func (r NoopRepository) Rejects() int {
	return 0
}
