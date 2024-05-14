// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// MediaTransferFunc is an autogenerated mock type for the MediaTransferFunc type
type MediaTransferFunc struct {
	mock.Mock
}

type MediaTransferFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *MediaTransferFunc) EXPECT() *MediaTransferFunc_Expecter {
	return &MediaTransferFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, records
func (_m *MediaTransferFunc) Execute(ctx context.Context, records catalog.MediaTransferRecords) error {
	ret := _m.Called(ctx, records)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.MediaTransferRecords) error); ok {
		r0 = rf(ctx, records)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MediaTransferFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MediaTransferFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - records catalog.MediaTransferRecords
func (_e *MediaTransferFunc_Expecter) Execute(ctx interface{}, records interface{}) *MediaTransferFunc_Execute_Call {
	return &MediaTransferFunc_Execute_Call{Call: _e.mock.On("Execute", ctx, records)}
}

func (_c *MediaTransferFunc_Execute_Call) Run(run func(ctx context.Context, records catalog.MediaTransferRecords)) *MediaTransferFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(catalog.MediaTransferRecords))
	})
	return _c
}

func (_c *MediaTransferFunc_Execute_Call) Return(_a0 error) *MediaTransferFunc_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MediaTransferFunc_Execute_Call) RunAndReturn(run func(context.Context, catalog.MediaTransferRecords) error) *MediaTransferFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMediaTransferFunc creates a new instance of MediaTransferFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMediaTransferFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *MediaTransferFunc {
	mock := &MediaTransferFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}