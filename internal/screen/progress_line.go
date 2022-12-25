package screen

type ProgressLine struct {
	SwapSpinner    func(int)
	SetLabel       func(string)
	SetBar         func(int, int)
	SetExplanation func(string)
}

func NewProgressLine(table *TableGenerator, initialLabel string) (*ProgressLine, Segment) {
	spinner, swapSpinner := NewSwitchSegment(NewSpinnerSegment(), NewGreenTickSegment())
	label, setLabel := NewUpdatableSegment(initialLabel)
	bar, setBar := NewProgressBarSegment()
	explanation, setExplanation := NewUpdatableSegment("")

	return &ProgressLine{
		SwapSpinner:    swapSpinner,
		SetLabel:       setLabel,
		SetBar:         setBar,
		SetExplanation: setExplanation,
	}, table.NewRow(spinner, label, bar, explanation)
}
