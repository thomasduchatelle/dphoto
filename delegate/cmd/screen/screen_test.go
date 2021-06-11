package screen

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleScreen_updateMaxLength(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name string
		args string
		want int
	}{
		{"it should find 0 for an empty string", "", 0},
		{"it should find 0 for an empty string in multiline", "\n", 0},
		{"it should find len of the single line string", "this is not long", 16},
		{"it should find max len of the last of a multi line string", "foo\nbar\nthis is not long", 16},
		{"it should find max len of the first of a multi line string", "foo\nbar\nbe", 3},
		{"it should find max len of multiline string with empty lines", "foo\n\nbe", 3},
	}

	for _, tt := range tests {
		s := &SimpleScreen{}
		s.updateMaxWidth(tt.args)
		a.Equal(tt.want, s.maxWidth, tt.name)
	}
}
