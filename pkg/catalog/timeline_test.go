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
	next   []string
}

func toSimplifiedSegment(s *segment) simplifiedSegment {
	name := "<empty>"
	if len(s.albums) > 0 {
		name = s.albums[0].FolderName.String()
	}

	var otherFolders []string
	if len(s.albums) > 1 {
		for _, a := range s.albums[1:] {
			otherFolders = append(otherFolders, a.FolderName.String())
		}
	}

	return simplifiedSegment{
		folder: name,
		start:  s.from.Format(layout),
		end:    s.to.Format(layout),
		next:   otherFolders,
	}
}

func TestNewTimeline(t *testing.T) {
	tests := []struct {
		name    string
		args    []*Album
		want    []simplifiedSegment
		wantErr assert.ErrorAssertionFunc
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
			assert.NoError,
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
			assert.NoError,
		},
		{
			"it should support albums overlapping - priority to the first",
			[]*Album{
				newAlbum("/A-01", "2020-12-01T00", "2020-12-11T00"),
				newAlbum("/A-02", "2020-12-10T12", "2020-12-21T12"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/A-01", "2020-12-01T00", "2020-12-11T00", "/A-02"),
				newSimplifiedSegment("/A-02", "2020-12-11T00", "2020-12-21T12"),
			},
			assert.NoError,
		},
		{
			"it should support albums overlapping - priority to the second",
			[]*Album{
				newAlbum("/A-01", "2020-12-01T00", "2020-12-11T00"),
				newAlbum("/A-02", "2020-12-10T12", "2020-12-15T12"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/A-01", "2020-12-01T00", "2020-12-10T12"),
				newSimplifiedSegment("/A-02", "2020-12-10T12", "2020-12-15T12", "/A-01"),
			},
			assert.NoError,
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
				newSimplifiedSegment("/A-02", "2020-12-02T00", "2020-12-03T00", "/A-01"),
				newSimplifiedSegment("/A-03", "2020-12-03T00", "2020-12-04T00", "/A-02", "/A-01"),
				newSimplifiedSegment("/A-02", "2020-12-04T00", "2020-12-05T00", "/A-01"),
				newSimplifiedSegment("/A-01", "2020-12-05T00", "2020-12-06T00"),
			},
			assert.NoError,
		},
		{
			"it should support albums starting at the same time",
			[]*Album{
				newAlbum("/A-01", "2020-12-01T00", "2020-12-06T00"),
				newAlbum("/A-02", "2020-12-01T00", "2020-12-05T00"),
				newAlbum("/A-03", "2020-12-01T00", "2020-12-07T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/A-02", "2020-12-01T00", "2020-12-05T00", "/A-01", "/A-03"),
				newSimplifiedSegment("/A-01", "2020-12-05T00", "2020-12-06T00", "/A-03"),
				newSimplifiedSegment("/A-03", "2020-12-06T00", "2020-12-07T00"),
			},
			assert.NoError,
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
				newSimplifiedSegment("/A-02", "2020-12-02T00", "2020-12-03T00", "/A-01"),
				newSimplifiedSegment("/A-03", "2020-12-03T00", "2020-12-06T00", "/A-02", "/A-01"),
			},
			assert.NoError,
		},
		{
			"it should not have empty segment [2 back to back albums - 01 = 02 > 99]",
			[]*Album{
				newAlbum("/A-01", "2024-01-01T00", "2024-04-01T00"),
				newAlbum("/A-99", "2024-01-05T00", "2024-07-01T00"),
				newAlbum("/A-02", "2024-04-01T00", "2024-07-01T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/A-01", "2024-01-01T00", "2024-04-01T00", "/A-99"),
				newSimplifiedSegment("/A-02", "2024-04-01T00", "2024-07-01T00", "/A-99"),
			},
			assert.NoError,
		},
		{
			"it should have a segment for 99 [2 back to back albums - 01 > 99 > 02]",
			[]*Album{
				newAlbum("/A-01", "2024-02-01T00", "2024-04-01T00"),
				newAlbum("/A-99", "2024-02-01T00", "2024-04-02T00"),
				newAlbum("/A-02", "2024-04-01T00", "2024-07-01T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/A-01", "2024-02-01T00", "2024-04-01T00", "/A-99"),
				newSimplifiedSegment("/A-99", "2024-04-01T00", "2024-04-02T00", "/A-02"),
				newSimplifiedSegment("/A-02", "2024-04-02T00", "2024-07-01T00"),
			},
			assert.NoError,
		},
		{
			"it should have a segment for 99 [2 back to back albums - 02 > 99 > 01]",
			[]*Album{
				newAlbum("/A-01", "2024-01-01T00", "2024-04-01T00"),
				newAlbum("/A-99", "2024-02-01T00", "2024-04-02T00"),
				newAlbum("/A-02", "2024-04-01T00", "2024-05-01T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/A-01", "2024-01-01T00", "2024-02-01T00"),
				newSimplifiedSegment("/A-99", "2024-02-01T00", "2024-04-01T00", "/A-01"),
				newSimplifiedSegment("/A-02", "2024-04-01T00", "2024-05-01T00", "/A-99"),
			},
			assert.NoError,
		},
		{
			"it should have a segment for 99 [2 back to back albums - 99 > 01 = 03]",
			[]*Album{
				newAlbum("/A-01", "2024-01-01T00", "2024-04-01T00"),
				newAlbum("/A-99", "2024-03-20T00", "2024-04-10T00"),
				newAlbum("/A-02", "2024-04-01T00", "2024-07-01T00"),
			},
			[]simplifiedSegment{
				newSimplifiedSegment("/A-01", "2024-01-01T00", "2024-03-20T00"),
				newSimplifiedSegment("/A-99", "2024-03-20T00", "2024-04-10T00", "/A-01", "/A-02"),
				newSimplifiedSegment("/A-02", "2024-04-10T00", "2024-07-01T00"),
			},
			assert.NoError,
		},
		{
			name: "it should not fail if two albums with same dates are added to the timeline",
			args: []*Album{
				newAlbum("/A-01", "2024-01-01T00", "2024-04-01T00"),
				newAlbum("/A-01", "2024-01-01T00", "2024-04-01T00"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, DuplicateError, "got:", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTimeline(tt.args)

			if !tt.wantErr(t, err, tt.name) || err != nil {
				return
			}

			segments := make([]simplifiedSegment, len(got.segments))
			for i, s := range got.segments {
				segments[i] = toSimplifiedSegment(&s)
			}

			assert.Equal(t, tt.want, segments, tt.name)
		})
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
				a.Equal(tt.want.name, album.FolderName.String(), tt.name)
			}

			var albums []string
			for _, a := range timeline.FindAllAt(tt.args) {
				albums = append(albums, a.FolderName.String())
			}
			sort.Strings(albums)

			a.Equal(tt.want.allNames, albums, tt.name)
		}
	}
}

