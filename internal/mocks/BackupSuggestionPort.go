// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	ui "github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
)

// BackupSuggestionPort is an autogenerated mock type for the BackupSuggestionPort type
type BackupSuggestionPort struct {
	mock.Mock
}

type BackupSuggestionPort_Expecter struct {
	mock *mock.Mock
}

func (_m *BackupSuggestionPort) EXPECT() *BackupSuggestionPort_Expecter {
	return &BackupSuggestionPort_Expecter{mock: &_m.Mock}
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

// BackupSuggestionPort_BackupSuggestion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BackupSuggestion'
type BackupSuggestionPort_BackupSuggestion_Call struct {
	*mock.Call
}

// BackupSuggestion is a helper method to define mock.On call
//   - record *ui.SuggestionRecord
//   - existing *ui.ExistingRecord
//   - listener ui.InteractiveRendererPort
func (_e *BackupSuggestionPort_Expecter) BackupSuggestion(record interface{}, existing interface{}, listener interface{}) *BackupSuggestionPort_BackupSuggestion_Call {
	return &BackupSuggestionPort_BackupSuggestion_Call{Call: _e.mock.On("BackupSuggestion", record, existing, listener)}
}

func (_c *BackupSuggestionPort_BackupSuggestion_Call) Run(run func(record *ui.SuggestionRecord, existing *ui.ExistingRecord, listener ui.InteractiveRendererPort)) *BackupSuggestionPort_BackupSuggestion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*ui.SuggestionRecord), args[1].(*ui.ExistingRecord), args[2].(ui.InteractiveRendererPort))
	})
	return _c
}

func (_c *BackupSuggestionPort_BackupSuggestion_Call) Return(_a0 error) *BackupSuggestionPort_BackupSuggestion_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BackupSuggestionPort_BackupSuggestion_Call) RunAndReturn(run func(*ui.SuggestionRecord, *ui.ExistingRecord, ui.InteractiveRendererPort) error) *BackupSuggestionPort_BackupSuggestion_Call {
	_c.Call.Return(run)
	return _c
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
