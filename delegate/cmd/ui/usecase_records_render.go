package ui

import (
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/logrusorgru/aurora/v3"
)

type recordsRenderer struct{}

func (r *recordsRenderer) Render(state *recordsState) (string, error) {
	const layout = "02/01/2006 (Mon)"

	table := simpletable.New()
	table.Header = &simpletable.Header{Cells: []*simpletable.Cell{
		{Text: "Name"},
		{Text: "Folder Name"},
		{Text: "Start"},
		{Text: "End"},
		{Text: "Files"},
	}}
	table.Body = &simpletable.Body{Cells: make([][]*simpletable.Cell, len(state.Records))}

	for idx, album := range state.Records {
		isSelected := state.Selected == idx

		countContent := "-"
		if album.Count > 0 {
			countContent = fmt.Sprint(album.Count)
		}

		table.Body.Cells[idx] = []*simpletable.Cell{
			{Text: r.applyStyleOnName(isSelected, album.Suggestion, album.Name)},
			{Text: r.applyStyle(isSelected, album.Suggestion, album.FolderName)},
			{Text: r.applyStyle(isSelected, album.Suggestion, album.Start.Format(layout))},
			{Text: r.applyStyle(isSelected, album.Suggestion, album.End.Format(layout))},
			{Text: r.applyStyle(isSelected, album.Suggestion, countContent), Align: simpletable.AlignRight},
		}
	}

	return table.String(), nil
}

func (r *recordsRenderer) applyStyle(selected, suggestion bool, args interface{}) string {
	if selected {
		return aurora.Black(aurora.Bold(aurora.BgWhite(args))).String()
	}

	return fmt.Sprint(args)
}

func (r *recordsRenderer) applyStyleOnName(selected, suggestion bool, args interface{}) string {
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
