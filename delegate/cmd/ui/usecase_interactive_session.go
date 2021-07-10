package ui

import (
	"fmt"
	"sync/atomic"
)

type InteractiveSession struct {
	actionsPort    InteractiveActionsPort
	catchError     atomic.Value
	removedRecords map[Record]interface{}
	renderer       *interactiveRender
	repositoryPort RecordRepositoryPort
	state          interactiveViewState
}

func NewInteractiveSession(actions InteractiveActionsPort, repositories ...RecordRepositoryPort) *InteractiveSession {
	return &InteractiveSession{
		actionsPort:    actions,
		removedRecords: make(map[Record]interface{}),
		renderer:       newInteractiveRender(),
		repositoryPort: NewRepositoryAggregator(repositories...),
		state: interactiveViewState{
			recordsState: recordsState{},
			Actions:      nil,
		},
	}
}

func (i *InteractiveSession) Start() error {
	i.state.PageSize = i.renderer.Height() - 15
	i.reloadRecords()
	if i.catchError.Load() != nil {
		return i.catchError.Load().(error)
	}
	i.updateActions()

	keyboardPort := keyboardInteractionAdaptor{
		session: i,
	}
	keyboardPort.startListening()

	err, _ := i.catchError.Load().(error)
	return err
}

func (i *InteractiveSession) reloadRecords() {
	records, err := i.repositoryPort.FindRecords()
	if err != nil {
		i.catchError.Store(err)
		return
	}

	i.state.Records = i.state.Records[:0]
	for _, record := range records {
		if _, present := i.removedRecords[*record]; !present {
			i.state.Records = append(i.state.Records, record)
		}
	}

	if i.state.Selected >= len(i.state.Records) {
		// safeguard that won't happen unless anther process deleted some albums
		i.state.Selected = 0
		i.state.FirstElement = 0
	}
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
	err := i.renderer.Render(&i.state)
	if err != nil {
		i.catchError.Store(err)
	}
}

func (i *InteractiveSession) HasError() bool {
	_, hasError := i.catchError.Load().(error)
	return hasError
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
		actions = append(actions, "C: create")
	} else {
		actions = append(actions, "DEL: delete", "E: edit name", "D: edit dates")
	}

	i.state.Actions = actions
}

func (i *InteractiveSession) CreateFromSelectedSuggestion() {
	record := *i.state.Records[i.state.Selected]
	if record.Suggestion {
		ok, err := NewCreateAlbumForm(i.actionsPort, i.renderer).AlbumForm(record)
		if err != nil {
			i.catchError.Store(err)
			return
		}

		if ok {
			i.removedRecords[record] = nil
			i.reloadRecords()
		}
	}
}

func (i *InteractiveSession) CreateNew() {
	_, err := NewCreateAlbumForm(i.actionsPort, i.renderer).AlbumForm(Record{})
	if err != nil {
		i.catchError.Store(err)
		return
	}

	i.reloadRecords()
}

func (i *InteractiveSession) DeleteSelectedAlbum() {
	record := *i.state.Records[i.state.Selected]
	if !record.Suggestion {
		deleted, err := NewDeleteAlbumForm(i.actionsPort, i.renderer).DeleteAlbum(record)
		if err != nil {
			i.catchError.Store(err)
			return
		}

		if deleted {
			i.state.Records = append(i.state.Records[:i.state.Selected], i.state.Records[i.state.Selected+1:]...)
		}
	}
}

func (i *InteractiveSession) EditSelectedAlbumDates() {
	record := *i.state.Records[i.state.Selected]
	if !record.Suggestion {
		err := NewEditAlbumDateForm(i.actionsPort, i.renderer).EditAlbumDates(record)
		if err != nil {
			i.catchError.Store(err)
			return
		}

		i.reloadRecords()
	}
}

func (i *InteractiveSession) EditSelectedAlbumName() {
	record := *i.state.Records[i.state.Selected]
	if !record.Suggestion {
		err := NewRenameAlbumForm(i.actionsPort, i.renderer).EditAlbumName(record)
		if err != nil {
			i.catchError.Store(err)
			return
		}

		i.reloadRecords()
	}
}
