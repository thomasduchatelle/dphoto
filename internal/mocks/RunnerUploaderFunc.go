// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// RunnerUploaderFunc is an autogenerated mock type for the RunnerUploaderFunc type
type RunnerUploaderFunc struct {
	mock.Mock
}

type RunnerUploaderFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *RunnerUploaderFunc) EXPECT() *RunnerUploaderFunc_Expecter {
	return &RunnerUploaderFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: buffer, progressChannel
func (_m *RunnerUploaderFunc) Execute(buffer []*backup.BackingUpMediaRequest, progressChannel chan *backup.ProgressEvent) error {
	ret := _m.Called(buffer, progressChannel)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]*backup.BackingUpMediaRequest, chan *backup.ProgressEvent) error); ok {
		r0 = rf(buffer, progressChannel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunnerUploaderFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type RunnerUploaderFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - buffer []*backup.BackingUpMediaRequest
//   - progressChannel chan *backup.ProgressEvent
func (_e *RunnerUploaderFunc_Expecter) Execute(buffer interface{}, progressChannel interface{}) *RunnerUploaderFunc_Execute_Call {
	return &RunnerUploaderFunc_Execute_Call{Call: _e.mock.On("Execute", buffer, progressChannel)}
}

func (_c *RunnerUploaderFunc_Execute_Call) Run(run func(buffer []*backup.BackingUpMediaRequest, progressChannel chan *backup.ProgressEvent)) *RunnerUploaderFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]*backup.BackingUpMediaRequest), args[1].(chan *backup.ProgressEvent))
	})
	return _c
}

func (_c *RunnerUploaderFunc_Execute_Call) Return(_a0 error) *RunnerUploaderFunc_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RunnerUploaderFunc_Execute_Call) RunAndReturn(run func([]*backup.BackingUpMediaRequest, chan *backup.ProgressEvent) error) *RunnerUploaderFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewRunnerUploaderFunc creates a new instance of RunnerUploaderFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRunnerUploaderFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *RunnerUploaderFunc {
	mock := &RunnerUploaderFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
