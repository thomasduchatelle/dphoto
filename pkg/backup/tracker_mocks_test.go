package backup

import (
	"fmt"
	"slices"
	"strings"
)

type capturedEvents struct {
	Captured map[ProgressEventType]eventSummary
}

type eventSummary struct {
	SumCount int
	SumSize  int
	Albums   []string
}

func (e *eventSummary) String() string {
	return fmt.Sprintf("SumCount=%d, SumSize=%d, Albums=%s", e.SumCount, e.SumSize, strings.Join(e.Albums, ", "))
}

func newEventCapture() *capturedEvents {
	return &capturedEvents{
		Captured: make(map[ProgressEventType]eventSummary),
	}
}

func (e *capturedEvents) OnEvent(event ProgressEvent) {
	capture, _ := e.Captured[event.Type]

	capture.SumCount += event.Count
	capture.SumSize += event.Size

	if event.Album != "" && !slices.Contains(capture.Albums, event.Album) {
		capture.Albums = append(capture.Albums, event.Album)
		slices.Sort(capture.Albums)
	}

	e.Captured[event.Type] = capture
}
