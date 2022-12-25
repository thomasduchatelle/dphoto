package ui

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_createFlattenTree(t *testing.T) {
	a := assert.New(t)

	singleSuggestion := []*SuggestionRecord{
		{
			FolderName: "suggestion-1",
			Name:       "",
			Start:      parseDate("2021-07-20"),
			End:        parseDate("2021-07-24"),
			Distribution: map[string]int{
				"2021-07-20": 4,
				"2021-07-21": 3,
				"2021-07-22": 12,
			},
		},
	}

	existingRecordsToMere := []*ExistingRecord{
		{
			FolderName: "fev-21",
			Name:       "",
			Start:      parseDate("2021-02-01"),
			End:        parseDate("2021-03-01"),
			Count:      2,
			ActivePeriods: []Period{
				{Start: parseDate("2021-02-01"), End: parseDate("2021-03-01")},
			},
		},
		{
			FolderName: "q1-21",
			Name:       "",
			Start:      parseDate("2021-01-01"),
			End:        parseDate("2021-04-01"),
			Count:      6,
			ActivePeriods: []Period{
				{Start: parseDate("2021-01-01"), End: parseDate("2021-02-01")},
				{Start: parseDate("2021-03-01"), End: parseDate("2021-04-01")},
			},
		},
	}
	suggestionRecordsToMerge := []*SuggestionRecord{
		{
			FolderName: "ski 21",
			Name:       "",
			Start:      parseDate("2021-02-12"),
			End:        parseDate("2021-02-19"),
			Distribution: map[string]int{
				"2021-02-12": 4,
				"2021-02-17": 3,
			},
		},
		{
			FolderName: "school q1",
			Name:       "",
			Start:      parseDate("2021-01-04"),
			End:        parseDate("2021-04-12"),
			Distribution: map[string]int{
				"2021-01-04": 1,
				"2021-02-17": 10,
				"2021-03-05": 100,
				"2021-04-10": 1000,
			},
		},
	}

	type args struct {
		existing    []*ExistingRecord
		suggestions []*SuggestionRecord
	}
	wantMergeRecords := []*Record{
		{
			Indent:     0,
			Suggestion: false,
			FolderName: "fev-21",
			Name:       "",
			Start:      parseDate("2021-02-01"),
			End:        parseDate("2021-03-01"),
			Count:      2,
			TotalCount: 2,
		},
		{
			Indent:               1,
			Suggestion:           true,
			FolderName:           "ski 21",
			Name:                 "",
			Start:                parseDate("2021-02-12"),
			End:                  parseDate("2021-02-19"),
			Count:                7,
			TotalCount:           7,
			ParentExistingRecord: existingRecordsToMere[0],
			SuggestionRecord:     suggestionRecordsToMerge[0],
		},
		{
			Indent:               1,
			Suggestion:           true,
			FolderName:           "school q1",
			Name:                 "",
			Start:                parseDate("2021-01-04"),
			End:                  parseDate("2021-04-12"),
			Count:                10,
			TotalCount:           1111,
			ParentExistingRecord: existingRecordsToMere[0],
			SuggestionRecord:     suggestionRecordsToMerge[1],
		},
		{
			Indent:           0,
			Suggestion:       true,
			FolderName:       "school q1",
			Name:             "",
			Start:            parseDate("2021-01-04"),
			End:              parseDate("2021-04-12"),
			Count:            1000,
			TotalCount:       1111,
			SuggestionRecord: suggestionRecordsToMerge[1],
		},
		{
			Indent:     0,
			Suggestion: false,
			FolderName: "q1-21",
			Name:       "",
			Start:      parseDate("2021-01-01"),
			End:        parseDate("2021-04-01"),
			Count:      6,
			TotalCount: 6,
		},
		{
			Indent:               1,
			Suggestion:           true,
			FolderName:           "school q1",
			Name:                 "",
			Start:                parseDate("2021-01-04"),
			End:                  parseDate("2021-04-12"),
			Count:                101,
			TotalCount:           1111,
			ParentExistingRecord: existingRecordsToMere[1],
			SuggestionRecord:     suggestionRecordsToMerge[1],
		},
	}
	tests := []struct {
		name string
		args args
		want []*Record
	}{
		{"it should give a list with only existing records",
			args{
				[]*ExistingRecord{
					{
						FolderName:    "exiting-1",
						Name:          "",
						Start:         parseDate("2021-07-20"),
						End:           parseDate("2021-07-24"),
						Count:         42,
						ActivePeriods: nil,
					},
				},
				nil,
			},
			[]*Record{
				{
					Indent:     0,
					Suggestion: false,
					FolderName: "exiting-1",
					Name:       "",
					Start:      parseDate("2021-07-20"),
					End:        parseDate("2021-07-24"),
					Count:      42,
					TotalCount: 42,
				},
			}},
		{"it should give a list with only suggestions",
			args{
				nil,
				singleSuggestion,
			},
			[]*Record{
				{
					Indent:           0,
					Suggestion:       true,
					FolderName:       "suggestion-1",
					Name:             "",
					Start:            parseDate("2021-07-20"),
					End:              parseDate("2021-07-24"),
					Count:            19,
					TotalCount:       19,
					SuggestionRecord: singleSuggestion[0],
				},
			}},
		{"it should merge suggestions with the existing records",
			args{
				existingRecordsToMere,
				suggestionRecordsToMerge,
			},
			wantMergeRecords},
	}
	for _, tt := range tests {
		got := createFlattenTree(tt.args.existing, tt.args.suggestions)
		a.Equal(tt.want, got)
	}
}

func parseDate(date string) time.Time {
	parse, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return parse
}
