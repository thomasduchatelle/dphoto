// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// TrackAnalysed is an autogenerated mock type for the TrackAnalysed type
type TrackAnalysed struct {
	mock.Mock
}

// OnAnalysed provides a mock function with given fields: done, total
func (_m *TrackAnalysed) OnAnalysed(done, total, cached backup.MediaCounter) {
	_m.Called(done, total)
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
