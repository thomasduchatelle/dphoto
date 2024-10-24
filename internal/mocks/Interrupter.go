// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	backup "github.com/thomasduchatelle/dphoto/pkg/backup"

	mock "github.com/stretchr/testify/mock"
)

// Interrupter is an autogenerated mock type for the Interrupter type
type Interrupter struct {
	mock.Mock
}

type Interrupter_Expecter struct {
	mock *mock.Mock
}

func (_m *Interrupter) EXPECT() *Interrupter_Expecter {
	return &Interrupter_Expecter{mock: &_m.Mock}
}

// Cancel provides a mock function with given fields:
func (_m *Interrupter) Cancel() {
	_m.Called()
}

// Interrupter_Cancel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Cancel'
type Interrupter_Cancel_Call struct {
	*mock.Call
}

// Cancel is a helper method to define mock.On call
func (_e *Interrupter_Expecter) Cancel() *Interrupter_Cancel_Call {
	return &Interrupter_Cancel_Call{Call: _e.mock.On("Cancel")}
}

func (_c *Interrupter_Cancel_Call) Run(run func()) *Interrupter_Cancel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Interrupter_Cancel_Call) Return() *Interrupter_Cancel_Call {
	_c.Call.Return()
	return _c
}

func (_c *Interrupter_Cancel_Call) RunAndReturn(run func()) *Interrupter_Cancel_Call {
	_c.Call.Return(run)
	return _c
}

// OnRejectedMedia provides a mock function with given fields: ctx, found, cause
func (_m *Interrupter) OnRejectedMedia(ctx context.Context, found backup.FoundMedia, cause error) {
	_m.Called(ctx, found, cause)
}

// Interrupter_OnRejectedMedia_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnRejectedMedia'
type Interrupter_OnRejectedMedia_Call struct {
	*mock.Call
}

// OnRejectedMedia is a helper method to define mock.On call
//   - ctx context.Context
//   - found backup.FoundMedia
//   - cause error
func (_e *Interrupter_Expecter) OnRejectedMedia(ctx interface{}, found interface{}, cause interface{}) *Interrupter_OnRejectedMedia_Call {
	return &Interrupter_OnRejectedMedia_Call{Call: _e.mock.On("OnRejectedMedia", ctx, found, cause)}
}

func (_c *Interrupter_OnRejectedMedia_Call) Run(run func(ctx context.Context, found backup.FoundMedia, cause error)) *Interrupter_OnRejectedMedia_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(backup.FoundMedia), args[2].(error))
	})
	return _c
}

func (_c *Interrupter_OnRejectedMedia_Call) Return() *Interrupter_OnRejectedMedia_Call {
	_c.Call.Return()
	return _c
}

func (_c *Interrupter_OnRejectedMedia_Call) RunAndReturn(run func(context.Context, backup.FoundMedia, error)) *Interrupter_OnRejectedMedia_Call {
	_c.Call.Return(run)
	return _c
}

// NewInterrupter creates a new instance of Interrupter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInterrupter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Interrupter {
	mock := &Interrupter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
