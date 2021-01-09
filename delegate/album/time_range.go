package album

import "time"

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func (t TimeRange) Plus(other TimeRange) (ranges []TimeRange) {
	first := t
	second := other
	if second.Start.Before(first.Start) {
		first, second = second, first
	}

	if second.Start.Before(first.End) {
		return []TimeRange{
			{first.Start, maxTime(first.End, second.End)},
		}
	}

	return []TimeRange{t, other}
}

func (t TimeRange) Minus(other TimeRange) (ranges []TimeRange) {
	if other.Start.After(t.Start) {
		ranges = append(ranges, TimeRange{Start: t.Start, End: minTime(t.End, other.Start)})
	}

	if other.End.Before(t.End) {
		ranges = append(ranges, TimeRange{Start: maxTime(t.Start, other.End), End: t.End})
	}

	return ranges
}

func (t TimeRange) Equals(other TimeRange) bool {
	return t.Start.Equal(other.Start) && t.End.Equal(other.End)
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
