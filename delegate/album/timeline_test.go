package album

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

const layout = "2006-01-02T15"

type simplifiedSegment struct {
	folder string
	start  string
	end    string
}

func toSimplifiedSegment(s *segment) simplifiedSegment {
	name := "<empty>"
	if len(s.albums) > 0 {
		name = s.albums[0].FolderName
	}

	return simplifiedSegment{
		folder: name,
		start:  s.from.Format(layout),
		end:    s.to.Format(layout),
	}
}

func TestNewTimeline(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name    string
		args    []Album
		want    []simplifiedSegment
		wantErr bool
	}{
		{
			"it should support albums following up each-other",
			[]Album{
				newAlbum("2020-Q3", "2020-07-01T00", "2020-10-01T00"),
				newAlbum("2020-Q4", "2020-10-01T00", "2021-01-01T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("2020-Q3", "2020-07-01T00", "2020-10-01T00"),
				newSimplifiedSegment("2020-Q4", "2020-10-01T00", "2021-01-01T00"),
			},
			false,
		},
		{
			"it should support albums with a gap",
			[]Album{
				newAlbum("2020-Q3", "2020-07-01T00", "2020-10-01T00"),
				newAlbum("2020-Q4", "2020-11-01T00", "2021-02-01T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("2020-Q3", "2020-07-01T00", "2020-10-01T00"),
				newSimplifiedSegment("2020-Q4", "2020-11-01T00", "2021-02-01T00"),
			},
			false,
		},
		{
			"it should support albums overlapping - priority to the first",
			[]Album{
				newAlbum("A-01", "2020-12-01T00", "2020-12-11T00"),
				newAlbum("A-02", "2020-12-10T12", "2020-12-21T12"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("A-01", "2020-12-01T00", "2020-12-11T00"),
				newSimplifiedSegment("A-02", "2020-12-11T00", "2020-12-21T12"),
			},
			false,
		},
		{
			"it should support albums overlapping - priority to the shortest",
			[]Album{
				newAlbum("A-01", "2020-12-01T00", "2020-12-06T00"),
				newAlbum("A-02", "2020-12-02T00", "2020-12-05T00"),
				newAlbum("A-03", "2020-12-03T00", "2020-12-04T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("A-01", "2020-12-01T00", "2020-12-02T00"),
				newSimplifiedSegment("A-02", "2020-12-02T00", "2020-12-03T00"),
				newSimplifiedSegment("A-03", "2020-12-03T00", "2020-12-04T00"),
				newSimplifiedSegment("A-02", "2020-12-04T00", "2020-12-05T00"),
				newSimplifiedSegment("A-01", "2020-12-05T00", "2020-12-06T00"),
			},
			false,
		},
		{
			"it should support albums starting at the same time",
			[]Album{
				newAlbum("A-01", "2020-12-01T00", "2020-12-06T00"),
				newAlbum("A-02", "2020-12-01T00", "2020-12-05T00"),
				newAlbum("A-03", "2020-12-01T00", "2020-12-07T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("A-02", "2020-12-01T00", "2020-12-05T00"),
				newSimplifiedSegment("A-01", "2020-12-05T00", "2020-12-06T00"),
				newSimplifiedSegment("A-03", "2020-12-06T00", "2020-12-07T00"),
			},
			false,
		},
		{
			"it should support albums ending at the same time",
			[]Album{
				newAlbum("A-01", "2020-12-01T00", "2020-12-06T00"),
				newAlbum("A-02", "2020-12-02T00", "2020-12-06T00"),
				newAlbum("A-03", "2020-12-03T00", "2020-12-06T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("A-01", "2020-12-01T00", "2020-12-02T00"),
				newSimplifiedSegment("A-02", "2020-12-02T00", "2020-12-03T00"),
				newSimplifiedSegment("A-03", "2020-12-03T00", "2020-12-06T00"),
			},
			false,
		},
	}

	for _, tt := range tests {
		got, err := NewTimeline(tt.args)

		if !tt.wantErr && a.NoError(err, tt.name) && err == nil {
			segments := make([]simplifiedSegment, len(got.segments))
			for i, s := range got.segments {
				segments[i] = toSimplifiedSegment(&s)
			}

			a.Equal(tt.want, segments, tt.name)

		} else if tt.wantErr && a.Error(err, tt.name) {
			// all good
		}
	}
}

func TestFindAt_FindAllAt(t *testing.T) {
	a := assert.New(t)

	type Want struct {
		name     string
		allNames []string
	}
	tests := []struct {
		name string
		args time.Time
		want Want
	}{
		{
			"it should return no albums when out of any range",
			time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
			Want{"", nil},
		},
		{
			"it should return album on first second (inclusive)",
			time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
			Want{"2020-Q3", []string{"2020-Q3"}},
		},
		{
			"it should return albums up to its last second (range's end is exclusive)",
			time.Date(2020, 9, 30, 23, 59, 59, 999999, time.UTC),
			Want{"2020-Q3", []string{"2020-Q3"}},
		},
		{
			"it should return next album on its first second (inclusive)",
			time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
			Want{"2020-Q4", []string{"2020-Q4"}},
		},
		{
			"it should return all albums matching the date (album added to timeline are taking the priority each time)",
			time.Date(2020, 12, 25, 10, 0, 0, 0, time.UTC),
			Want{"Christmas Day", []string{"2020-Q4", "Christmas Day", "Christmas Holidays"}},
		},
		{
			"it should return all albums matching the date (2021-Q1 album added to the timeline is not taking the priority immediately)",
			time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC),
			Want{"New Year", []string{"2021-Q1", "Christmas Holidays", "New Year"}},
		},
	}

	timeline, err := NewTimeline(albumCollection())
	if a.NoError(err) {
		for _, tt := range tests {
			album := timeline.FindAt(tt.args)
			albumName := ""
			if album != nil {
				albumName = album.FolderName
			}

			a.Equal(tt.want.name, albumName, tt.name)

			var albums []string
			for _, a := range timeline.FindAllAt(tt.args) {
				albums = append(albums, a.FolderName)
			}
			sort.Strings(albums)

			a.Equal(tt.want.allNames, albums, tt.name)
		}
	}
}

func TestFindAt_FindBetween(t *testing.T) {
	a := assert.New(t)

	type Args struct {
		start string
		end   string
	}
	type Want struct {
		start    string
		end      string
		allNames []string
	}
	tests := []struct {
		name string
		args Args
		want []Want
	}{
		{
			"it should return no segment outside timeline boundaries",
			Args{"2019-01-01T00", "2019-01-01T00"},
			nil,
		},
		{
			"it should return a segment with dates updated to match the request",
			Args{"2020-01-01T00", "2020-09-01T00"},
			[]Want{{"2020-07-01T00", "2020-09-01T00", []string{"2020-Q3"}}},
		},
		{
			"it should return segments within the request",
			Args{"2020-12-31T00", "2021-01-02T00"},
			[]Want{
				{"2020-12-31T00", "2020-12-31T18", []string{"Christmas Holidays", "2020-Q4"}},
				{"2020-12-31T18", "2021-01-01T18", []string{"New Year", "Christmas Holidays", "2021-Q1", "2020-Q4"}},
				{"2021-01-01T18", "2021-01-02T00", []string{"Christmas Holidays", "2021-Q1"}},
			},
		},
		{
			"it should return segments within the request",
			Args{"2020-12-18T00", "2020-12-26T00"},
			[]Want{
				{"2020-12-18T00", "2020-12-24T00", []string{"Christmas First Week", "Christmas Holidays", "2020-Q4"}},
				{"2020-12-24T00", "2020-12-26T00", []string{"Christmas Day", "Christmas First Week", "Christmas Holidays", "2020-Q4"}},
			},
		},
	}

	timeline, err := NewTimeline(albumCollection())

	if a.NoError(err) {
		for _, tt := range tests {
			segments := timeline.FindBetween(mustParse(layout, tt.args.start), mustParse(layout, tt.args.end))
			var got []Want
			for _, seg := range segments {
				var names []string
				for _, a := range seg.Albums {
					names = append(names, a.FolderName)
				}

				got = append(got, Want{
					start:    seg.Start.Format(layout),
					end:      seg.End.Format(layout),
					allNames: names,
				})
			}

			a.Equal(tt.want, got, tt.name)
		}
	}
}

func TestTimeline_FindForAlbum(t *testing.T) {
	a := assert.New(t)

	type Want struct {
		start    string
		end      string
		allNames []string
	}

	timeline, err := NewTimeline(albumCollection())
	if a.NoError(err) {
		segments := timeline.FindForAlbum("Christmas Holidays")

		var got []Want
		for _, seg := range segments {
			var names []string
			for _, a := range seg.Albums {
				names = append(names, a.FolderName)
			}

			got = append(got, Want{
				start:    seg.Start.Format(layout),
				end:      seg.End.Format(layout),
				allNames: names,
			})
		}

		a.Equal([]Want{
			{"2020-12-26T00", "2020-12-31T18", []string{"Christmas Holidays", "2020-Q4"}},
			{"2021-01-01T18", "2021-01-04T00", []string{"Christmas Holidays", "2021-Q1"}},
		}, got)
	}
}

func newSimplifiedSegment(folder, start, end string) simplifiedSegment {
	return simplifiedSegment{
		folder: folder,
		start:  start,
		end:    end,
	}
}

func newAlbum(folder, start, end string) Album {
	startTime, err := time.Parse(layout, start)
	if err != nil {
		panic(err)
	}
	endTime, err := time.Parse(layout, end)
	if err != nil {
		panic(err)
	}
	return Album{
		Name:       folder,
		FolderName: folder,
		Start:      startTime,
		End:        endTime,
	}
}
