package screen

var SpinnerDefaultStyle = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type LineSpinnerSegment struct {
	count int
}

func NewSpinnerSegment() *LineSpinnerSegment {
	return new(LineSpinnerSegment)
}

func (l *LineSpinnerSegment) Content(options RenderingOptions) string {
	l.count++
	return utf8Fill(SpinnerDefaultStyle[l.count%len(SpinnerDefaultStyle)], " ", options.Width)
}
