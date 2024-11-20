package backup

import (
	"fmt"
	"slices"
	"strings"
)

func newEventCapture() *trackerEventsCapture {
	return &trackerEventsCapture{
		Captured: make(map[trackEvent]eventSummary),
	}
}

type trackerEventsCapture struct {
	Captured map[trackEvent]eventSummary
}

type eventSummary struct {
	SumCount int
	SumSize  int
	Albums   []string
}

func (e *eventSummary) String() string {
	return fmt.Sprintf("SumCount=%d, SumSize=%d, Albums=%s", e.SumCount, e.SumSize, strings.Join(e.Albums, ", "))
}

func (e *trackerEventsCapture) OnEvent(event progressEvent) {
	if e.Captured == nil {
		e.Captured = make(map[trackEvent]eventSummary)
	}

	capture, _ := e.Captured[event.Type]

	capture.SumCount += event.Count
	capture.SumSize += event.Size

	if event.Album != "" && !slices.Contains(capture.Albums, event.Album) {
		capture.Albums = append(capture.Albums, event.Album)
		slices.Sort(capture.Albums)
	}

	e.Captured[event.Type] = capture
}
