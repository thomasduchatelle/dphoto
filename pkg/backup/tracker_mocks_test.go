package backup_test

import (
	"fmt"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"strings"
)

type capturedEvents struct {
	Captured map[backup.ProgressEventType]eventSummary
}

type eventSummary struct {
	Number   int
	SumCount int
	SumSize  int
	Albums   []string
}

func (e *eventSummary) String() string {
	return fmt.Sprintf("Number=%d , SumCount=%d, SumSize=%d, Albums=%s", e.Number, e.SumCount, e.SumSize, strings.Join(e.Albums, ", "))
}

func newEventCapture() *capturedEvents {
	return &capturedEvents{
		Captured: make(map[backup.ProgressEventType]eventSummary),
	}
}

func (e *capturedEvents) OnEvent(event backup.ProgressEvent) {
	capture, found := e.Captured[event.Type]
	if !found {
		capture = eventSummary{}
	}

	capture.Number++
	capture.SumCount += int(event.Count)
	capture.SumSize += int(event.Size)

	if event.Album != "" {
		capture.Albums = append(capture.Albums, event.Album)
	}

	e.Captured[event.Type] = capture
}
