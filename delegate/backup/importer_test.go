package backup

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const layout = "2006-01-02"

func TestFoundAlbum_pushBoundaries(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name               string
		dates              []time.Time
		wantStart, wantEnd time.Time
	}{
		{"it should use the only date for start and end", []time.Time{mustParse("2021-06-02")}, mustParse("2021-06-02"), mustParse("2021-06-03")},
		{"it should push boundaries in both directions", []time.Time{mustParse("2021-06-02"), mustParse("2021-06-03"), mustParse("2021-06-01")}, mustParse("2021-06-01"), mustParse("2021-06-04")},
		{"it should push boundaries in start directions", []time.Time{mustParse("2021-06-02"), mustParse("2021-05-21"), mustParse("2021-06-01")}, mustParse("2021-05-21"), mustParse("2021-06-03")},
	}
	for _, tt := range tests {
		alb := newFoundAlbum("", tt.dates[0])
		for _, date := range tt.dates[1:] {
			alb.pushBoundaries(date)
		}

		a.Equal(tt.wantStart.Format(layout), alb.Start.Format(layout), tt.name)
		a.Equal(tt.wantEnd.Format(layout), alb.End.Format(layout), tt.name)
	}
}

func mustParse(value string) time.Time {
	date, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return date
}
