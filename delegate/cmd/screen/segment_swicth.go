package screen

type SwitchSegment struct {
	current   int
	delegates []Segment
}

// NewSwitchSegment creates a segment that can switch from one implementation to the other ; returns func to swap that will panic if index out of bound.
func NewSwitchSegment(delegates ...Segment) (Segment, func(int)) {
	segment := &SwitchSegment{
		current:   0,
		delegates: delegates,
	}
	return segment, func(newCurrent int) {
		segment.current = newCurrent
	}
}

func (s *SwitchSegment) Content(options RenderingOptions) string {
	return s.delegates[s.current].Content(options)
}
