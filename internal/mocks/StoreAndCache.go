// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	io "io"

	archive "github.com/thomasduchatelle/dphoto/pkg/archive"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// StoreAndCache is an autogenerated mock type for the StoreAndCache type
type StoreAndCache struct {
	mock.Mock
}

// Copy provides a mock function with given fields: origin, destination
func (_m *StoreAndCache) Copy(origin string, destination archive.DestructuredKey) (string, error) {
	ret := _m.Called(origin, destination)

	if len(ret) == 0 {
		panic("no return value specified for Copy")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, archive.DestructuredKey) (string, error)); ok {
		return rf(origin, destination)
	}
	if rf, ok := ret.Get(0).(func(string, archive.DestructuredKey) string); ok {
		r0 = rf(origin, destination)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, archive.DestructuredKey) error); ok {
		r1 = rf(origin, destination)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: locations
func (_m *StoreAndCache) Delete(locations []string) error {
	ret := _m.Called(locations)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(locations)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Download provides a mock function with given fields: key
func (_m *StoreAndCache) Download(key string) (io.ReadCloser, error) {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Download")
	}

	var r0 io.ReadCloser
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (io.ReadCloser, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) io.ReadCloser); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: key
func (_m *StoreAndCache) Get(key string) (io.ReadCloser, int, string, error) {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 io.ReadCloser
	var r1 int
	var r2 string
	var r3 error
	if rf, ok := ret.Get(0).(func(string) (io.ReadCloser, int, string, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) io.ReadCloser); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	if rf, ok := ret.Get(1).(func(string) int); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(string) string); ok {
		r2 = rf(key)
	} else {
		r2 = ret.Get(2).(string)
	}

	if rf, ok := ret.Get(3).(func(string) error); ok {
		r3 = rf(key)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// Put provides a mock function with given fields: key, mediaType, content
func (_m *StoreAndCache) Put(key string, mediaType string, content io.Reader) error {
	ret := _m.Called(key, mediaType, content)

	if len(ret) == 0 {
		panic("no return value specified for Put")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, io.Reader) error); ok {
		r0 = rf(key, mediaType, content)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SignedURL provides a mock function with given fields: key, duration
func (_m *StoreAndCache) SignedURL(key string, duration time.Duration) (string, error) {
	ret := _m.Called(key, duration)

	if len(ret) == 0 {
		panic("no return value specified for SignedURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, time.Duration) (string, error)); ok {
		return rf(key, duration)
	}
	if rf, ok := ret.Get(0).(func(string, time.Duration) string); ok {
		r0 = rf(key, duration)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, time.Duration) error); ok {
		r1 = rf(key, duration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upload provides a mock function with given fields: values, content
func (_m *StoreAndCache) Upload(values archive.DestructuredKey, content io.Reader) (string, error) {
	ret := _m.Called(values, content)

	if len(ret) == 0 {
		panic("no return value specified for Upload")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(archive.DestructuredKey, io.Reader) (string, error)); ok {
		return rf(values, content)
	}
	if rf, ok := ret.Get(0).(func(archive.DestructuredKey, io.Reader) string); ok {
		r0 = rf(values, content)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(archive.DestructuredKey, io.Reader) error); ok {
		r1 = rf(values, content)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WalkCacheByPrefix provides a mock function with given fields: prefix, observer
func (_m *StoreAndCache) WalkCacheByPrefix(prefix string, observer func(string)) error {
	ret := _m.Called(prefix, observer)

	if len(ret) == 0 {
		panic("no return value specified for WalkCacheByPrefix")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, func(string)) error); ok {
		r0 = rf(prefix, observer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStoreAndCache creates a new instance of StoreAndCache. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStoreAndCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *StoreAndCache {
	mock := &StoreAndCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
