// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	backupmodel "github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	mock "github.com/stretchr/testify/mock"
)

// TrackAnalysed is an autogenerated mock type for the TrackAnalysed type
type TrackAnalysed struct {
	mock.Mock
}

// OnAnalysed provides a mock function with given fields: done, total
func (_m *TrackAnalysed) OnAnalysed(done backupmodel.MediaCounter, total backupmodel.MediaCounter) {
	_m.Called(done, total)
}