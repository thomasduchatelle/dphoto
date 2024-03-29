// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	io "io"

	backup "github.com/thomasduchatelle/dphoto/pkg/backup"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// FoundMedia is an autogenerated mock type for the FoundMedia type
type FoundMedia struct {
	mock.Mock
}

// LastModification provides a mock function with given fields:
func (_m *FoundMedia) LastModification() time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for LastModification")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// MediaPath provides a mock function with given fields:
func (_m *FoundMedia) MediaPath() backup.MediaPath {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MediaPath")
	}

	var r0 backup.MediaPath
	if rf, ok := ret.Get(0).(func() backup.MediaPath); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(backup.MediaPath)
	}

	return r0
}

// ReadMedia provides a mock function with given fields:
func (_m *FoundMedia) ReadMedia() (io.ReadCloser, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadMedia")
	}

	var r0 io.ReadCloser
	var r1 error
	if rf, ok := ret.Get(0).(func() (io.ReadCloser, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() io.ReadCloser); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Size provides a mock function with given fields:
func (_m *FoundMedia) Size() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Size")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// String provides a mock function with given fields:
func (_m *FoundMedia) String() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for String")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewFoundMedia creates a new instance of FoundMedia. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFoundMedia(t interface {
	mock.TestingT
	Cleanup(func())
}) *FoundMedia {
	mock := &FoundMedia{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
