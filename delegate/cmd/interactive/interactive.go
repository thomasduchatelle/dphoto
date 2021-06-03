package interactive

import (
	"bufio"
	"duchatelle.io/dphoto/dphoto/cmd/screen"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/eiannone/keyboard"
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	"os"
	"strings"
	"time"
)

const layout = "02/01/2006"

type AlbumRecord struct {
	Suggestion bool   // Suggestion true means the album does not exists
	FolderName string // FolderName is a suggested name when Suggestion is true, not a unique key
	Name       string
	Start, End time.Time
	Count      uint
	Size       uint
}

type List struct {
	suggestions []*AlbumRecord
	selected    int
	done        chan struct{}
	screen      *screen.SimpleScreen
	operations  CatalogOperations
}

type CatalogOperations interface {
	FindAllAlbumsWithStats() ([]*AlbumRecord, error)
	Create(createRequest AlbumRecord) error
	RenameAlbum(folderName, newName string, renameFolder bool) error
	UpdateAlbum(folderName string, start, end time.Time) error
}

func StartInteractiveSession(suggestions []*AlbumRecord, operations CatalogOperations) {
	it := &List{
		suggestions: suggestions,
		done:        make(chan struct{}),
		operations:  operations,
	}
	it.screen = screen.NewSimpleScreen(screen.RenderingOptions{Width: 180}, it, screen.NewConstantSegment("C: create new album ; N: update name ; D: update dates"))

	it.start()
	return
}

func (l *List) Content(options screen.RenderingOptions) string {
	table := simpletable.New()
	table.Header = &simpletable.Header{Cells: []*simpletable.Cell{
		{Text: ""},
		{Text: "Name"},
		{Text: "Start"},
		{Text: "End"},
	}}
	table.Body = &simpletable.Body{Cells: make([][]*simpletable.Cell, len(l.suggestions))}

	for idx, album := range l.suggestions {
		selected := l.selected == idx
		selectedMark := ""
		if selected {
			selectedMark = "*"
		}

		table.Body.Cells[idx] = []*simpletable.Cell{
			{Text: l.applyStyle(selected, album.Suggestion, selectedMark)},
			{Text: l.applyStyleOnName(selected, album.Suggestion, album.Name)},
			{Text: l.applyStyle(selected, album.Suggestion, album.Start.Format(layout))},
			{Text: l.applyStyle(selected, album.Suggestion, album.End.Format(layout))},
		}
	}

	return table.String()
}

func (l *List) start() {
	go func() {
		defer close(l.done)

		for {
			ck, key, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}

			switch key {
			case keyboard.KeyEsc:
				return

			case keyboard.KeyArrowDown:
				l.selected = (l.selected + 1) % len(l.suggestions)

			case keyboard.KeyArrowUp:
				l.selected = (len(l.suggestions) + l.selected - 1) % len(l.suggestions)

			default:
				switch ck {
				case 'c':
					_ = keyboard.Close()

					l.screen.Clear()
					err = CreateAlbumForm(l.operations, *l.suggestions[l.selected])
					if err != nil {
						panic(err)
					}
					l.screen.Refresh()
				}

			}

			l.screen.Refresh()
		}
	}()

	l.screen.Refresh()
	<-l.done
}

func (l *List) applyStyle(selected, suggestion bool, args interface{}) string {
	if selected {
		return aurora.Black(aurora.Bold(aurora.BgWhite(args))).String()
	}

	return fmt.Sprint(args)
}

func (l *List) applyStyleOnName(selected, suggestion bool, args interface{}) string {
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

func CreateAlbumForm(operations CatalogOperations, record AlbumRecord) error {
	creation := AlbumRecord{}
	ok := true

	creation.Name, ok = scanString("Name of the album", record.Name)
	if !ok {
		return nil
	}

	creation.FolderName, _ = scanString("Folder name (leave blank for automatically generated)", "")

	creation.Start, ok = scanDate("Start date", record.Start)
	if !ok {
		return nil
	}

	creation.End, ok = scanDate("End date", record.End)
	if !ok {
		return nil
	}

	return operations.Create(creation)
}

func scanString(label string, defaultValue string) (string, bool) {
	reader := bufio.NewReader(os.Stdin)
	printedDefaultValue := ""
	if defaultValue != "" {
		printedDefaultValue = fmt.Sprintf(" [%s]", defaultValue)
	}
	fmt.Printf("%s%s: ", aurora.Bold(aurora.White(label)), aurora.Gray(12, printedDefaultValue))

	value, err := reader.ReadString('\n')
	value = strings.Trim(strings.TrimSuffix(value, "\n"), " ")
	if (err != nil || value == "") && defaultValue != "" {
		return defaultValue, true
	}
	return value, err == nil
}

func scanDate(label string, defaultValue time.Time) (time.Time, bool) {
	reader := bufio.NewReader(os.Stdin)
	printedDefaultValue := ""
	if !defaultValue.IsZero() {
		printedDefaultValue = fmt.Sprintf(" [%s]", defaultValue.Format("2006-01-02"))
	}
	fmt.Printf("%s%s: ", aurora.Bold(aurora.White(label)), aurora.Gray(12, printedDefaultValue))

	value, err := reader.ReadString('\n')
	value = strings.Trim(strings.TrimSuffix(value, "\n"), " ")
	if err != nil || value == "" {
		return defaultValue, !defaultValue.IsZero()
	}

	date, err := parseDate(value)
	return date, err == nil
}

func parseDate(value string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02T15:04:05", "2006-01-02"} {
		parse, err := time.Parse(layout, value)
		if err == nil {
			return parse, nil
		}
	}

	return time.Time{}, errors.Errorf("'%s' is not a valid date, or datetime, format.", value)
}
