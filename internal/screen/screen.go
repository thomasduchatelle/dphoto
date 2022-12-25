// Package screen provides building blocks to represent progress bars.
package screen

import (
	"context"
	"fmt"
	"github.com/buger/goterm"
	"strings"
	"time"
)

type SimpleScreen struct {
	lines                []Segment
	numberOfPrintedLines int
	contentHeight        int
	renderingOptions     RenderingOptions
	maxWidth             int
}

type AutoRefreshScreen struct {
	SimpleScreen
	cancel context.CancelFunc
	done   chan struct{}
}

type PagePrint struct {
	Content []Segment
	Footer  []Segment
}

type Segment interface {
	Content(options RenderingOptions) string
}

func NewSimpleScreen(options RenderingOptions, bars ...Segment) *SimpleScreen {
	if options.Full {
		goterm.Clear()
	}
	return &SimpleScreen{
		lines:            bars,
		renderingOptions: options,
	}
}

func NewAutoRefreshScreen(options RenderingOptions, bars ...Segment) *AutoRefreshScreen {
	if options.Full {
		goterm.Clear()
	}
	ctx, cancel := context.WithCancel(context.Background())
	progress := &AutoRefreshScreen{
		SimpleScreen: SimpleScreen{
			lines:            bars,
			renderingOptions: options,
		},
		cancel: cancel,
		done:   make(chan struct{}),
	}

	progress.start(ctx)
	return progress
}

// ForceClear goes back to the beginning after writing white spaces everywhere
func (s *SimpleScreen) ForceClear() {
	if s.numberOfPrintedLines > 0 {
		fmt.Printf("\033[%dA", s.numberOfPrintedLines)
		fmt.Printf(strings.Repeat(strings.Repeat(" ", s.maxWidth)+"\n", s.numberOfPrintedLines))
		fmt.Printf("\033[%dA", s.numberOfPrintedLines)

		s.numberOfPrintedLines = 0
	}
}

// Clear goes back to the beginning
func (s *SimpleScreen) Clear() {
	if s.renderingOptions.Full {
		goterm.MoveCursor(1, 1)
		for l := 1; l < goterm.Height(); l++ {
			_, _ = goterm.Println(strings.Repeat(" ", goterm.Width()))
		}
		goterm.Flush()
	} else if s.numberOfPrintedLines > 0 {
		fmt.Printf("\033[%dA", s.numberOfPrintedLines)
	}
	s.numberOfPrintedLines = 0
}

func (s *SimpleScreen) Refresh() {
	s.Clear()
	s.Print(PagePrint{Content: s.lines})
}

// Print writes given segments in the output ; it doesn't keep the segments and doesn't refresh the page.
func (s *SimpleScreen) Print(page PagePrint) {
	footer := make([]string, len(page.Footer))
	footerHeight := 0
	for i, line := range page.Footer {
		footer[i] = strings.Trim(line.Content(s.renderingOptions), "\n")
		footerHeight += strings.Count(footer[i], "\n") + 1
	}

	content := make([]string, len(page.Content), len(page.Content))
	contentHeight := 0
	for i, line := range page.Content {
		content[i] = strings.Trim(line.Content(s.renderingOptions), "\n")
		s.updateMaxWidth(content[i])
		contentHeight += strings.Count(content[i], "\n") + 1
	}

	s.contentHeight = contentHeight
	if s.renderingOptions.Full {
		goterm.MoveCursor(1, 1)
		_, _ = goterm.Println(strings.Join(content, "\n"))
		goterm.MoveCursor(1, goterm.Height()-footerHeight)
		_, _ = goterm.Println(strings.Join(footer, "\n"))
		goterm.Flush()
		s.numberOfPrintedLines = goterm.Height()

	} else {
		s.numberOfPrintedLines = contentHeight + footerHeight + 1
		fmt.Println(strings.Join(content, "\n"))
		fmt.Println(strings.Join(footer, "\n"))
	}
}

func (s *SimpleScreen) updateMaxWidth(content string) {
	substr := content
	for len(substr) > 0 {
		idx := strings.Index(substr, "\n")
		if idx > s.maxWidth {
			s.maxWidth = idx
		} else if idx == -1 && len(substr) > s.maxWidth {
			s.maxWidth = len(substr)
		}

		if idx >= 0 {
			substr = substr[idx+1:]
		} else {
			substr = ""
		}
	}

	if goterm.Width() > 0 && s.maxWidth > goterm.Width() {
		s.maxWidth = goterm.Width()
	}
}

// TermSize returns the size of the terminal, in chars
func (s *SimpleScreen) TermSize() (int, int) {
	return goterm.Width(), goterm.Height()
}

// ContentHeight returns number of lines used in content (excluding footer)
func (s *SimpleScreen) ContentHeight() int {
	return s.contentHeight
}

func (s *AutoRefreshScreen) start(ctx context.Context) {
	go func() {
		defer close(s.done)
		tick := time.Tick(100 * time.Millisecond)

		for {
			select {
			case <-ctx.Done():
				s.Refresh()
				return
			case <-tick:
				s.Refresh()
			}
		}
	}()
}

func (s *AutoRefreshScreen) Stop() {
	s.cancel()
	<-s.done
}
