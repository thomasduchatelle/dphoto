// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// TrackUploaded is an autogenerated mock type for the TrackUploaded type
type TrackUploaded struct {
	mock.Mock
}

// OnUploaded provides a mock function with given fields: done, total
func (_m *TrackUploaded) OnUploaded(done backup.MediaCounter, total backup.MediaCounter) {
	_m.Called(done, total)
}

// NewTrackUploaded creates a new instance of TrackUploaded. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTrackUploaded(t interface {
	mock.TestingT
	Cleanup(func())
}) *TrackUploaded {
	mock := &TrackUploaded{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
