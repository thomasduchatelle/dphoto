// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	io "io"

	archive "github.com/thomasduchatelle/dphoto/pkg/archive"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// StoreAdapter is an autogenerated mock type for the StoreAdapter type
type StoreAdapter struct {
	mock.Mock
}

// Copy provides a mock function with given fields: origin, destination
func (_m *StoreAdapter) Copy(origin string, destination archive.DestructuredKey) (string, error) {
	ret := _m.Called(origin, destination)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, archive.DestructuredKey) string); ok {
		r0 = rf(origin, destination)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, archive.DestructuredKey) error); ok {
		r1 = rf(origin, destination)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: locations
func (_m *StoreAdapter) Delete(locations []string) error {
	ret := _m.Called(locations)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(locations)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Download provides a mock function with given fields: key
func (_m *StoreAdapter) Download(key string) (io.ReadCloser, error) {
	ret := _m.Called(key)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(string) io.ReadCloser); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignedURL provides a mock function with given fields: key, duration
func (_m *StoreAdapter) SignedURL(key string, duration time.Duration) (string, error) {
	ret := _m.Called(key, duration)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, time.Duration) string); ok {
		r0 = rf(key, duration)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, time.Duration) error); ok {
		r1 = rf(key, duration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upload provides a mock function with given fields: values, content
func (_m *StoreAdapter) Upload(values archive.DestructuredKey, content io.Reader) (string, error) {
	ret := _m.Called(values, content)

	var r0 string
	if rf, ok := ret.Get(0).(func(archive.DestructuredKey, io.Reader) string); ok {
		r0 = rf(values, content)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(archive.DestructuredKey, io.Reader) error); ok {
		r1 = rf(values, content)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewStoreAdapter interface {
	mock.TestingT
	Cleanup(func())
}

// NewStoreAdapter creates a new instance of StoreAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStoreAdapter(t mockConstructorTestingTNewStoreAdapter) *StoreAdapter {
	mock := &StoreAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