func TestFindAt_FindBetween(t *testing.T) {
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
			name:       "it should return no segment outside timeline boundaries",
			args:       Args{"2019-01-01T00", "2019-01-01T00"},
			wantMissed: []Want{{start: "2019-01-01T00", end: "2019-01-01T00"}},
		},
		{
			name:       "it should return a segment with dates updated to match the request",
			args:       Args{"2020-01-01T00", "2020-09-01T00"},
			want:       []Want{{"2020-07-01T00", "2020-09-01T00", []string{"/2020-Q3"}}},
			wantMissed: []Want{{start: "2020-01-01T00", end: "2020-07-01T00"}},
		},
		{
			name: "it should return segments within the request",
			args: Args{"2020-12-31T00", "2021-01-02T00"},
			want: []Want{
				{"2020-12-31T00", "2020-12-31T18", []string{"/Christmas_Holidays", "/2020-Q4"}},
				{"2020-12-31T18", "2021-01-01T18", []string{"/New_Year", "/Christmas_Holidays", "/2021-Q1", "/2020-Q4"}},
				{"2021-01-01T18", "2021-01-02T00", []string{"/Christmas_Holidays", "/2021-Q1"}},
			},
		},
		{
			name: "it should return segments within the request",
			args: Args{"2020-12-18T00", "2020-12-26T00"},
			want: []Want{
				{"2020-12-18T00", "2020-12-24T00", []string{"/Christmas_First_Week", "/Christmas_Holidays", "/2020-Q4"}},
				{"2020-12-24T00", "2020-12-26T00", []string{"/Christmas_Day", "/Christmas_First_Week", "/Christmas_Holidays", "/2020-Q4"}},
			},
		},
		{
			name: "it should notice the gap between mars and may, and the missing dates at the end of the request",
			args: Args{"2021-03-23T00", "2021-06-26T00"},
			want: []Want{
				{"2021-03-23T00", "2021-04-01T00", []string{"/2021-Q1"}},
				{"2021-05-01T00", "2021-06-01T00", []string{"/2021-May"}},
			},
			wantMissed: []Want{
				{"2021-04-01T00", "2021-05-01T00", nil},
				{"2021-06-01T00", "2021-06-26T00", nil},
			},
		},
	}

	timeline, err := NewTimeline(AlbumCollection())
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments, missed := timeline.FindBetween(MustParse(layout, tt.args.start), MustParse(layout, tt.args.end))
			var got []Want
			var gotMissed []Want
			for _, seg := range segments {
				var names []string
				for _, a := range seg.Albums {
					names = append(names, a.FolderName.String())
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

			assert.Equal(t, tt.want, got, tt.name)
			assert.Equal(t, tt.wantMissed, gotMissed, tt.name)
		})
	}
}

