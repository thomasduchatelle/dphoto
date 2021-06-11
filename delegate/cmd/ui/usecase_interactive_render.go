package ui

import (
	"duchatelle.io/dphoto/dphoto/cmd/screen"
	"strings"
)

type interactiveRender struct {
	recordsRenderer *recordsRenderer
	screen          *screen.SimpleScreen
}

func newInteractiveRender() *interactiveRender {
	return &interactiveRender{
		recordsRenderer: new(recordsRenderer),
		screen:          screen.NewSimpleScreen(screen.RenderingOptions{Full: true}),
	}
}

func (i *interactiveRender) Render(state *interactiveViewState) error {
	table, err := i.recordsRenderer.Render(&state.recordsState)
	if err != nil {
		return nil
	}

	i.screen.Clear()
	i.screen.Print(screen.PagePrint{
		Content: []screen.Segment{screen.NewConstantSegment(table)},
		Footer:  []screen.Segment{screen.NewConstantSegment(strings.Join(state.Actions, " ; "))},
	})

	return nil
}
