package screen

import (
	"fmt"
	"strings"
	"sync"
)

type ProgressBarSegment struct {
	done  int
	total int
	lock  sync.RWMutex
}

func NewProgressBarSegment() (Segment, func(int, int)) {
	segment := &ProgressBarSegment{lock: sync.RWMutex{}}
	return segment, func(done, total int) {
		segment.lock.Lock()
		defer segment.lock.Unlock()
		segment.done = done
		segment.total = total
	}
}

func (p *ProgressBarSegment) Content(options RenderingOptions) string {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if p.total == 0 {
		return strings.Repeat(" ", options.Width)
	}

	if p.done == p.total {
		return fmt.Sprintf("[%s]", strings.Repeat("=", options.Width-2))
	}

	barWidth := options.Width - 3
	progressWidth := barWidth * int(p.done) / int(p.total)
	if progressWidth < 0 {
		progressWidth = 0
	}

	restWidth := barWidth - progressWidth
	if restWidth < 0 {
		restWidth = 0
	}

	return fmt.Sprintf("[%s>%s]", strings.Repeat("=", progressWidth), strings.Repeat(" ", restWidth))
}
