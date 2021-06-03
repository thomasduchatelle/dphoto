package screen

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type SimpleScreen struct {
	lines                []Segment
	numberOfPrintedLines int
	renderingOptions     RenderingOptions
	maxWidth             int
}

type AutoRefreshScreen struct {
	SimpleScreen
	cancel context.CancelFunc
	done   chan struct{}
}

type Segment interface {
	Content(options RenderingOptions) string
}

func NewSimpleScreen(options RenderingOptions, bars ...Segment) *SimpleScreen {
	return &SimpleScreen{
		lines:            bars,
		renderingOptions: options,
	}
}

func NewAutoRefreshScreen(options RenderingOptions, bars ...Segment) *AutoRefreshScreen {
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

func (s *SimpleScreen) Clear() {
	if s.numberOfPrintedLines > 0 {
		fmt.Printf("\033[%dA", s.numberOfPrintedLines)
		fmt.Printf(strings.Repeat(strings.Repeat(" ", s.maxWidth)+"\n", s.numberOfPrintedLines))
		fmt.Printf("\033[%dA", s.numberOfPrintedLines)

		s.numberOfPrintedLines = 0
	}
}

func (s *SimpleScreen) Refresh() {
	if s.numberOfPrintedLines > 0 {
		fmt.Printf("\033[%dA", s.numberOfPrintedLines)
	}
	s.numberOfPrintedLines = 0
	for _, line := range s.lines {
		content := strings.Trim(line.Content(s.renderingOptions), "\n")
		s.numberOfPrintedLines += strings.Count(content, "\n") + 1
		s.updateMaxLength(content)
		fmt.Println(content)
	}
}

func (s *SimpleScreen) updateMaxLength(content string) {
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
