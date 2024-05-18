// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MediaCounterFunc is an autogenerated mock type for the MediaCounterFunc type
type MediaCounterFunc struct {
	mock.Mock
}

type MediaCounterFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *MediaCounterFunc) EXPECT() *MediaCounterFunc_Expecter {
	return &MediaCounterFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, album
func (_m *MediaCounterFunc) Execute(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	_va := make([]interface{}, len(album))
	for _i := range album {
		_va[_i] = album[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
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

// MediaCounterFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MediaCounterFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - album ...catalog.AlbumId
func (_e *MediaCounterFunc_Expecter) Execute(ctx interface{}, album ...interface{}) *MediaCounterFunc_Execute_Call {
	return &MediaCounterFunc_Execute_Call{Call: _e.mock.On("Execute",
		append([]interface{}{ctx}, album...)...)}
}

func (_c *MediaCounterFunc_Execute_Call) Run(run func(ctx context.Context, album ...catalog.AlbumId)) *MediaCounterFunc_Execute_Call {
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

func (_c *MediaCounterFunc_Execute_Call) Return(_a0 map[catalog.AlbumId]int, _a1 error) *MediaCounterFunc_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MediaCounterFunc_Execute_Call) RunAndReturn(run func(context.Context, ...catalog.AlbumId) (map[catalog.AlbumId]int, error)) *MediaCounterFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMediaCounterFunc creates a new instance of MediaCounterFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMediaCounterFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *MediaCounterFunc {
	mock := &MediaCounterFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
