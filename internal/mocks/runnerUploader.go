// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// RunnerUploader is an autogenerated mock type for the RunnerUploader type
type RunnerUploader struct {
	mock.Mock
}

type RunnerUploader_Expecter struct {
	mock *mock.Mock
}

func (_m *RunnerUploader) EXPECT() *RunnerUploader_Expecter {
	return &RunnerUploader_Expecter{mock: &_m.Mock}
}

// Upload provides a mock function with given fields: buffer, progressChannel
func (_m *RunnerUploader) Upload(buffer []*backup.BackingUpMediaRequest, progressChannel chan *backup.ProgressEvent) error {
	ret := _m.Called(buffer, progressChannel)

	if len(ret) == 0 {
		panic("no return value specified for Upload")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]*backup.BackingUpMediaRequest, chan *backup.ProgressEvent) error); ok {
		r0 = rf(buffer, progressChannel)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunnerUploader_Upload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upload'
type RunnerUploader_Upload_Call struct {
	*mock.Call
}

// Upload is a helper method to define mock.On call
//   - buffer []*backup.BackingUpMediaRequest
//   - progressChannel chan *backup.ProgressEvent
func (_e *RunnerUploader_Expecter) Upload(buffer interface{}, progressChannel interface{}) *RunnerUploader_Upload_Call {
	return &RunnerUploader_Upload_Call{Call: _e.mock.On("Upload", buffer, progressChannel)}
}

func (_c *RunnerUploader_Upload_Call) Run(run func(buffer []*backup.BackingUpMediaRequest, progressChannel chan *backup.ProgressEvent)) *RunnerUploader_Upload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]*backup.BackingUpMediaRequest), args[1].(chan *backup.ProgressEvent))
	})
	return _c
}

func (_c *RunnerUploader_Upload_Call) Return(_a0 error) *RunnerUploader_Upload_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RunnerUploader_Upload_Call) RunAndReturn(run func([]*backup.BackingUpMediaRequest, chan *backup.ProgressEvent) error) *RunnerUploader_Upload_Call {
	_c.Call.Return(run)
	return _c
}

// NewRunnerUploader creates a new instance of RunnerUploader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRunnerUploader(t interface {
	mock.TestingT
	Cleanup(func())
}) *RunnerUploader {
	mock := &RunnerUploader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
