package ui

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInteractiveSession_MoveDown(t *testing.T) {
	a := assert.New(t)

	const downAction = "DOWN: next page"
	const upAction = "UP: previous page"
	tests := []struct {
		name         string
		initialState InteractiveViewState
		expected     InteractiveViewState
	}{
		{
			name:         "it should take second line from the first [no paging]",
			initialState: InteractiveViewState{newTestRecords(0, 42, 0), nil},
			expected:     InteractiveViewState{newTestRecords(1, 42, 0), newNoSuggestionActions()},
		},
		{
			name:         "it should loop to the first when on the last [no paging]",
			initialState: InteractiveViewState{newTestRecords(8, 42, 0), nil},
			expected:     InteractiveViewState{newTestRecords(0, 42, 0), newNoSuggestionActions()},
		},
		{
			name:         "it should move the last page when on last of the second page [last page not full]",
			initialState: InteractiveViewState{newTestRecords(7, 4, 4), nil},
			expected:     InteractiveViewState{newTestRecords(8, 4, 8), newSuggestionActions("page 3/3", downAction, upAction)},
		},
		{
			name:         "it should move the second page when on last of the first page [exact page size]",
			initialState: InteractiveViewState{newTestRecords(2, 3, 0), nil},
			expected:     InteractiveViewState{newTestRecords(3, 3, 3), newNoSuggestionActions("page 2/3", downAction, upAction)},
		},
		{
			name:         "it should move the first page when on last of the last page",
			initialState: InteractiveViewState{newTestRecords(8, 4, 8), nil},
			expected:     InteractiveViewState{newTestRecords(0, 4, 0), newNoSuggestionActions("page 1/3", downAction, upAction)},
		},
	}

	for _, tt := range tests {
		sess := InteractiveSession{state: tt.initialState}
		sess.MoveDown()
		a.Equal(tt.expected, sess.state, tt.name)
	}
}

func TestInteractiveSession_MoveUp(t *testing.T) {
	a := assert.New(t)

	const downAction = "DOWN: next page"
	const upAction = "UP: previous page"
	tests := []struct {
		name         string
		initialState InteractiveViewState
		expected     InteractiveViewState
	}{
		{
			name:         "it should take first line from the second [no paging]",
			initialState: InteractiveViewState{newTestRecords(1, 42, 0), nil},
			expected:     InteractiveViewState{newTestRecords(0, 42, 0), newNoSuggestionActions()},
		},
		{
			name:         "it should loop to the last when on the first [no paging]",
			initialState: InteractiveViewState{newTestRecords(0, 42, 0), nil},
			expected:     InteractiveViewState{newTestRecords(8, 42, 0), newSuggestionActions()},
		},
		{
			name:         "it should move the last page when on last of the second page [last page not full]",
			initialState: InteractiveViewState{newTestRecords(8, 4, 8), nil},
			expected:     InteractiveViewState{newTestRecords(7, 4, 4), newSuggestionActions("page 2/3", downAction, upAction)},
		},
		{
			name:         "it should loop to last element and last page when on the first element",
			initialState: InteractiveViewState{newTestRecords(0, 4, 0), nil},
			expected:     InteractiveViewState{newTestRecords(8, 4, 8), newSuggestionActions("page 3/3", downAction, upAction)},
		},
	}

	for _, tt := range tests {
		sess := InteractiveSession{state: tt.initialState}
		sess.MoveUp()
		a.Equal(tt.expected, sess.state, tt.name)
	}
}

