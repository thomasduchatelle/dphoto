package ui

import (
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
				i.state.Selected = (i.state.Selected + 1) % len(i.state.Records)
				i.refresh()

			case keyboard.KeyArrowUp:
				i.state.Selected = (len(i.state.Records) + i.state.Selected - 1) % len(i.state.Records)
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

// refresh update available actions depending on selected line, and refresh the screen
func (i *InteractiveSession) refresh() {
	actions := []string{
		"ESC: exit",
		"N: new",
	}
	if i.state.Selected < 0 || i.state.Selected >= len(i.state.Records) {
		i.state.Selected = 0
	}
	if i.state.Records[i.state.Selected].Suggestion {
		actions = append(actions, "C: create")
	} else {
		actions = append(actions, "D: delete", "F: edit name", "E: edit dates")
	}

	i.state.Actions = actions
	err := i.renderer.Render(&i.state)
	if err != nil {
		i.catchError.Store(err)
	}
}
