// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// SourceVolume is an autogenerated mock type for the SourceVolume type
type SourceVolume struct {
	mock.Mock
}

// FindMedias provides a mock function with given fields:
func (_m *SourceVolume) FindMedias() ([]backup.FoundMedia, error) {
	ret := _m.Called()

	var r0 []backup.FoundMedia
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]backup.FoundMedia, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []backup.FoundMedia); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]backup.FoundMedia)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// String provides a mock function with given fields:
func (_m *SourceVolume) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewSourceVolume interface {
	mock.TestingT
	Cleanup(func())
}

// NewSourceVolume creates a new instance of SourceVolume. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSourceVolume(t mockConstructorTestingTNewSourceVolume) *SourceVolume {
	mock := &SourceVolume{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
