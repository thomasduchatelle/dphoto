// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"

	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// InsertMediasRepositoryPort is an autogenerated mock type for the InsertMediasRepositoryPort type
type InsertMediasRepositoryPort struct {
	mock.Mock
}

type InsertMediasRepositoryPort_Expecter struct {
	mock *mock.Mock
}

func (_m *InsertMediasRepositoryPort) EXPECT() *InsertMediasRepositoryPort_Expecter {
	return &InsertMediasRepositoryPort_Expecter{mock: &_m.Mock}
}

// InsertMedias provides a mock function with given fields: ctx, owner, media
func (_m *InsertMediasRepositoryPort) InsertMedias(ctx context.Context, owner ownermodel.Owner, media []catalog.CreateMediaRequest) error {
	ret := _m.Called(ctx, owner, media)

	if len(ret) == 0 {
		panic("no return value specified for InsertMedias")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner, []catalog.CreateMediaRequest) error); ok {
		r0 = rf(ctx, owner, media)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertMediasRepositoryPort_InsertMedias_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertMedias'
type InsertMediasRepositoryPort_InsertMedias_Call struct {
	*mock.Call
}

// InsertMedias is a helper method to define mock.On call
//   - ctx context.Context
//   - owner ownermodel.Owner
//   - media []catalog.CreateMediaRequest
func (_e *InsertMediasRepositoryPort_Expecter) InsertMedias(ctx interface{}, owner interface{}, media interface{}) *InsertMediasRepositoryPort_InsertMedias_Call {
	return &InsertMediasRepositoryPort_InsertMedias_Call{Call: _e.mock.On("InsertMedias", ctx, owner, media)}
}

func (_c *InsertMediasRepositoryPort_InsertMedias_Call) Run(run func(ctx context.Context, owner ownermodel.Owner, media []catalog.CreateMediaRequest)) *InsertMediasRepositoryPort_InsertMedias_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ownermodel.Owner), args[2].([]catalog.CreateMediaRequest))
	})
	return _c
}

func (_c *InsertMediasRepositoryPort_InsertMedias_Call) Return(_a0 error) *InsertMediasRepositoryPort_InsertMedias_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InsertMediasRepositoryPort_InsertMedias_Call) RunAndReturn(run func(context.Context, ownermodel.Owner, []catalog.CreateMediaRequest) error) *InsertMediasRepositoryPort_InsertMedias_Call {
	_c.Call.Return(run)
	return _c
}

// NewInsertMediasRepositoryPort creates a new instance of InsertMediasRepositoryPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInsertMediasRepositoryPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *InsertMediasRepositoryPort {
	mock := &InsertMediasRepositoryPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