func TestTimeline_FindForAlbum(t1 *testing.T) {
	collection := AlbumCollection()
	christmasHolidays := collection[2]

	type Want struct {
		start    string
		end      string
		allNames []string
	}

	type fields struct {
		AlbumCollection []*Album
	}
	type args struct {
		albumId AlbumId
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantSegments []Want
	}{
		{
			name: "it should return no segment if there is no albums",
			fields: fields{
				AlbumCollection: nil,
			},
			args: args{
				albumId: christmasHolidays.AlbumId,
			},
			wantSegments: nil,
		},
		{
			name: "it should return one segment if there is a single albums",
			fields: fields{
				AlbumCollection: []*Album{christmasHolidays},
			},
			args: args{
				albumId: christmasHolidays.AlbumId,
			},
			wantSegments: []Want{
				{"2020-12-18T00", "2021-01-04T00", []string{"/Christmas_Holidays"}},
			},
		},
		{
			name: "it should return one segment when another quarter is added",
			fields: fields{
				AlbumCollection: []*Album{collection[0], collection[1]},
			},
			args: args{
				albumId: collection[1].AlbumId,
			},
			wantSegments: []Want{
				{"2020-10-01T00", "2021-01-01T00", []string{"/2020-Q4"}},
			},
		},
		{
			name: "it should find all segments for the Christmas album",
			fields: fields{
				AlbumCollection: collection,
			},
			args: args{
				albumId: christmasHolidays.AlbumId,
			},
			wantSegments: []Want{
				{"2020-12-26T00", "2020-12-31T18", []string{"/Christmas_Holidays", "/2020-Q4"}},
				{"2021-01-01T18", "2021-01-04T00", []string{"/Christmas_Holidays", "/2021-Q1"}},
			},
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			timeline, err := NewTimeline(tt.fields.AlbumCollection)
			if !assert.NoError(t, err) {
				return
			}

			gotSegments := timeline.FindForAlbum(tt.args.albumId)

			var got []Want
			for _, seg := range gotSegments {
				var names []string
				for _, a := range seg.Albums {
					names = append(names, a.FolderName.String())
				}

				got = append(got, Want{
					start:    seg.Start.Format(layout),
					end:      seg.End.Format(layout),
					allNames: names,
				})
			}
			assert.Equalf(t, tt.wantSegments, got, "FindForAlbum(%v)", tt.args.albumId)
		})
	}
}

