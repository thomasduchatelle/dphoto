// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// RunnerAnalyserFunc is an autogenerated mock type for the RunnerAnalyserFunc type
type RunnerAnalyserFunc struct {
	mock.Mock
}

type RunnerAnalyserFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *RunnerAnalyserFunc) EXPECT() *RunnerAnalyserFunc_Expecter {
	return &RunnerAnalyserFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: found, progressChannel
func (_m *RunnerAnalyserFunc) Execute(found backup.FoundMedia, progressChannel chan *backup.ProgressEvent) (*backup.AnalysedMedia, error) {
	ret := _m.Called(found, progressChannel)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 *backup.AnalysedMedia
	var r1 error
	if rf, ok := ret.Get(0).(func(backup.FoundMedia, chan *backup.ProgressEvent) (*backup.AnalysedMedia, error)); ok {
		return rf(found, progressChannel)
	}
	if rf, ok := ret.Get(0).(func(backup.FoundMedia, chan *backup.ProgressEvent) *backup.AnalysedMedia); ok {
		r0 = rf(found, progressChannel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*backup.AnalysedMedia)
		}
	}

	if rf, ok := ret.Get(1).(func(backup.FoundMedia, chan *backup.ProgressEvent) error); ok {
		r1 = rf(found, progressChannel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RunnerAnalyserFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type RunnerAnalyserFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - found backup.FoundMedia
//   - progressChannel chan *backup.ProgressEvent
func (_e *RunnerAnalyserFunc_Expecter) Execute(found interface{}, progressChannel interface{}) *RunnerAnalyserFunc_Execute_Call {
	return &RunnerAnalyserFunc_Execute_Call{Call: _e.mock.On("Execute", found, progressChannel)}
}

func (_c *RunnerAnalyserFunc_Execute_Call) Run(run func(found backup.FoundMedia, progressChannel chan *backup.ProgressEvent)) *RunnerAnalyserFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(backup.FoundMedia), args[1].(chan *backup.ProgressEvent))
	})
	return _c
}

func (_c *RunnerAnalyserFunc_Execute_Call) Return(_a0 *backup.AnalysedMedia, _a1 error) *RunnerAnalyserFunc_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RunnerAnalyserFunc_Execute_Call) RunAndReturn(run func(backup.FoundMedia, chan *backup.ProgressEvent) (*backup.AnalysedMedia, error)) *RunnerAnalyserFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewRunnerAnalyserFunc creates a new instance of RunnerAnalyserFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRunnerAnalyserFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *RunnerAnalyserFunc {
	mock := &RunnerAnalyserFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
