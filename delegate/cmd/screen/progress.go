package screen

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Screen struct {
	cancel               context.CancelFunc
	done                 chan struct{}
	lines                []Segment
	numberOfPrintedLines int
	renderingOptions     RenderingOptions
}

type Segment interface {
	Content(options RenderingOptions) string
}

func NewScreen(options RenderingOptions, bars ...Segment) *Screen {
	ctx, cancel := context.WithCancel(context.Background())
	progress := &Screen{
		cancel:           cancel,
		lines:            bars,
		done:             make(chan struct{}),
		renderingOptions: options,
	}

	progress.start(ctx)
	return progress
}

func (s *Screen) start(ctx context.Context) {
	go func() {
		defer close(s.done)
		tick := time.Tick(100 * time.Millisecond)

		for {
			select {
			case <-ctx.Done():
				s.refresh()
				return
			case <-tick:
				s.refresh()
			}
		}
	}()
}

func (s *Screen) refresh() {
	if s.numberOfPrintedLines > 0 {
		fmt.Printf("\033[%dA", s.numberOfPrintedLines)
	}

	s.numberOfPrintedLines = 0
	for _, line := range s.lines {
		content := strings.Trim(line.Content(s.renderingOptions), "\n")
		s.numberOfPrintedLines += strings.Count(content, "\n") + 1
		fmt.Println(content)
	}
}

func (s *Screen) Stop() {
	s.cancel()
	<-s.done
}
