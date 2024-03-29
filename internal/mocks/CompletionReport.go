// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// CompletionReport is an autogenerated mock type for the CompletionReport type
type CompletionReport struct {
	mock.Mock
}

// CountPerAlbum provides a mock function with given fields:
func (_m *CompletionReport) CountPerAlbum() map[string]*backup.TypeCounter {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for CountPerAlbum")
	}

	var r0 map[string]*backup.TypeCounter
	if rf, ok := ret.Get(0).(func() map[string]*backup.TypeCounter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*backup.TypeCounter)
		}
	}

	return r0
}

// NewAlbums provides a mock function with given fields:
func (_m *CompletionReport) NewAlbums() []string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewAlbums")
	}

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Skipped provides a mock function with given fields:
func (_m *CompletionReport) Skipped() backup.MediaCounter {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Skipped")
	}

	var r0 backup.MediaCounter
	if rf, ok := ret.Get(0).(func() backup.MediaCounter); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(backup.MediaCounter)
	}

	return r0
}

// NewCompletionReport creates a new instance of CompletionReport. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCompletionReport(t interface {
	mock.TestingT
	Cleanup(func())
}) *CompletionReport {
	mock := &CompletionReport{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
