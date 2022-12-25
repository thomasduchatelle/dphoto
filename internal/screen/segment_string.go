package screen

import "sync"

type StringSegment struct {
	content string
	lock    sync.RWMutex
}

// NewConstantSegment creates a constant segment
func NewConstantSegment(content string) Segment {
	return &StringSegment{
		content: content,
		lock:    sync.RWMutex{},
	}
}

// NewUpdatableSegment returns the segment and the function to update it
func NewUpdatableSegment(content string) (Segment, func(string)) {
	segment := &StringSegment{
		content: content,
	}

	return segment, func(newValue string) {
		segment.lock.Lock()
		defer segment.lock.Unlock()
		segment.content = newValue
	}
}

func (l *StringSegment) Content(options RenderingOptions) string {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return utf8Fill(l.content, " ", options.Width)
}
