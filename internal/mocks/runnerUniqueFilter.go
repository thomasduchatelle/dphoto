// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// runnerUniqueFilter is an autogenerated mock type for the runnerUniqueFilter type
type runnerUniqueFilter struct {
	mock.Mock
}

// Execute provides a mock function with given fields: medias, progressChannel
func (_m *runnerUniqueFilter) Execute(medias *backup.BackingUpMediaRequest, progressChannel chan *backup.ProgressEvent) bool {
	ret := _m.Called(medias, progressChannel)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*backup.BackingUpMediaRequest, chan *backup.ProgressEvent) bool); ok {
		r0 = rf(medias, progressChannel)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// newRunnerUniqueFilter creates a new instance of runnerUniqueFilter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newRunnerUniqueFilter(t interface {
	mock.TestingT
	Cleanup(func())
}) *runnerUniqueFilter {
	mock := &runnerUniqueFilter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