func TestInteractiveSession_NextPage(t *testing.T) {
	a := assert.New(t)

	const downAction = "DOWN: next page"
	const upAction = "UP: previous page"
	tests := []struct {
		name         string
		initialState InteractiveViewState
		expected     InteractiveViewState
	}{
		{
			name:         "it should do nothing when no pagination",
			initialState: InteractiveViewState{newTestRecords(2, 42, 0), nil},
			expected:     InteractiveViewState{newTestRecords(2, 42, 0), nil},
		},
		{
			name:         "it should move the next page and select first of the page when on a first page",
			initialState: InteractiveViewState{newTestRecords(1, 3, 0), nil},
			expected:     InteractiveViewState{newTestRecords(3, 3, 3), newNoSuggestionActions("page 2/3", downAction, upAction)},
		},
		{
			name:         "it should loop to first page and first element when on the last page",
			initialState: InteractiveViewState{newTestRecords(8, 3, 6), nil},
			expected:     InteractiveViewState{newTestRecords(0, 3, 0), newNoSuggestionActions("page 1/3", downAction, upAction)},
		},
	}

	for _, tt := range tests {
		sess := InteractiveSession{state: tt.initialState}
		sess.NextPage()
		a.Equal(tt.expected, sess.state, tt.name)
	}
}

func TestInteractiveSession_PreviousPage(t *testing.T) {
	a := assert.New(t)

	const downAction = "DOWN: next page"
	const upAction = "UP: previous page"
	tests := []struct {
		name         string
		initialState InteractiveViewState
		expected     InteractiveViewState
	}{
		{
			name:         "it should do nothing when no pagination",
			initialState: InteractiveViewState{newTestRecords(3, 42, 0), nil},
			expected:     InteractiveViewState{newTestRecords(3, 42, 0), nil},
		},
		{
			name:         "it should move the previous page and select first of the page when on a 3rd page",
			initialState: InteractiveViewState{newTestRecords(7, 3, 6), nil},
			expected:     InteractiveViewState{newTestRecords(3, 3, 3), newNoSuggestionActions("page 2/3", downAction, upAction)},
		},
		{
			name:         "it should loop to last page and first element of the page when on the first page",
			initialState: InteractiveViewState{newTestRecords(1, 3, 1), nil},
			expected:     InteractiveViewState{newTestRecords(6, 3, 6), newSuggestionActions("page 3/3", downAction, upAction)},
		},
	}

	for _, tt := range tests {
		sess := InteractiveSession{state: tt.initialState}
		sess.PreviousPage()
		a.Equal(tt.expected, sess.state, tt.name)
	}
}

var records = []*Record{
	{Suggestion: false, FolderName: "2020-Q1", Name: "Q1 2020", Start: time.Now(), End: time.Now(), Count: 4},
	{Suggestion: false, FolderName: "2020-Q2", Name: "Q2 2020", Start: time.Now(), End: time.Now(), Count: 4},
	{Suggestion: false, FolderName: "2020-Q3", Name: "Q3 2020", Start: time.Now(), End: time.Now(), Count: 4},
	{Suggestion: false, FolderName: "2020-Q4", Name: "Q4 2020", Start: time.Now(), End: time.Now(), Count: 4},
	{Suggestion: true, FolderName: "2021-Q1", Name: "Q1 2021", Start: time.Now(), End: time.Now(), Count: 4},
	{Suggestion: true, FolderName: "2021-Q2", Name: "Q2 2021", Start: time.Now(), End: time.Now(), Count: 4},
	{Suggestion: true, FolderName: "2021-Q3", Name: "Q3 2021", Start: time.Now(), End: time.Now(), Count: 4},
	{Suggestion: true, FolderName: "2021-Q4", Name: "Q4 2021", Start: time.Now(), End: time.Now(), Count: 4},
	{Suggestion: true, FolderName: "2022-Q1", Name: "Q1 2022", Start: time.Now(), End: time.Now(), Count: 4},
}

func newTestRecords(selected, size, first int) RecordsState {
	return RecordsState{
		Records:      records,
		Selected:     selected,
		PageSize:     size,
		FirstElement: first,
	}
}

func newNoSuggestionActions(paging ...string) []string {
	return append(paging, "ESC: exit", "N: new", "DEL: delete", "E: edit name", "D: edit dates")
}

func newSuggestionActions(paging ...string) []string {
	return append(paging, "ESC: exit", "N: new", "C: create", "B: backup")
}
