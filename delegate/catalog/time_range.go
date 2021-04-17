package catalog

import (
	"fmt"
	"time"
)

type timeRange struct {
	Start time.Time
	End   time.Time
}

func (t timeRange) Plus(other timeRange) (ranges []timeRange) {
	first := t
	second := other
	if second.Start.Before(first.Start) {
		first, second = second, first
	}

	if second.Start.Before(first.End) {
		return []timeRange{
			{first.Start, maxTime(first.End, second.End)},
		}
	}

	return []timeRange{t, other}
}

func (t timeRange) Minus(other timeRange) (ranges []timeRange) {
	if other.Start.After(t.Start) {
		ranges = append(ranges, timeRange{Start: t.Start, End: minTime(t.End, other.Start)})
	}

	if other.End.Before(t.End) {
		ranges = append(ranges, timeRange{Start: maxTime(t.Start, other.End), End: t.End})
	}

	return ranges
}

func (t timeRange) Equals(other timeRange) bool {
	return t.Start.Equal(other.Start) && t.End.Equal(other.End)
}

func (t timeRange) String() string {
	return fmt.Sprintf("%s -> %s", t.Start.Format(time.RFC3339), t.End.Format(time.RFC3339))
}

func minTime(a, b time.Time) time.Time {
	if a.Unix() < b.Unix() {
		return a
	}

	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.Unix() > b.Unix() {
		return a
	}

	return b
}
