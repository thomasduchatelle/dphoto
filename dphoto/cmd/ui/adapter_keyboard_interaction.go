package ui

import (
	"github.com/eiannone/keyboard"
)

type keyboardInteractionAdaptor struct {
	session *InteractiveSession
}

// StartListening listens keyboard events until session has an error or ESC key is hit.
func (a *keyboardInteractionAdaptor) StartListening() {
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

		case keyboard.KeyDelete:
			a.session.DeleteSelectedAlbum()
			a.session.Refresh()

		default:
			switch ck {
			case 'c':
				a.session.CreateFromSelectedSuggestion(a.session.owner)
				a.session.Refresh()

			case 'n':
				a.session.CreateNew(a.session.owner)
				a.session.Refresh()

			case 'd':
				a.session.EditSelectedAlbumDates()
				a.session.Refresh()

			case 'e':
				a.session.EditSelectedAlbumName()
				a.session.Refresh()

			case 'b':
				a.session.BackupSelected()
				a.session.Refresh()
			}

		}
	}

	a.session.Refresh()
}
