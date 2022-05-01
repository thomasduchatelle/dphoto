// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backupmodel "github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"

	testing "testing"
)

// TrackAnalysed is an autogenerated mock type for the TrackAnalysed type
type TrackAnalysed struct {
	mock.Mock
}

// OnAnalysed provides a mock function with given fields: done, total
func (_m *TrackAnalysed) OnAnalysed(done backupmodel.MediaCounter, total backupmodel.MediaCounter) {
	_m.Called(done, total)
}

// NewTrackAnalysed creates a new instance of TrackAnalysed. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTrackAnalysed(t testing.TB) *TrackAnalysed {
	mock := &TrackAnalysed{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
