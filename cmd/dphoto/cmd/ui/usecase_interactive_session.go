package ui

import (
	"fmt"
	"sync/atomic"
)

type InteractiveSession struct {
	actionsPort          InteractiveActionsPort
	catchError           atomic.Value
	existingRepository   ExistingRecordRepositoryPort
	owner                string
	renderer             InteractiveRendererPort
	state                InteractiveViewState
	suggestionRepository SuggestionRecordRepositoryPort
}

type recordNode struct {
	ExistingRecord *ExistingRecord // ExistingRecord is only present if record is not a suggestion
	record         *Record
	activeDays     map[string]interface{} // activeDays has the day ('YYYY-MM-dd' format) as key
	children       []*Record
}

func NewInteractiveSession(actions InteractiveActionsPort, existingRepository ExistingRecordRepositoryPort, suggestionRepository SuggestionRecordRepositoryPort, owner string) *InteractiveSession {
	return &InteractiveSession{
		actionsPort:          actions,
		renderer:             newInteractiveRender(),
		existingRepository:   existingRepository,
		owner:                owner,
		suggestionRepository: suggestionRepository,
		state: InteractiveViewState{
			RecordsState: RecordsState{
				Rejected: suggestionRepository.Rejects(),
			},
			Actions: nil,
		},
	}
}

func (i *InteractiveSession) Start() error {
	i.state.PageSize = i.renderer.Height() - 15
	err := i.reloadRecords()
	if err != nil {
		return err
	}
	i.updateActions()

	keyboardPort := keyboardInteractionAdaptor{
		session: i,
	}
	keyboardPort.StartListening()

	err, _ = i.catchError.Load().(error)
	return err
}

func (i *InteractiveSession) MoveDown() {
	if len(i.state.Records) > 0 {
		i.state.Selected = (i.state.Selected + 1) % len(i.state.Records)
		if i.state.Selected >= i.state.FirstElement+i.state.PageSize || i.state.Selected < i.state.FirstElement {
			i.state.FirstElement = i.state.PageSize * (i.state.Selected / i.state.PageSize)
		}
		i.updateActions()
	}
}

func (i *InteractiveSession) MoveUp() {
	if len(i.state.Records) > 0 {
		i.state.Selected = (len(i.state.Records) + i.state.Selected - 1) % len(i.state.Records)
		if i.state.Selected >= i.state.FirstElement+i.state.PageSize || i.state.Selected < i.state.FirstElement {
			i.state.FirstElement = i.state.PageSize * (i.state.Selected / i.state.PageSize)
		}
		i.updateActions()
	}
}

func (i *InteractiveSession) NextPage() {
	if i.state.PageSize < len(i.state.Records) {
		i.state.FirstElement += i.state.PageSize
		if i.state.FirstElement >= len(i.state.Records) {
			i.state.FirstElement = 0
		}
		i.state.Selected = i.state.FirstElement
		i.updateActions()
	}
}

func (i *InteractiveSession) PreviousPage() {
	if i.state.PageSize < len(i.state.Records) {
		i.state.FirstElement -= i.state.PageSize
		if i.state.FirstElement < 0 {
			i.state.FirstElement = i.state.PageSize * ((len(i.state.Records) - 1) / i.state.PageSize)
		}
		i.state.Selected = i.state.FirstElement
		i.updateActions()
	}
}

// Refresh updates the screen
func (i *InteractiveSession) Refresh() {
	i.must(i.renderer.Render(&i.state))
}

func (i *InteractiveSession) HasError() bool {
	_, hasError := i.catchError.Load().(error)
	return hasError
}

func (i *InteractiveSession) CreateFromSelectedSuggestion(owner string) {
	record := *i.state.Records[i.state.Selected]
	if record.Suggestion {
		ok, err := NewCreateAlbumForm(i.actionsPort, i.renderer).AlbumForm(owner, record)
		if err != nil {
			i.catchError.Store(err)
			return
		}

		if i.must(err) && ok {
			i.must(i.reloadRecords())
		}
	}
}

func (i *InteractiveSession) CreateNew(owner string) {
	_, err := NewCreateAlbumForm(i.actionsPort, i.renderer).AlbumForm(owner, Record{})
	if i.must(err) {
		i.must(i.reloadRecords())
	}
}

func (i *InteractiveSession) DeleteSelectedAlbum() {
	record := *i.state.Records[i.state.Selected]
	if !record.Suggestion {
		deleted, err := NewDeleteAlbumForm(i.actionsPort, i.renderer).DeleteAlbum(record)
		if i.must(err) && deleted {
			i.must(i.reloadRecords())
		}
	}
}

func (i *InteractiveSession) EditSelectedAlbumDates() {
	record := *i.state.Records[i.state.Selected]
	if !record.Suggestion {
		err := NewEditAlbumDateForm(i.actionsPort, i.renderer).EditAlbumDates(record)
		if i.must(err) {
			i.must(i.reloadRecords())
		}

	}
}

func (i *InteractiveSession) EditSelectedAlbumName() {
	record := *i.state.Records[i.state.Selected]
	if !record.Suggestion {
		err := NewRenameAlbumForm(i.actionsPort, i.renderer).EditAlbumName(record)
		if i.must(err) {
			i.must(i.reloadRecords())
		}
	}
}

func (i *InteractiveSession) BackupSelected() {
	record := *i.state.Records[i.state.Selected]
	if record.Suggestion {
		// take control of the screen
		if i.must(i.actionsPort.BackupSuggestion(record.SuggestionRecord, record.ParentExistingRecord, i.renderer)) {

			// hand over screen control
			i.must(i.reloadRecords())
		}
	}
}

func (i *InteractiveSession) reloadRecords() error {
	existing, err := i.existingRepository.FindExistingRecords()
	if err != nil {
		return err
	}

	suggestions := i.suggestionRepository.FindSuggestionRecords()

	i.state.Records = createFlattenTree(existing, suggestions)

	if i.state.Selected >= len(i.state.Records) {
		i.state.Selected = 0
		i.state.FirstElement = 0
	}

	return nil
}

func (i *InteractiveSession) updateActions() {
	actions := []string{
		"ESC: exit",
		"N: new",
	}
	if len(i.state.Records) > i.state.PageSize {
		incompletePage := 0
		if len(i.state.Records)%i.state.PageSize > 0 {
			incompletePage = 1
		}
		actions = append([]string{
			fmt.Sprintf("page %d/%d", 1+i.state.FirstElement/i.state.PageSize, incompletePage+len(i.state.Records)/i.state.PageSize),
			"DOWN: next page",
			"UP: previous page",
		}, actions...)
	}

	if len(i.state.Records) == 0 {
		// no specific action
	} else if i.state.Records[i.state.Selected].Suggestion {
		actions = append(actions, "C: create", "B: backup")
	} else {
		actions = append(actions, "DEL: delete", "E: edit name", "D: edit dates")
	}

	i.state.Actions = actions
}

// must return TRUE is there is no error ; otherwise catch the error and returns false.
func (i *InteractiveSession) must(err error) bool {
	if err != nil {
		i.catchError.Store(err)
		return false
	}

	return true
}
