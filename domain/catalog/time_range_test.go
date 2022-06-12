package catalog

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeRange_Plus(t *testing.T) {
	a := assert.New(t)

	type fields struct {
		Start string
		End   string
	}
	tests := []struct {
		name       string
		timeRange  fields
		other      fields
		wantRanges []fields
	}{
		{"it should return both ranges if other is before", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-04", "2020-12-05"}, []fields{{"2020-12-05", "2020-12-10"}, {"2020-12-04", "2020-12-05"}}},
		{"it should return both ranges if other is after", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-10", "2020-12-15"}, []fields{{"2020-12-05", "2020-12-10"}, {"2020-12-10", "2020-12-15"}}},
		{"it should return first range if second included (starts)", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-05", "2020-12-08"}, []fields{{"2020-12-05", "2020-12-10"}}},
		{"it should return first range if second included (ends)", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-08", "2020-12-10"}, []fields{{"2020-12-05", "2020-12-10"}}},
		{"it should return second range if first is included", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-01", "2020-12-31"}, []fields{{"2020-12-01", "2020-12-31"}}},
		{"it should return first range when other is equals", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-05", "2020-12-10"}, []fields{{"2020-12-05", "2020-12-10"}}},
	}

	const layout = "2006-01-02"
	for _, tt := range tests {
		t := TimeRange{
			Start: MustParse(layout, tt.timeRange.Start),
			End:   MustParse(layout, tt.timeRange.End),
		}

		gotRanges := t.Plus(TimeRange{Start: MustParse(layout, tt.other.Start), End: MustParse(layout, tt.other.End)})

		var expected []TimeRange
		for _, r := range tt.wantRanges {
			expected = append(expected, TimeRange{Start: MustParse(layout, r.Start), End: MustParse(layout, r.End)})
		}

		a.Equal(expected, gotRanges, tt.name)
	}
}
func TestTimeRange_Equal(t *testing.T) {
	a := assert.New(t)

	type fields struct {
		Start string
		End   string
	}
	tests := []struct {
		name      string
		timeRange fields
		other     fields
		want      bool
	}{
		{"it should return false if other is before", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-04", "2020-12-05"}, false},
		{"it should return false if other is after", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-10", "2020-12-15"}, false},
		{"it should return false if second included (starts)", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-05", "2020-12-08"}, false},
		{"it should return false if second included (ends)", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-08", "2020-12-10"}, false},
		{"it should return false if first is included", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-01", "2020-12-31"}, false},
		{"it should return true when they are equals", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-05", "2020-12-10"}, true},
	}

	const layout = "2006-01-02"
	for _, tt := range tests {
		t := TimeRange{
			Start: MustParse(layout, tt.timeRange.Start),
			End:   MustParse(layout, tt.timeRange.End),
		}

		got := t.Equals(TimeRange{Start: MustParse(layout, tt.other.Start), End: MustParse(layout, tt.other.End)})
		a.Equal(tt.want, got, tt.name)
	}
}

func TestTimeRange_Minus(t *testing.T) {
	a := assert.New(t)

	type fields struct {
		Start string
		End   string
	}
	tests := []struct {
		name       string
		timeRange  fields
		other      fields
		wantRanges []fields
	}{
		{"it should return full range if other is before", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-04", "2020-12-05"}, []fields{{"2020-12-05", "2020-12-10"}}},
		{"it should return full range if other is after", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-10", "2020-12-15"}, []fields{{"2020-12-05", "2020-12-10"}}},
		{"it should return range beginning if other is at the beginning", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-05", "2020-12-08"}, []fields{{"2020-12-08", "2020-12-10"}}},
		{"it should return range ending if other is at the end", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-08", "2020-12-10"}, []fields{{"2020-12-05", "2020-12-08"}}},
		{"it should return empty when fully covered by other", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-01", "2020-12-31"}, []fields{}},
		{"it should return empty when other is equals", fields{"2020-12-05", "2020-12-10"}, fields{"2020-12-05", "2020-12-10"}, []fields{}},
	}

	const layout = "2006-01-02"
	for _, tt := range tests {
		t := TimeRange{
			Start: MustParse(layout, tt.timeRange.Start),
			End:   MustParse(layout, tt.timeRange.End),
		}

		gotRanges := t.Minus(TimeRange{Start: MustParse(layout, tt.other.Start), End: MustParse(layout, tt.other.End)})

		var expected []TimeRange
		for _, r := range tt.wantRanges {
			expected = append(expected, TimeRange{Start: MustParse(layout, r.Start), End: MustParse(layout, r.End)})
		}

		a.Equal(expected, gotRanges, tt.name)
	}
}

func MustParse(layout string, date string) time.Time {
	parse, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}
	return parse
}
