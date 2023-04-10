// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// TimelineAdapter is an autogenerated mock type for the TimelineAdapter type
type TimelineAdapter struct {
	mock.Mock
}

// FindAlbum provides a mock function with given fields: dateTime
func (_m *TimelineAdapter) FindAlbum(dateTime time.Time) (string, bool, error) {
	ret := _m.Called(dateTime)

	var r0 string
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(time.Time) (string, bool, error)); ok {
		return rf(dateTime)
	}
	if rf, ok := ret.Get(0).(func(time.Time) string); ok {
		r0 = rf(dateTime)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(time.Time) bool); ok {
		r1 = rf(dateTime)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(time.Time) error); ok {
		r2 = rf(dateTime)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FindOrCreateAlbum provides a mock function with given fields: mediaTime
func (_m *TimelineAdapter) FindOrCreateAlbum(mediaTime time.Time) (string, bool, error) {
	ret := _m.Called(mediaTime)

	var r0 string
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(time.Time) (string, bool, error)); ok {
		return rf(mediaTime)
	}
	if rf, ok := ret.Get(0).(func(time.Time) string); ok {
		r0 = rf(mediaTime)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(time.Time) bool); ok {
		r1 = rf(mediaTime)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(time.Time) error); ok {
		r2 = rf(mediaTime)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewTimelineAdapter interface {
	mock.TestingT
	Cleanup(func())
}

// NewTimelineAdapter creates a new instance of TimelineAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTimelineAdapter(t mockConstructorTestingTNewTimelineAdapter) *TimelineAdapter {
	mock := &TimelineAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
