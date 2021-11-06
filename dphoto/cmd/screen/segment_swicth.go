package screen

import "sync"

type SwitchSegment struct {
	current   int
	delegates []Segment
	lock      sync.RWMutex
}

// NewSwitchSegment creates a segment that can switch from one implementation to the other ; returns func to swap that will panic if index out of bound.
func NewSwitchSegment(delegates ...Segment) (Segment, func(int)) {
	segment := &SwitchSegment{
		current:   0,
		delegates: delegates,
		lock:      sync.RWMutex{},
	}
	return segment, func(newCurrent int) {
		segment.lock.Lock()
		defer segment.lock.Unlock()
		segment.current = newCurrent
	}
}

func (s *SwitchSegment) Content(options RenderingOptions) string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.delegates[s.current].Content(options)
}
