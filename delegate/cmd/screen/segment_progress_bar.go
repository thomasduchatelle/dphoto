package screen

import (
	"fmt"
	"strings"
)

type ProgressBarSegment struct {
	done  uint
	total uint
}

func NewProgressBarSegment() (Segment, func(uint, uint)) {
	segment := &ProgressBarSegment{}
	return segment, func(done, total uint) {
		segment.done = done
		segment.total = total
	}
}

func (p *ProgressBarSegment) Content(options RenderingOptions) string {
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
