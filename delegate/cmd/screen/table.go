package screen

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type TableGenerator struct {
	separator string
	widths    []int
}

type row struct {
	separator string
	widths    []int
	delegates []Segment
}

func NewTable(separator string, widths ...int) *TableGenerator {
	return &TableGenerator{
		separator: separator,
		widths:    widths,
	}
}

func (t *TableGenerator) NewRow(delegates ...Segment) Segment {
	if len(delegates) != len(t.widths) {
		panic(fmt.Sprintf("Number of delegates (%d) must be the same as number of widths: %d", len(delegates), len(t.widths)))
	}

	return &row{
		separator: t.separator,
		widths:    t.widths,
		delegates: delegates,
	}
}

func (r *row) Content(options RenderingOptions) string {
	content := make([]string, len(r.delegates))
	for i, delegate := range r.delegates {
		options.Width = r.widths[i]
		content[i] = delegate.Content(options)
	}

	return strings.Join(content, r.separator)
}

func utf8Fill(value, filler string, count int) string {
	actualLength := utf8.RuneCountInString(value)
	if actualLength >= count {
		return value
	}

	return value + strings.Repeat(filler, count-actualLength)
}
