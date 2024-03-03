// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	ui "github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
)

// BackupSuggestionPort is an autogenerated mock type for the BackupSuggestionPort type
type BackupSuggestionPort struct {
	mock.Mock
}

// BackupSuggestion provides a mock function with given fields: record, existing, listener
func (_m *BackupSuggestionPort) BackupSuggestion(record *ui.SuggestionRecord, existing *ui.ExistingRecord, listener ui.InteractiveRendererPort) error {
	ret := _m.Called(record, existing, listener)

	if len(ret) == 0 {
		panic("no return value specified for BackupSuggestion")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*ui.SuggestionRecord, *ui.ExistingRecord, ui.InteractiveRendererPort) error); ok {
		r0 = rf(record, existing, listener)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewBackupSuggestionPort creates a new instance of BackupSuggestionPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBackupSuggestionPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *BackupSuggestionPort {
	mock := &BackupSuggestionPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
