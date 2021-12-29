// Code generated by mockery 2.9.4. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	ui "github.com/thomasduchatelle/dphoto/dphoto/cmd/ui"
)

// BackupSuggestionPort is an autogenerated mock type for the BackupSuggestionPort type
type BackupSuggestionPort struct {
	mock.Mock
}

// BackupSuggestion provides a mock function with given fields: record, existing, listener
func (_m *BackupSuggestionPort) BackupSuggestion(record *ui.SuggestionRecord, existing *ui.ExistingRecord, listener ui.InteractiveRendererPort) error {
	ret := _m.Called(record, existing, listener)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ui.SuggestionRecord, *ui.ExistingRecord, ui.InteractiveRendererPort) error); ok {
		r0 = rf(record, existing, listener)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}