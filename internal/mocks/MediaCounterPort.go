// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MediaCounterPort is an autogenerated mock type for the MediaCounterPort type
type MediaCounterPort struct {
	mock.Mock
}

type MediaCounterPort_Expecter struct {
	mock *mock.Mock
}

func (_m *MediaCounterPort) EXPECT() *MediaCounterPort_Expecter {
	return &MediaCounterPort_Expecter{mock: &_m.Mock}
}

// CountMedia provides a mock function with given fields: ctx, album
func (_m *MediaCounterPort) CountMedia(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	_va := make([]interface{}, len(album))
	for _i := range album {
		_va[_i] = album[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CountMedia")
	}

	var r0 map[catalog.AlbumId]int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...catalog.AlbumId) (map[catalog.AlbumId]int, error)); ok {
		return rf(ctx, album...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...catalog.AlbumId) map[catalog.AlbumId]int); ok {
		r0 = rf(ctx, album...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[catalog.AlbumId]int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...catalog.AlbumId) error); ok {
		r1 = rf(ctx, album...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MediaCounterPort_CountMedia_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountMedia'
type MediaCounterPort_CountMedia_Call struct {
	*mock.Call
}

// CountMedia is a helper method to define mock.On call
//   - ctx context.Context
//   - album ...catalog.AlbumId
func (_e *MediaCounterPort_Expecter) CountMedia(ctx interface{}, album ...interface{}) *MediaCounterPort_CountMedia_Call {
	return &MediaCounterPort_CountMedia_Call{Call: _e.mock.On("CountMedia",
		append([]interface{}{ctx}, album...)...)}
}

func (_c *MediaCounterPort_CountMedia_Call) Run(run func(ctx context.Context, album ...catalog.AlbumId)) *MediaCounterPort_CountMedia_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]catalog.AlbumId, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(catalog.AlbumId)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *MediaCounterPort_CountMedia_Call) Return(_a0 map[catalog.AlbumId]int, _a1 error) *MediaCounterPort_CountMedia_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MediaCounterPort_CountMedia_Call) RunAndReturn(run func(context.Context, ...catalog.AlbumId) (map[catalog.AlbumId]int, error)) *MediaCounterPort_CountMedia_Call {
	_c.Call.Return(run)
	return _c
}

// NewMediaCounterPort creates a new instance of MediaCounterPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMediaCounterPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *MediaCounterPort {
	mock := &MediaCounterPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}