// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalogviews "github.com/thomasduchatelle/dphoto/pkg/catalogviews"

	mock "github.com/stretchr/testify/mock"

	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"

	usermodel "github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// GetCurrentAlbumSizesPort is an autogenerated mock type for the GetCurrentAlbumSizesPort type
type GetCurrentAlbumSizesPort struct {
	mock.Mock
}

type GetCurrentAlbumSizesPort_Expecter struct {
	mock *mock.Mock
}

func (_m *GetCurrentAlbumSizesPort) EXPECT() *GetCurrentAlbumSizesPort_Expecter {
	return &GetCurrentAlbumSizesPort_Expecter{mock: &_m.Mock}
}

// GetAlbumSizes provides a mock function with given fields: ctx, userId, owner
func (_m *GetCurrentAlbumSizesPort) GetAlbumSizes(ctx context.Context, userId usermodel.UserId, owner ...ownermodel.Owner) ([]catalogviews.UserAlbumSize, error) {
	_va := make([]interface{}, len(owner))
	for _i := range owner {
		_va[_i] = owner[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, userId)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetAlbumSizes")
	}

	var r0 []catalogviews.UserAlbumSize
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usermodel.UserId, ...ownermodel.Owner) ([]catalogviews.UserAlbumSize, error)); ok {
		return rf(ctx, userId, owner...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usermodel.UserId, ...ownermodel.Owner) []catalogviews.UserAlbumSize); ok {
		r0 = rf(ctx, userId, owner...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]catalogviews.UserAlbumSize)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usermodel.UserId, ...ownermodel.Owner) error); ok {
		r1 = rf(ctx, userId, owner...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCurrentAlbumSizesPort_GetAlbumSizes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAlbumSizes'
type GetCurrentAlbumSizesPort_GetAlbumSizes_Call struct {
	*mock.Call
}

// GetAlbumSizes is a helper method to define mock.On call
//   - ctx context.Context
//   - userId usermodel.UserId
//   - owner ...ownermodel.Owner
func (_e *GetCurrentAlbumSizesPort_Expecter) GetAlbumSizes(ctx interface{}, userId interface{}, owner ...interface{}) *GetCurrentAlbumSizesPort_GetAlbumSizes_Call {
	return &GetCurrentAlbumSizesPort_GetAlbumSizes_Call{Call: _e.mock.On("GetAlbumSizes",
		append([]interface{}{ctx, userId}, owner...)...)}
}

func (_c *GetCurrentAlbumSizesPort_GetAlbumSizes_Call) Run(run func(ctx context.Context, userId usermodel.UserId, owner ...ownermodel.Owner)) *GetCurrentAlbumSizesPort_GetAlbumSizes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]ownermodel.Owner, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(ownermodel.Owner)
			}
		}
		run(args[0].(context.Context), args[1].(usermodel.UserId), variadicArgs...)
	})
	return _c
}

func (_c *GetCurrentAlbumSizesPort_GetAlbumSizes_Call) Return(_a0 []catalogviews.UserAlbumSize, _a1 error) *GetCurrentAlbumSizesPort_GetAlbumSizes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GetCurrentAlbumSizesPort_GetAlbumSizes_Call) RunAndReturn(run func(context.Context, usermodel.UserId, ...ownermodel.Owner) ([]catalogviews.UserAlbumSize, error)) *GetCurrentAlbumSizesPort_GetAlbumSizes_Call {
	_c.Call.Return(run)
	return _c
}

// NewGetCurrentAlbumSizesPort creates a new instance of GetCurrentAlbumSizesPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGetCurrentAlbumSizesPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *GetCurrentAlbumSizesPort {
	mock := &GetCurrentAlbumSizesPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}