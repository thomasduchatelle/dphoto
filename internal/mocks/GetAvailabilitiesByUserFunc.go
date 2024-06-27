// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalogviews "github.com/thomasduchatelle/dphoto/pkg/catalogviews"

	mock "github.com/stretchr/testify/mock"

	usermodel "github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// GetAvailabilitiesByUserFunc is an autogenerated mock type for the GetAvailabilitiesByUserFunc type
type GetAvailabilitiesByUserFunc struct {
	mock.Mock
}

type GetAvailabilitiesByUserFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *GetAvailabilitiesByUserFunc) EXPECT() *GetAvailabilitiesByUserFunc_Expecter {
	return &GetAvailabilitiesByUserFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, user
func (_m *GetAvailabilitiesByUserFunc) Execute(ctx context.Context, user usermodel.UserId) ([]catalogviews.AlbumSize, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 []catalogviews.AlbumSize
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usermodel.UserId) ([]catalogviews.AlbumSize, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usermodel.UserId) []catalogviews.AlbumSize); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]catalogviews.AlbumSize)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usermodel.UserId) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAvailabilitiesByUserFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type GetAvailabilitiesByUserFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - user usermodel.UserId
func (_e *GetAvailabilitiesByUserFunc_Expecter) Execute(ctx interface{}, user interface{}) *GetAvailabilitiesByUserFunc_Execute_Call {
	return &GetAvailabilitiesByUserFunc_Execute_Call{Call: _e.mock.On("Execute", ctx, user)}
}

func (_c *GetAvailabilitiesByUserFunc_Execute_Call) Run(run func(ctx context.Context, user usermodel.UserId)) *GetAvailabilitiesByUserFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(usermodel.UserId))
	})
	return _c
}

func (_c *GetAvailabilitiesByUserFunc_Execute_Call) Return(_a0 []catalogviews.AlbumSize, _a1 error) *GetAvailabilitiesByUserFunc_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GetAvailabilitiesByUserFunc_Execute_Call) RunAndReturn(run func(context.Context, usermodel.UserId) ([]catalogviews.AlbumSize, error)) *GetAvailabilitiesByUserFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewGetAvailabilitiesByUserFunc creates a new instance of GetAvailabilitiesByUserFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGetAvailabilitiesByUserFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *GetAvailabilitiesByUserFunc {
	mock := &GetAvailabilitiesByUserFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}