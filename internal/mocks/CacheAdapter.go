// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	io "io"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// CacheAdapter is an autogenerated mock type for the CacheAdapter type
type CacheAdapter struct {
	mock.Mock
}

type CacheAdapter_Expecter struct {
	mock *mock.Mock
}

func (_m *CacheAdapter) EXPECT() *CacheAdapter_Expecter {
	return &CacheAdapter_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: key
func (_m *CacheAdapter) Get(key string) (io.ReadCloser, int, string, error) {
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

// CacheAdapter_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type CacheAdapter_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - key string
func (_e *CacheAdapter_Expecter) Get(key interface{}) *CacheAdapter_Get_Call {
	return &CacheAdapter_Get_Call{Call: _e.mock.On("Get", key)}
}

func (_c *CacheAdapter_Get_Call) Run(run func(key string)) *CacheAdapter_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *CacheAdapter_Get_Call) Return(_a0 io.ReadCloser, _a1 int, _a2 string, _a3 error) *CacheAdapter_Get_Call {
	_c.Call.Return(_a0, _a1, _a2, _a3)
	return _c
}

func (_c *CacheAdapter_Get_Call) RunAndReturn(run func(string) (io.ReadCloser, int, string, error)) *CacheAdapter_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Put provides a mock function with given fields: key, mediaType, content
func (_m *CacheAdapter) Put(key string, mediaType string, content io.Reader) error {
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

// CacheAdapter_Put_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Put'
type CacheAdapter_Put_Call struct {
	*mock.Call
}

// Put is a helper method to define mock.On call
//   - key string
//   - mediaType string
//   - content io.Reader
func (_e *CacheAdapter_Expecter) Put(key interface{}, mediaType interface{}, content interface{}) *CacheAdapter_Put_Call {
	return &CacheAdapter_Put_Call{Call: _e.mock.On("Put", key, mediaType, content)}
}

func (_c *CacheAdapter_Put_Call) Run(run func(key string, mediaType string, content io.Reader)) *CacheAdapter_Put_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(io.Reader))
	})
	return _c
}

func (_c *CacheAdapter_Put_Call) Return(_a0 error) *CacheAdapter_Put_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CacheAdapter_Put_Call) RunAndReturn(run func(string, string, io.Reader) error) *CacheAdapter_Put_Call {
	_c.Call.Return(run)
	return _c
}

// SignedURL provides a mock function with given fields: key, duration
func (_m *CacheAdapter) SignedURL(key string, duration time.Duration) (string, error) {
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

// CacheAdapter_SignedURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SignedURL'
type CacheAdapter_SignedURL_Call struct {
	*mock.Call
}

// SignedURL is a helper method to define mock.On call
//   - key string
//   - duration time.Duration
func (_e *CacheAdapter_Expecter) SignedURL(key interface{}, duration interface{}) *CacheAdapter_SignedURL_Call {
	return &CacheAdapter_SignedURL_Call{Call: _e.mock.On("SignedURL", key, duration)}
}

func (_c *CacheAdapter_SignedURL_Call) Run(run func(key string, duration time.Duration)) *CacheAdapter_SignedURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(time.Duration))
	})
	return _c
}

func (_c *CacheAdapter_SignedURL_Call) Return(_a0 string, _a1 error) *CacheAdapter_SignedURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CacheAdapter_SignedURL_Call) RunAndReturn(run func(string, time.Duration) (string, error)) *CacheAdapter_SignedURL_Call {
	_c.Call.Return(run)
	return _c
}

// WalkCacheByPrefix provides a mock function with given fields: prefix, observer
func (_m *CacheAdapter) WalkCacheByPrefix(prefix string, observer func(string)) error {
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

// CacheAdapter_WalkCacheByPrefix_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WalkCacheByPrefix'
type CacheAdapter_WalkCacheByPrefix_Call struct {
	*mock.Call
}

// WalkCacheByPrefix is a helper method to define mock.On call
//   - prefix string
//   - observer func(string)
func (_e *CacheAdapter_Expecter) WalkCacheByPrefix(prefix interface{}, observer interface{}) *CacheAdapter_WalkCacheByPrefix_Call {
	return &CacheAdapter_WalkCacheByPrefix_Call{Call: _e.mock.On("WalkCacheByPrefix", prefix, observer)}
}

func (_c *CacheAdapter_WalkCacheByPrefix_Call) Run(run func(prefix string, observer func(string))) *CacheAdapter_WalkCacheByPrefix_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(func(string)))
	})
	return _c
}

func (_c *CacheAdapter_WalkCacheByPrefix_Call) Return(_a0 error) *CacheAdapter_WalkCacheByPrefix_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CacheAdapter_WalkCacheByPrefix_Call) RunAndReturn(run func(string, func(string)) error) *CacheAdapter_WalkCacheByPrefix_Call {
	_c.Call.Return(run)
	return _c
}

// NewCacheAdapter creates a new instance of CacheAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCacheAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *CacheAdapter {
	mock := &CacheAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
