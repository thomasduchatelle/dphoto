package catalog

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
		args    []*Album
		want    []simplifiedSegment
		wantErr bool
	}{
		{
			"it should support albums following up each-other",
			[]*Album{
				newAlbum("/2020-Q3", "2020-07-01T00", "2020-10-01T00"),
				newAlbum("/2020-Q4", "2020-10-01T00", "2021-01-01T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/2020-Q3", "2020-07-01T00", "2020-10-01T00"),
				newSimplifiedSegment("/2020-Q4", "2020-10-01T00", "2021-01-01T00"),
			},
			false,
		},
		{
			"it should support albums with a gap",
			[]*Album{
				newAlbum("/2020-Q3", "2020-07-01T00", "2020-10-01T00"),
				newAlbum("/2020-Q4", "2020-11-01T00", "2021-02-01T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/2020-Q3", "2020-07-01T00", "2020-10-01T00"),
				newSimplifiedSegment("/2020-Q4", "2020-11-01T00", "2021-02-01T00"),
			},
			false,
		},
		{
			"it should support albums overlapping - priority to the first",
			[]*Album{
				newAlbum("/A-01", "2020-12-01T00", "2020-12-11T00"),
				newAlbum("/A-02", "2020-12-10T12", "2020-12-21T12"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/A-01", "2020-12-01T00", "2020-12-11T00"),
				newSimplifiedSegment("/A-02", "2020-12-11T00", "2020-12-21T12"),
			},
			false,
		},
		{
			"it should support albums overlapping - priority to the shortest",
			[]*Album{
				newAlbum("/A-01", "2020-12-01T00", "2020-12-06T00"),
				newAlbum("/A-02", "2020-12-02T00", "2020-12-05T00"),
				newAlbum("/A-03", "2020-12-03T00", "2020-12-04T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/A-01", "2020-12-01T00", "2020-12-02T00"),
				newSimplifiedSegment("/A-02", "2020-12-02T00", "2020-12-03T00"),
				newSimplifiedSegment("/A-03", "2020-12-03T00", "2020-12-04T00"),
				newSimplifiedSegment("/A-02", "2020-12-04T00", "2020-12-05T00"),
				newSimplifiedSegment("/A-01", "2020-12-05T00", "2020-12-06T00"),
			},
			false,
		},
		{
			"it should support albums starting at the same time",
			[]*Album{
				newAlbum("/A-01", "2020-12-01T00", "2020-12-06T00"),
				newAlbum("/A-02", "2020-12-01T00", "2020-12-05T00"),
				newAlbum("/A-03", "2020-12-01T00", "2020-12-07T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/A-02", "2020-12-01T00", "2020-12-05T00"),
				newSimplifiedSegment("/A-01", "2020-12-05T00", "2020-12-06T00"),
				newSimplifiedSegment("/A-03", "2020-12-06T00", "2020-12-07T00"),
			},
			false,
		},
		{
			"it should support albums ending at the same time",
			[]*Album{
				newAlbum("/A-01", "2020-12-01T00", "2020-12-06T00"),
				newAlbum("/A-02", "2020-12-02T00", "2020-12-06T00"),
				newAlbum("/A-03", "2020-12-03T00", "2020-12-06T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/A-01", "2020-12-01T00", "2020-12-02T00"),
				newSimplifiedSegment("/A-02", "2020-12-02T00", "2020-12-03T00"),
				newSimplifiedSegment("/A-03", "2020-12-03T00", "2020-12-06T00"),
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
			Want{"/2020-Q3", []string{"/2020-Q3"}},
		},
		{
			"it should return albums up to its last second (range's end is exclusive)",
			time.Date(2020, 9, 30, 23, 59, 59, 999999, time.UTC),
			Want{"/2020-Q3", []string{"/2020-Q3"}},
		},
		{
			"it should return next album on its first second (inclusive)",
			time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
			Want{"/2020-Q4", []string{"/2020-Q4"}},
		},
		{
			"it should return all albums matching the date (album added to timeline are taking the priority each time)",
			time.Date(2020, 12, 25, 10, 0, 0, 0, time.UTC),
			Want{"/Christmas_Day", []string{"/2020-Q4", "/Christmas_Day", "/Christmas_Holidays"}},
		},
		{
			"it should return all albums matching the date (2021-Q1 album added to the timeline is not taking the priority immediately)",
			time.Date(2021, 1, 1, 10, 0, 0, 0, time.UTC),
			Want{"/New_Year", []string{"/2021-Q1", "/Christmas_Holidays", "/New_Year"}},
		},
	}

	timeline, err := NewTimeline(AlbumCollection())
	if a.NoError(err) {
		for _, tt := range tests {
			album, ok := timeline.FindAt(tt.args)
			if tt.want.name == "" {
				a.False(ok, tt.name)
			} else {
				a.Equal(tt.want.name, album.FolderName, tt.name)
			}

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
		name       string
		args       Args
		want       []Want
		wantMissed []Want
	}{
		{
			"it should return no segment outside timeline boundaries",
			Args{"2019-01-01T00", "2019-01-01T00"},
			nil,
			[]Want{{start: "2019-01-01T00", end: "2019-01-01T00"}},
		},
		{
			"it should return a segment with dates updated to match the request",
			Args{"2020-01-01T00", "2020-09-01T00"},
			[]Want{{"2020-07-01T00", "2020-09-01T00", []string{"/2020-Q3"}}},
			[]Want{{start: "2020-01-01T00", end: "2020-07-01T00"}},
		},
		{
			"it should return segments within the request",
			Args{"2020-12-31T00", "2021-01-02T00"},
			[]Want{
				{"2020-12-31T00", "2020-12-31T18", []string{"/Christmas_Holidays", "/2020-Q4"}},
				{"2020-12-31T18", "2021-01-01T18", []string{"/New_Year", "/Christmas_Holidays", "/2021-Q1", "/2020-Q4"}},
				{"2021-01-01T18", "2021-01-02T00", []string{"/Christmas_Holidays", "/2021-Q1"}},
			},
			nil,
		},
		{
			"it should return segments within the request",
			Args{"2020-12-18T00", "2020-12-26T00"},
			[]Want{
				{"2020-12-18T00", "2020-12-24T00", []string{"/Christmas_First_Week", "/Christmas_Holidays", "/2020-Q4"}},
				{"2020-12-24T00", "2020-12-26T00", []string{"/Christmas_Day", "/Christmas_First_Week", "/Christmas_Holidays", "/2020-Q4"}},
			},
			nil,
		},
		{
			"it should notice the gap between mars and may, and the missing dates at the end of the request",
			Args{"2021-03-23T00", "2021-06-26T00"},
			[]Want{
				{"2021-03-23T00", "2021-04-01T00", []string{"/2021-Q1"}},
				{"2021-05-01T00", "2021-06-01T00", []string{"/2021-May"}},
			},
			[]Want{
				{"2021-04-01T00", "2021-05-01T00", nil},
				{"2021-06-01T00", "2021-06-26T00", nil},
			},
		},
	}

	timeline, err := NewTimeline(AlbumCollection())

	if a.NoError(err) {
		for _, tt := range tests {
			segments, missed := timeline.FindBetween(MustParse(layout, tt.args.start), MustParse(layout, tt.args.end))
			var got []Want
			var gotMissed []Want
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
			for _, seg := range missed {
				gotMissed = append(gotMissed, Want{
					start: seg.Start.Format(layout),
					end:   seg.End.Format(layout),
				})
			}

			a.Equal(tt.want, got, tt.name)
			a.Equal(tt.wantMissed, gotMissed, tt.name)
		}
	}
}

func TestTimeline_FindForAlbum(t *testing.T) {
	a := assert.New(t)
	const owner = "ironman"

	type Want struct {
		start    string
		end      string
		allNames []string
	}

	timeline, err := NewTimeline(AlbumCollection())
	if a.NoError(err) {
		segments := timeline.FindForAlbum(owner, "/Christmas_Holidays")

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
			{"2020-12-26T00", "2020-12-31T18", []string{"/Christmas_Holidays", "/2020-Q4"}},
			{"2021-01-01T18", "2021-01-04T00", []string{"/Christmas_Holidays", "/2021-Q1"}},
		}, got)
	}
}