func TestTimeline_AppendAlbum(t *testing.T) {
	a := assert.New(t)
	const owner = "ironman"

	albums := []*Album{
		{
			AlbumId: AlbumId{
				Owner:      owner,
				FolderName: NewFolderName("2020-Q3"),
			},
			Start: time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			AlbumId: AlbumId{
				Owner:      owner,
				FolderName: NewFolderName("/2020-Q4"),
			},
			Start: time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			AlbumId: AlbumId{
				Owner:      owner,
				FolderName: NewFolderName("/Christmas_Holidays"),
			},
			Start: time.Date(2020, 12, 18, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
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

func newSimplifiedSegment(folder, start, end string, next ...string) simplifiedSegment {
	return simplifiedSegment{
		folder: folder,
		start:  start,
		end:    end,
		next:   next,
	}
}

func newSegment(start, end time.Time, albums ...*Album) segment {
	return segment{
		from:   start,
		to:     end,
		albums: albums,
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
		AlbumId: AlbumId{
			Owner:      "stark",
			FolderName: NewFolderName(folder),
		},
		Name:  folder,
		Start: startTime,
		End:   endTime,
	}
}

func TestTimeline_FindSegmentsBetween(t1 *testing.T) {
	jan23 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	jan24 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	jul24 := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	oct24 := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	jan25 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	jan26 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	album24 := Album{
		AlbumId: AlbumId{
			Owner:      "stark",
			FolderName: NewFolderName("/2024"),
		},
		Name:  "2024",
		Start: jan24,
		End:   jan25,
	}
	album24FirstHalf := Album{
		AlbumId: AlbumId{
			Owner:      "stark",
			FolderName: NewFolderName("/2024-H1"),
		},
		Name:  "2024-H1",
		Start: jan24,
		End:   jul24,
	}
	album24SecondHalf := Album{
		AlbumId: AlbumId{
			Owner:      "stark",
			FolderName: NewFolderName("/2024-H2"),
		},
		Name:  "2024-H2",
		Start: jul24,
		End:   jan25,
	}
	album24q4 := Album{
		AlbumId: AlbumId{
			Owner:      "stark",
			FolderName: NewFolderName("/2024-Q4"),
		},
		Name:  "2024-Q4",
		Start: oct24,
		End:   jan25,
	}

	type fields struct {
		segments []segment
		albums   []*Album
	}
	type args struct {
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantSegments []PrioritySegment
	}{
		{
			name:   "it should return a single segment when the timeline is empty",
			fields: fields{},
			args: args{
				start: jan24,
				end:   jan25,
			},
			wantSegments: []PrioritySegment{
				{
					Start: jan24,
					End:   jan25,
				},
			},
		},
		{
			name: "it should return the single segment when requested range is the same as the segment",
			fields: fields{
				segments: []segment{
					newSegment(jan24, jan25, &album24),
				},
			},
			args: args{
				start: jan24,
				end:   jan25,
			},
			wantSegments: []PrioritySegment{
				{
					Start:  jan24,
					End:    jan25,
					Albums: []Album{album24},
				},
			},
		},
		{
			name: "it should return all segments within the 2 dates",
			fields: fields{
				segments: []segment{
					newSegment(jan24, jul24, &album24FirstHalf),
					newSegment(jul24, jan25, &album24SecondHalf),
				},
			},
			args: args{
				start: jan24,
				end:   jan25,
			},
			wantSegments: []PrioritySegment{
				{
					Start:  jan24,
					End:    jul24,
					Albums: []Album{album24FirstHalf},
				},
				{
					Start:  jul24,
					End:    jan25,
					Albums: []Album{album24SecondHalf},
				},
			},
		},
		{
			name: "it should create missing segments before, between, and after existing segments",
			fields: fields{
				segments: []segment{
					newSegment(jan24, jul24, &album24FirstHalf),
					newSegment(oct24, jan25, &album24q4),
				},
			},
			args: args{
				start: jan23,
				end:   jan26,
			},
			wantSegments: []PrioritySegment{
				{
					Start: jan23,
					End:   jan24,
				},
				{
					Start:  jan24,
					End:    jul24,
					Albums: []Album{album24FirstHalf},
				},
				{
					Start: jul24,
					End:   oct24,
				},
				{
					Start:  oct24,
					End:    jan25,
					Albums: []Album{album24q4},
				},
				{
					Start: jan25,
					End:   jan26,
				},
			},
		},
		{
			name: "it should take partially segments",
			fields: fields{
				segments: []segment{
					newSegment(jan24, jan25, &album24),
				},
			},
			args: args{
				start: jul24,
				end:   oct24,
			},
			wantSegments: []PrioritySegment{
				{
					Start:  jul24,
					End:    oct24,
					Albums: []Album{album24},
				},
			},
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Timeline{
				segments: tt.fields.segments,
				albums:   tt.fields.albums,
			}
			assert.Equalf(t1, tt.wantSegments, t.FindSegmentsBetween(tt.args.start, tt.args.end), "FindSegmentsBetween(%v, %v)", tt.args.start, tt.args.end)
		})
	}
}
