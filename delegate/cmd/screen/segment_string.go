package screen

type StringSegment struct {
	content string
}

// NewConstantSegment creates a constant segment
func NewConstantSegment(content string) Segment {
	return &StringSegment{
		content: content,
	}
}

// NewUpdatableSegment returns the segment and the function to update it
func NewUpdatableSegment(content string) (Segment, func(string)) {
	segment := &StringSegment{
		content: content,
	}

	return segment, func(newValue string) {
		segment.content = newValue
	}
}

func (l *StringSegment) Content(options RenderingOptions) string {
	return utf8Fill(l.content, " ", options.Width)
}