func TestTimeline_AppendAlbum(t *testing.T) {
	a := assert.New(t)
	const owner = "ironman"

	albums := []*Album{
		{
			Owner:      owner,
			FolderName: "2020-Q3",
			Start:      time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
			End:        time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Owner:      owner,
			FolderName: "/2020-Q4",
			Start:      time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
			End:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Owner:      owner,
			FolderName: "/Christmas_Holidays",
			Start:      time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC),
			End:        time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
		},
	}

	reference, err := NewTimeline(albums)
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	timeline, err := NewTimeline([]*Album{albums[0]})
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	timeline, err = timeline.AppendAlbum(albums[2])
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	timeline, err = timeline.AppendAlbum(albums[1])
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	a.Equal(reference, timeline, "it should create the same timeline than if if it was created with all albums from the beginning")
}

func newSimplifiedSegment(folder, start, end string) simplifiedSegment {
	return simplifiedSegment{
		folder: folder,
		start:  start,
		end:    end,
	}
}

func newAlbum(folder, start, end string) *Album {
	startTime, err := time.Parse(layout, start)
	if err != nil {
		panic(err)
	}
	endTime, err := time.Parse(layout, end)
	if err != nil {
		panic(err)
	}
	return &Album{
		Owner:      "stark",
		Name:       folder,
		FolderName: folder,
		Start:      startTime,
		End:        endTime,
	}
}
