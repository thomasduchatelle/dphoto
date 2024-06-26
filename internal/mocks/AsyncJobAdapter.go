// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	archive "github.com/thomasduchatelle/dphoto/pkg/archive"
)

// AsyncJobAdapter is an autogenerated mock type for the AsyncJobAdapter type
type AsyncJobAdapter struct {
	mock.Mock
}

type AsyncJobAdapter_Expecter struct {
	mock *mock.Mock
}

func (_m *AsyncJobAdapter) EXPECT() *AsyncJobAdapter_Expecter {
	return &AsyncJobAdapter_Expecter{mock: &_m.Mock}
}

// LoadImagesInCache provides a mock function with given fields: images
func (_m *AsyncJobAdapter) LoadImagesInCache(images ...*archive.ImageToResize) error {
	_va := make([]interface{}, len(images))
	for _i := range images {
		_va[_i] = images[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for LoadImagesInCache")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...*archive.ImageToResize) error); ok {
		r0 = rf(images...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AsyncJobAdapter_LoadImagesInCache_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoadImagesInCache'
type AsyncJobAdapter_LoadImagesInCache_Call struct {
	*mock.Call
}

// LoadImagesInCache is a helper method to define mock.On call
//   - images ...*archive.ImageToResize
func (_e *AsyncJobAdapter_Expecter) LoadImagesInCache(images ...interface{}) *AsyncJobAdapter_LoadImagesInCache_Call {
	return &AsyncJobAdapter_LoadImagesInCache_Call{Call: _e.mock.On("LoadImagesInCache",
		append([]interface{}{}, images...)...)}
}

func (_c *AsyncJobAdapter_LoadImagesInCache_Call) Run(run func(images ...*archive.ImageToResize)) *AsyncJobAdapter_LoadImagesInCache_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*archive.ImageToResize, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(*archive.ImageToResize)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *AsyncJobAdapter_LoadImagesInCache_Call) Return(_a0 error) *AsyncJobAdapter_LoadImagesInCache_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AsyncJobAdapter_LoadImagesInCache_Call) RunAndReturn(run func(...*archive.ImageToResize) error) *AsyncJobAdapter_LoadImagesInCache_Call {
	_c.Call.Return(run)
	return _c
}

// WarmUpCacheByFolder provides a mock function with given fields: owner, missedStoreKey, width
func (_m *AsyncJobAdapter) WarmUpCacheByFolder(owner string, missedStoreKey string, width int) error {
	ret := _m.Called(owner, missedStoreKey, width)

	if len(ret) == 0 {
		panic("no return value specified for WarmUpCacheByFolder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int) error); ok {
		r0 = rf(owner, missedStoreKey, width)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AsyncJobAdapter_WarmUpCacheByFolder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WarmUpCacheByFolder'
type AsyncJobAdapter_WarmUpCacheByFolder_Call struct {
	*mock.Call
}

// WarmUpCacheByFolder is a helper method to define mock.On call
//   - owner string
//   - missedStoreKey string
//   - width int
func (_e *AsyncJobAdapter_Expecter) WarmUpCacheByFolder(owner interface{}, missedStoreKey interface{}, width interface{}) *AsyncJobAdapter_WarmUpCacheByFolder_Call {
	return &AsyncJobAdapter_WarmUpCacheByFolder_Call{Call: _e.mock.On("WarmUpCacheByFolder", owner, missedStoreKey, width)}
}

func (_c *AsyncJobAdapter_WarmUpCacheByFolder_Call) Run(run func(owner string, missedStoreKey string, width int)) *AsyncJobAdapter_WarmUpCacheByFolder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(int))
	})
	return _c
}

func (_c *AsyncJobAdapter_WarmUpCacheByFolder_Call) Return(_a0 error) *AsyncJobAdapter_WarmUpCacheByFolder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AsyncJobAdapter_WarmUpCacheByFolder_Call) RunAndReturn(run func(string, string, int) error) *AsyncJobAdapter_WarmUpCacheByFolder_Call {
	_c.Call.Return(run)
	return _c
}

// NewAsyncJobAdapter creates a new instance of AsyncJobAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAsyncJobAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *AsyncJobAdapter {
	mock := &AsyncJobAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
