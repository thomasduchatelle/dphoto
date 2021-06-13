package ui

import "github.com/eiannone/keyboard"

type keyboardInteractionAdaptor struct {
	session *InteractiveSession
}

// startListening listens keyboard events until session has an error or ESC key is hit.
func (a *keyboardInteractionAdaptor) startListening() {
	a.session.Refresh()

	for !a.session.HasError() {
		ck, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		switch key {
		case keyboard.KeyEsc:
			return

		case keyboard.KeyArrowDown:
			a.session.MoveDown()
			a.session.Refresh()

		case keyboard.KeyArrowUp:
			a.session.MoveUp()
			a.session.Refresh()

		case keyboard.KeyPgdn:
			a.session.NextPage()
			a.session.Refresh()

		case keyboard.KeyPgup:
			a.session.PreviousPage()
			a.session.Refresh()

		default:
			switch ck {
			//case 'c':
			//	record := *a.state.Records[a.state.Selected]
			//	if record.Suggestion {
			//		a.screen.Keep(1)
			//		fmt.Println()
			//
			//		err = CreateAlbumForm(a.operations, record)
			//		a.state.Records = append(a.state.Records[:a.state.Selected], a.state.Records[a.state.Selected+1:]...)
			//		if err != nil {
			//			panic(err)
			//		}
			//		a.reloadExistingAlbum()
			//		a.Refresh()
			//	}
			//
			//case 'n':
			//	a.screen.Keep(1)
			//	fmt.Println()
			//
			//	err = CreateAlbumForm(a.operations, Record{})
			//	if err != nil {
			//		panic(err)
			//	}
			//	a.reloadExistingAlbum()
			//	a.Refresh()
			//
			//case 'd':
			//	record := *a.state.Records[a.state.Selected]
			//	if !record.Suggestion {
			//		a.screen.Keep(1)
			//		fmt.Println()
			//		err = DeleteAlbum(a.operations, record)
			//		if err != nil {
			//			panic(err)
			//		}
			//		a.reloadExistingAlbum()
			//		a.Refresh()
			//	}
			//
			//case 'e':
			//	record := *a.state.Records[a.state.Selected]
			//	if !record.Suggestion {
			//		a.screen.Keep(1)
			//		fmt.Println()
			//		err = EditAlbumDates(a.operations, record)
			//		if err != nil {
			//			panic(err)
			//		}
			//		a.reloadExistingAlbum()
			//		a.Refresh()
			//	}
			//
			//case 'f':
			//	record := *a.state.Records[a.state.Selected]
			//	if !record.Suggestion {
			//		a.screen.Keep(1)
			//		fmt.Println()
			//		err = EditAlbumName(a.operations, record)
			//		if err != nil {
			//			panic(err)
			//		}
			//		a.reloadExistingAlbum()
			//		a.Refresh()
			//	}
			}

		}
	}

	a.session.Refresh()
}
