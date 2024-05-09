// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// TimelineMutationObserverFunc is an autogenerated mock type for the TimelineMutationObserverFunc type
type TimelineMutationObserverFunc struct {
	mock.Mock
}

type TimelineMutationObserverFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *TimelineMutationObserverFunc) EXPECT() *TimelineMutationObserverFunc_Expecter {
	return &TimelineMutationObserverFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, transfers
func (_m *TimelineMutationObserverFunc) Execute(ctx context.Context, transfers catalog.TransferredMedias) error {
	ret := _m.Called(ctx, transfers)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.TransferredMedias) error); ok {
		r0 = rf(ctx, transfers)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TimelineMutationObserverFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type TimelineMutationObserverFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - transfers catalog.TransferredMedias
func (_e *TimelineMutationObserverFunc_Expecter) Execute(ctx interface{}, transfers interface{}) *TimelineMutationObserverFunc_Execute_Call {
	return &TimelineMutationObserverFunc_Execute_Call{Call: _e.mock.On("Execute", ctx, transfers)}
}

func (_c *TimelineMutationObserverFunc_Execute_Call) Run(run func(ctx context.Context, transfers catalog.TransferredMedias)) *TimelineMutationObserverFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(catalog.TransferredMedias))
	})
	return _c
}

func (_c *TimelineMutationObserverFunc_Execute_Call) Return(_a0 error) *TimelineMutationObserverFunc_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TimelineMutationObserverFunc_Execute_Call) RunAndReturn(run func(context.Context, catalog.TransferredMedias) error) *TimelineMutationObserverFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewTimelineMutationObserverFunc creates a new instance of TimelineMutationObserverFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTimelineMutationObserverFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *TimelineMutationObserverFunc {
	mock := &TimelineMutationObserverFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
