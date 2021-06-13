package ui

import (
	"fmt"
	"sync/atomic"
)

type InteractiveSession struct {
	actionsPort    InteractiveActionsPort
	catchError     atomic.Value
	renderer       *interactiveRender
	repositoryPort RecordRepositoryPort
	state          interactiveViewState
	userInputPort  UserInputPort
}

func NewInteractiveSession(actions InteractiveActionsPort, repositories ...RecordRepositoryPort) *InteractiveSession {
	return &InteractiveSession{
		actionsPort:    actions,
		renderer:       newInteractiveRender(),
		repositoryPort: NewRepositoryAggregator(repositories...),
		state: interactiveViewState{
			recordsState: recordsState{},
			Actions:      nil,
		},
		userInputPort: new(keyboardInteractionAdaptor),
	}
}

func (i *InteractiveSession) Start() error {
	records, err := i.repositoryPort.FindRecords()
	if err != nil {
		return err
	}
	i.state.Records = records
	i.state.PageSize = i.renderer.Height() - 15

	i.userInputPort.startListening()

	err, _ = i.catchError.Load().(error)
	return err
}

func (i *InteractiveSession) MoveDown() {
	i.state.Selected = (i.state.Selected + 1) % len(i.state.Records)
	if i.state.Selected >= i.state.FirstElement+i.state.PageSize || i.state.Selected < i.state.FirstElement {
		i.state.FirstElement = i.state.PageSize * (i.state.Selected / i.state.PageSize)
	}
	i.updateActions()
}

func (i *InteractiveSession) MoveUp() {
	i.state.Selected = (len(i.state.Records) + i.state.Selected - 1) % len(i.state.Records)
	if i.state.Selected >= i.state.FirstElement+i.state.PageSize || i.state.Selected < i.state.FirstElement {
		i.state.FirstElement = i.state.PageSize * (i.state.Selected / i.state.PageSize)
	}
	i.updateActions()
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
	if i.state.Records[i.state.Selected].Suggestion {
		actions = append(actions, "C: create")
	} else {
		actions = append(actions, "D: delete", "F: edit name", "E: edit dates")
	}

	i.state.Actions = actions
}
