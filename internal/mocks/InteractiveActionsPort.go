// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"

	ui "github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
)

// InteractiveActionsPort is an autogenerated mock type for the InteractiveActionsPort type
type InteractiveActionsPort struct {
	mock.Mock
}

// BackupSuggestion provides a mock function with given fields: record, existing, listener
func (_m *InteractiveActionsPort) BackupSuggestion(record *ui.SuggestionRecord, existing *ui.ExistingRecord, listener ui.InteractiveRendererPort) error {
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

// Create provides a mock function with given fields: createRequest
func (_m *InteractiveActionsPort) Create(createRequest ui.RecordCreation) error {
	ret := _m.Called(createRequest)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(ui.RecordCreation) error); ok {
		r0 = rf(createRequest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAlbum provides a mock function with given fields: folderName
func (_m *InteractiveActionsPort) DeleteAlbum(folderName string) error {
	ret := _m.Called(folderName)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAlbum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(folderName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RenameAlbum provides a mock function with given fields: folderName, newName, renameFolder
func (_m *InteractiveActionsPort) RenameAlbum(folderName string, newName string, renameFolder bool) error {
	ret := _m.Called(folderName, newName, renameFolder)

	if len(ret) == 0 {
		panic("no return value specified for RenameAlbum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, bool) error); ok {
		r0 = rf(folderName, newName, renameFolder)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAlbum provides a mock function with given fields: folderName, start, end
func (_m *InteractiveActionsPort) UpdateAlbum(folderName string, start time.Time, end time.Time) error {
	ret := _m.Called(folderName, start, end)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAlbum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, time.Time, time.Time) error); ok {
		r0 = rf(folderName, start, end)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewInteractiveActionsPort creates a new instance of InteractiveActionsPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInteractiveActionsPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *InteractiveActionsPort {
	mock := &InteractiveActionsPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
