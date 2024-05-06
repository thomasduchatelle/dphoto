// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// TrackAnalysed is an autogenerated mock type for the TrackAnalysed type
type TrackAnalysed struct {
	mock.Mock
}

type TrackAnalysed_Expecter struct {
	mock *mock.Mock
}

func (_m *TrackAnalysed) EXPECT() *TrackAnalysed_Expecter {
	return &TrackAnalysed_Expecter{mock: &_m.Mock}
}

// OnAnalysed provides a mock function with given fields: done, total, cached
func (_m *TrackAnalysed) OnAnalysed(done backup.MediaCounter, total backup.MediaCounter, cached backup.MediaCounter) {
	_m.Called(done, total, cached)
}

// TrackAnalysed_OnAnalysed_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnAnalysed'
type TrackAnalysed_OnAnalysed_Call struct {
	*mock.Call
}

// OnAnalysed is a helper method to define mock.On call
//   - done backup.MediaCounter
//   - total backup.MediaCounter
//   - cached backup.MediaCounter
func (_e *TrackAnalysed_Expecter) OnAnalysed(done interface{}, total interface{}, cached interface{}) *TrackAnalysed_OnAnalysed_Call {
	return &TrackAnalysed_OnAnalysed_Call{Call: _e.mock.On("OnAnalysed", done, total, cached)}
}

func (_c *TrackAnalysed_OnAnalysed_Call) Run(run func(done backup.MediaCounter, total backup.MediaCounter, cached backup.MediaCounter)) *TrackAnalysed_OnAnalysed_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(backup.MediaCounter), args[1].(backup.MediaCounter), args[2].(backup.MediaCounter))
	})
	return _c
}

func (_c *TrackAnalysed_OnAnalysed_Call) Return() *TrackAnalysed_OnAnalysed_Call {
	_c.Call.Return()
	return _c
}

func (_c *TrackAnalysed_OnAnalysed_Call) RunAndReturn(run func(backup.MediaCounter, backup.MediaCounter, backup.MediaCounter)) *TrackAnalysed_OnAnalysed_Call {
	_c.Call.Return(run)
	return _c
}

// NewTrackAnalysed creates a new instance of TrackAnalysed. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTrackAnalysed(t interface {
	mock.TestingT
	Cleanup(func())
}) *TrackAnalysed {
	mock := &TrackAnalysed{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
