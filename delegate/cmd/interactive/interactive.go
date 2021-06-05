package interactive

import (
	"duchatelle.io/dphoto/dphoto/cmd/screen"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/eiannone/keyboard"
	"github.com/logrusorgru/aurora/v3"
	"sort"
	"strings"
	"time"
)

type AlbumRecord struct {
	Suggestion bool   // Suggestion true means the album does not exists
	FolderName string // FolderName is a suggested name when Suggestion is true, not a unique key
	Name       string
	Start, End time.Time
	Count      uint
	Size       uint
}

type ui struct {
	operations  CatalogOperations
	screen      *screen.SimpleScreen
	suggestions []*AlbumRecord
	selected    int
	done        chan struct{}
	setActions  func(string)
}

type CatalogOperations interface {
	FindAllAlbumsWithStats() ([]*AlbumRecord, error)
	Create(createRequest AlbumRecord) error
	RenameAlbum(folderName, newName string, renameFolder bool) error
	UpdateAlbum(folderName string, start, end time.Time) error
	DeleteAlbum(folderName string) error
}

func StartInteractiveSession(suggestions []*AlbumRecord, operations CatalogOperations) {
	segment, setActions := screen.NewUpdatableSegment("ESC: exit")
	it := &ui{
		operations:  operations,
		done:        make(chan struct{}),
		setActions:  setActions,
		suggestions: suggestions,
	}

	it.screen = screen.NewSimpleScreen(screen.RenderingOptions{Width: 180}, it, segment)
	it.reloadExistingAlbum()

	it.start()
	return
}

func (u *ui) Content(options screen.RenderingOptions) string {
	const layout = "02/01/2006 (Mon)"

	table := simpletable.New()
	table.Header = &simpletable.Header{Cells: []*simpletable.Cell{
		{Text: ""},
		{Text: "Name"},
		{Text: "Folder Name"},
		{Text: "Start"},
		{Text: "End"},
		{Text: "Files (size)"},
	}}
	table.Body = &simpletable.Body{Cells: make([][]*simpletable.Cell, len(u.suggestions))}

	for idx, album := range u.suggestions {
		selected := u.selected == idx
		selectedMark := ""
		if selected {
			selectedMark = "*"
		}

		countContent := "-"
		if album.Count > 0 {
			countContent = fmt.Sprint(album.Count)
			if album.Size > 0 {
				countContent = fmt.Sprintf("%s (%s)", countContent, byteCountIEC(int64(album.Size)))
			}
		}

		table.Body.Cells[idx] = []*simpletable.Cell{
			{Text: u.applyStyle(selected, album.Suggestion, selectedMark)},
			{Text: u.applyStyleOnName(selected, album.Suggestion, album.Name)},
			{Text: u.applyStyle(selected, album.Suggestion, album.FolderName)},
			{Text: u.applyStyle(selected, album.Suggestion, album.Start.Format(layout))},
			{Text: u.applyStyle(selected, album.Suggestion, album.End.Format(layout))},
			{Text: u.applyStyle(selected, album.Suggestion, countContent)},
		}
	}

	return table.String()
}

func (u *ui) start() {
	go func() {
		defer close(u.done)

		for {
			ck, key, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}

			switch key {
			case keyboard.KeyEsc:
				return

			case keyboard.KeyArrowDown:
				u.selected = (u.selected + 1) % len(u.suggestions)
				u.refresh()

			case keyboard.KeyArrowUp:
				u.selected = (len(u.suggestions) + u.selected - 1) % len(u.suggestions)
				u.refresh()

			default:
				switch ck {
				case 'c':
					record := *u.suggestions[u.selected]
					if record.Suggestion {
						u.screen.Keep(1)
						fmt.Println()

						err = CreateAlbumForm(u.operations, record)
						u.suggestions = append(u.suggestions[:u.selected], u.suggestions[u.selected+1:]...)
						if err != nil {
							panic(err)
						}
						u.reloadExistingAlbum()
						u.refresh()
					}

				case 'n':
					u.screen.Keep(1)
					fmt.Println()

					err = CreateAlbumForm(u.operations, AlbumRecord{})
					if err != nil {
						panic(err)
					}
					u.reloadExistingAlbum()
					u.refresh()

				case 'd':
					record := *u.suggestions[u.selected]
					if !record.Suggestion {
						u.screen.Keep(1)
						fmt.Println()
						err = DeleteAlbum(u.operations, record)
						if err != nil {
							panic(err)
						}
						u.reloadExistingAlbum()
						u.refresh()
					}

				case 'e':
					record := *u.suggestions[u.selected]
					if !record.Suggestion {
						u.screen.Keep(1)
						fmt.Println()
						err = EditAlbumDates(u.operations, record)
						if err != nil {
							panic(err)
						}
						u.reloadExistingAlbum()
						u.refresh()
					}
				}

			}
		}
	}()

	u.screen.Refresh()
	<-u.done
}

func (u *ui) applyStyle(selected, suggestion bool, args interface{}) string {
	if selected {
		return aurora.Black(aurora.Bold(aurora.BgWhite(args))).String()
	}

	return fmt.Sprint(args)
}

func (u *ui) applyStyleOnName(selected, suggestion bool, args interface{}) string {
	switch {
	case selected && suggestion:
		return aurora.White(aurora.Bold(aurora.BgYellow(args))).String()

	case selected && !suggestion:
		return aurora.White(aurora.Bold(aurora.BgCyan(args))).String()

	case !selected && suggestion:
		return aurora.Yellow(args).String()

	case !selected && !suggestion:
		return aurora.Cyan(args).String()

	}

	return fmt.Sprint(args)
}

func (u *ui) reloadExistingAlbum() {
	albums, err := u.operations.FindAllAlbumsWithStats()
	if err != nil {
		panic(err)
	}

	suggestions := make([]*AlbumRecord, len(albums))
	for i, a := range albums {
		suggestions[i] = &AlbumRecord{
			Suggestion: false,
			FolderName: a.FolderName,
			Name:       a.Name,
			Start:      a.Start,
			End:        a.End,
			Count:      a.Count,
			Size:       a.Size,
		}
	}

	for _, a := range u.suggestions {
		if a.Suggestion {
			suggestions = append(suggestions, &AlbumRecord{
				Suggestion: true,
				FolderName: a.FolderName,
				Name:       a.Name,
				Start:      a.Start,
				End:        a.End,
				Count:      a.Count,
				Size:       0,
			})
		}
	}

	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Start == suggestions[j].Start {
			return suggestions[i].End.Before(suggestions[j].End)
		}

		return suggestions[i].Start.Before(suggestions[j].Start)
	})

	u.suggestions = suggestions
}

// refresh update available actions depending on selected line, and refresh the screen
func (u *ui) refresh() {
	actions := []string{
		"ESC: exit",
		"n: new",
	}
	if u.selected < 0 || u.selected >= len(u.suggestions) {
		u.selected = 0
	}
	if u.suggestions[u.selected].Suggestion {
		actions = append(actions, "c: create")
	} else {
		actions = append(actions, "d: delete", "n: edit name", "e: edit dates")
	}
	u.setActions(strings.Join(actions, " ; "))

	u.screen.Refresh()
}

func byteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
