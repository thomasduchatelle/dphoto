// Code generated by mockery v2.43.0. DO NOT EDIT.

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

type StoreAdapter_Expecter struct {
	mock *mock.Mock
}

func (_m *StoreAdapter) EXPECT() *StoreAdapter_Expecter {
	return &StoreAdapter_Expecter{mock: &_m.Mock}
}

// Copy provides a mock function with given fields: origin, destination
func (_m *StoreAdapter) Copy(origin string, destination archive.DestructuredKey) (string, error) {
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

// StoreAdapter_Copy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Copy'
type StoreAdapter_Copy_Call struct {
	*mock.Call
}

// Copy is a helper method to define mock.On call
//   - origin string
//   - destination archive.DestructuredKey
func (_e *StoreAdapter_Expecter) Copy(origin interface{}, destination interface{}) *StoreAdapter_Copy_Call {
	return &StoreAdapter_Copy_Call{Call: _e.mock.On("Copy", origin, destination)}
}

func (_c *StoreAdapter_Copy_Call) Run(run func(origin string, destination archive.DestructuredKey)) *StoreAdapter_Copy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(archive.DestructuredKey))
	})
	return _c
}

func (_c *StoreAdapter_Copy_Call) Return(_a0 string, _a1 error) *StoreAdapter_Copy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StoreAdapter_Copy_Call) RunAndReturn(run func(string, archive.DestructuredKey) (string, error)) *StoreAdapter_Copy_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: locations
func (_m *StoreAdapter) Delete(locations []string) error {
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

// StoreAdapter_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type StoreAdapter_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - locations []string
func (_e *StoreAdapter_Expecter) Delete(locations interface{}) *StoreAdapter_Delete_Call {
	return &StoreAdapter_Delete_Call{Call: _e.mock.On("Delete", locations)}
}

func (_c *StoreAdapter_Delete_Call) Run(run func(locations []string)) *StoreAdapter_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string))
	})
	return _c
}

func (_c *StoreAdapter_Delete_Call) Return(_a0 error) *StoreAdapter_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StoreAdapter_Delete_Call) RunAndReturn(run func([]string) error) *StoreAdapter_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Download provides a mock function with given fields: key
func (_m *StoreAdapter) Download(key string) (io.ReadCloser, error) {
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

// StoreAdapter_Download_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Download'
type StoreAdapter_Download_Call struct {
	*mock.Call
}

// Download is a helper method to define mock.On call
//   - key string
func (_e *StoreAdapter_Expecter) Download(key interface{}) *StoreAdapter_Download_Call {
	return &StoreAdapter_Download_Call{Call: _e.mock.On("Download", key)}
}

func (_c *StoreAdapter_Download_Call) Run(run func(key string)) *StoreAdapter_Download_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *StoreAdapter_Download_Call) Return(_a0 io.ReadCloser, _a1 error) *StoreAdapter_Download_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StoreAdapter_Download_Call) RunAndReturn(run func(string) (io.ReadCloser, error)) *StoreAdapter_Download_Call {
	_c.Call.Return(run)
	return _c
}

// SignedURL provides a mock function with given fields: key, duration
func (_m *StoreAdapter) SignedURL(key string, duration time.Duration) (string, error) {
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

// StoreAdapter_SignedURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SignedURL'
type StoreAdapter_SignedURL_Call struct {
	*mock.Call
}

// SignedURL is a helper method to define mock.On call
//   - key string
//   - duration time.Duration
func (_e *StoreAdapter_Expecter) SignedURL(key interface{}, duration interface{}) *StoreAdapter_SignedURL_Call {
	return &StoreAdapter_SignedURL_Call{Call: _e.mock.On("SignedURL", key, duration)}
}

func (_c *StoreAdapter_SignedURL_Call) Run(run func(key string, duration time.Duration)) *StoreAdapter_SignedURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(time.Duration))
	})
	return _c
}

func (_c *StoreAdapter_SignedURL_Call) Return(_a0 string, _a1 error) *StoreAdapter_SignedURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StoreAdapter_SignedURL_Call) RunAndReturn(run func(string, time.Duration) (string, error)) *StoreAdapter_SignedURL_Call {
	_c.Call.Return(run)
	return _c
}

// Upload provides a mock function with given fields: values, content
func (_m *StoreAdapter) Upload(values archive.DestructuredKey, content io.Reader) (string, error) {
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

// StoreAdapter_Upload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upload'
type StoreAdapter_Upload_Call struct {
	*mock.Call
}

// Upload is a helper method to define mock.On call
//   - values archive.DestructuredKey
//   - content io.Reader
func (_e *StoreAdapter_Expecter) Upload(values interface{}, content interface{}) *StoreAdapter_Upload_Call {
	return &StoreAdapter_Upload_Call{Call: _e.mock.On("Upload", values, content)}
}

func (_c *StoreAdapter_Upload_Call) Run(run func(values archive.DestructuredKey, content io.Reader)) *StoreAdapter_Upload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(archive.DestructuredKey), args[1].(io.Reader))
	})
	return _c
}

func (_c *StoreAdapter_Upload_Call) Return(_a0 string, _a1 error) *StoreAdapter_Upload_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StoreAdapter_Upload_Call) RunAndReturn(run func(archive.DestructuredKey, io.Reader) (string, error)) *StoreAdapter_Upload_Call {
	_c.Call.Return(run)
	return _c
}

// NewStoreAdapter creates a new instance of StoreAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStoreAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *StoreAdapter {
	mock := &StoreAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
