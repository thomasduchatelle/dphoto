package backup

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const layout = "2006-01-02"

func TestFoundAlbum_pushBoundaries(t *testing.T) {
	a := assert.New(t)

	type found struct {
		date time.Time
		size int
	}
	tests := []struct {
		name               string
		found              []found
		wantStart, wantEnd time.Time
		wantSize           int
	}{
		{"it should use the only date for start and end", []found{{mustParse("2021-06-02"), 42}}, mustParse("2021-06-02"), mustParse("2021-06-03"), 1},
		{"it should push boundaries in both directions", []found{{mustParse("2021-06-02"), 1}, {mustParse("2021-06-03"), 2}, {mustParse("2021-06-01"), 4}}, mustParse("2021-06-01"), mustParse("2021-06-04"), 3},
		{"it should push boundaries in start directions", []found{{mustParse("2021-06-02"), 1}, {mustParse("2021-05-21"), 2}, {mustParse("2021-06-01"), 4}}, mustParse("2021-05-21"), mustParse("2021-06-03"), 3},
	}
	for _, tt := range tests {
		alb := &ScannedFolder{Distribution: make(map[string]MediaCounter)}
		for _, f := range tt.found {
			alb.PushBoundaries(f.date, f.size)
		}

		a.Equal(tt.wantStart.Format(layout), alb.Start.Format(layout), tt.name)
		a.Equal(tt.wantEnd.Format(layout), alb.End.Format(layout), tt.name)

		count := 0
		for _, d := range alb.Distribution {
			count += d.Count
		}
		a.Equal(tt.wantSize, count, tt.name)
	}
}

func mustParse(value string) time.Time {
	date, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return date
}
