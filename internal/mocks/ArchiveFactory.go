// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	archive "github.com/thomasduchatelle/dphoto/pkg/archive"

	mock "github.com/stretchr/testify/mock"
)

// ArchiveFactory is an autogenerated mock type for the ArchiveFactory type
type ArchiveFactory struct {
	mock.Mock
}

type ArchiveFactory_Expecter struct {
	mock *mock.Mock
}

func (_m *ArchiveFactory) EXPECT() *ArchiveFactory_Expecter {
	return &ArchiveFactory_Expecter{mock: &_m.Mock}
}

// ArchiveAsyncJobAdapter provides a mock function with given fields: ctx
func (_m *ArchiveFactory) ArchiveAsyncJobAdapter(ctx context.Context) archive.AsyncJobAdapter {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ArchiveAsyncJobAdapter")
	}

	var r0 archive.AsyncJobAdapter
	if rf, ok := ret.Get(0).(func(context.Context) archive.AsyncJobAdapter); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(archive.AsyncJobAdapter)
		}
	}

	return r0
}

// ArchiveFactory_ArchiveAsyncJobAdapter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ArchiveAsyncJobAdapter'
type ArchiveFactory_ArchiveAsyncJobAdapter_Call struct {
	*mock.Call
}

// ArchiveAsyncJobAdapter is a helper method to define mock.On call
//   - ctx context.Context
func (_e *ArchiveFactory_Expecter) ArchiveAsyncJobAdapter(ctx interface{}) *ArchiveFactory_ArchiveAsyncJobAdapter_Call {
	return &ArchiveFactory_ArchiveAsyncJobAdapter_Call{Call: _e.mock.On("ArchiveAsyncJobAdapter", ctx)}
}

func (_c *ArchiveFactory_ArchiveAsyncJobAdapter_Call) Run(run func(ctx context.Context)) *ArchiveFactory_ArchiveAsyncJobAdapter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *ArchiveFactory_ArchiveAsyncJobAdapter_Call) Return(_a0 archive.AsyncJobAdapter) *ArchiveFactory_ArchiveAsyncJobAdapter_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ArchiveFactory_ArchiveAsyncJobAdapter_Call) RunAndReturn(run func(context.Context) archive.AsyncJobAdapter) *ArchiveFactory_ArchiveAsyncJobAdapter_Call {
	_c.Call.Return(run)
	return _c
}

// NewArchiveFactory creates a new instance of ArchiveFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewArchiveFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *ArchiveFactory {
	mock := &ArchiveFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
