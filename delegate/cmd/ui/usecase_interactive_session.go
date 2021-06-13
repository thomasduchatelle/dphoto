package ui

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"sync/atomic"
)

type InteractiveSession struct {
	actionsPort    InteractiveActionsPort
	catchError     atomic.Value
	done           chan struct{}
	renderer       *interactiveRender
	repositoryPort RecordRepositoryPort
	state          interactiveViewState
}

func NewInteractiveSession(actions InteractiveActionsPort, repositories ...RecordRepositoryPort) *InteractiveSession {
	return &InteractiveSession{
		actionsPort:    actions,
		done:           make(chan struct{}),
		renderer:       newInteractiveRender(),
		repositoryPort: NewRepositoryAggregator(repositories...),
		state: interactiveViewState{
			recordsState: recordsState{},
			Actions:      nil,
		},
	}
}

func (i *InteractiveSession) Start() error {
	records, err := i.repositoryPort.FindRecords()
	if err != nil {
		return err
	}
	i.state.Records = records
	i.state.PageSize = i.renderer.Height() - 15

	i.actionListener()

	err, _ = i.catchError.Load().(error)
	return err
}

func (i *InteractiveSession) actionListener() {
	go func() {
		defer close(i.done)
		i.refresh()

		for i.catchError.Load() == nil {
			ck, key, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}

			switch key {
			case keyboard.KeyEsc:
				return

			case keyboard.KeyArrowDown:
				i.MoveDown()
				i.refresh()

			case keyboard.KeyArrowUp:
				i.MoveUp()
				i.refresh()

			case keyboard.KeyPgdn:
				i.NextPage()
				i.refresh()

			case keyboard.KeyPgup:
				i.PreviousPage()
				i.refresh()

			default:
				switch ck {
				//case 'c':
				//	record := *i.state.Records[i.state.Selected]
				//	if record.Suggestion {
				//		i.screen.Keep(1)
				//		fmt.Println()
				//
				//		err = CreateAlbumForm(i.operations, record)
				//		i.state.Records = append(i.state.Records[:i.state.Selected], i.state.Records[i.state.Selected+1:]...)
				//		if err != nil {
				//			panic(err)
				//		}
				//		i.reloadExistingAlbum()
				//		i.refresh()
				//	}
				//
				//case 'n':
				//	i.screen.Keep(1)
				//	fmt.Println()
				//
				//	err = CreateAlbumForm(i.operations, Record{})
				//	if err != nil {
				//		panic(err)
				//	}
				//	i.reloadExistingAlbum()
				//	i.refresh()
				//
				//case 'd':
				//	record := *i.state.Records[i.state.Selected]
				//	if !record.Suggestion {
				//		i.screen.Keep(1)
				//		fmt.Println()
				//		err = DeleteAlbum(i.operations, record)
				//		if err != nil {
				//			panic(err)
				//		}
				//		i.reloadExistingAlbum()
				//		i.refresh()
				//	}
				//
				//case 'e':
				//	record := *i.state.Records[i.state.Selected]
				//	if !record.Suggestion {
				//		i.screen.Keep(1)
				//		fmt.Println()
				//		err = EditAlbumDates(i.operations, record)
				//		if err != nil {
				//			panic(err)
				//		}
				//		i.reloadExistingAlbum()
				//		i.refresh()
				//	}
				//
				//case 'f':
				//	record := *i.state.Records[i.state.Selected]
				//	if !record.Suggestion {
				//		i.screen.Keep(1)
				//		fmt.Println()
				//		err = EditAlbumName(i.operations, record)
				//		if err != nil {
				//			panic(err)
				//		}
				//		i.reloadExistingAlbum()
				//		i.refresh()
				//	}
				}

			}
		}

		i.refresh()
	}()

	<-i.done
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

// refresh update available actions depending on selected line, and refresh the screen
func (i *InteractiveSession) refresh() {
	err := i.renderer.Render(&i.state)
	if err != nil {
		i.catchError.Store(err)
	}
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
