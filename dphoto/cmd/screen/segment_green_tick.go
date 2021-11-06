package screen

import "github.com/logrusorgru/aurora/v3"

func NewGreenTickSegment() Segment {
	return NewConstantSegment(aurora.Green("\xE2\x9C\x93 ").String())
}
