// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// CountMediasBySelectorsFunc is an autogenerated mock type for the CountMediasBySelectorsFunc type
type CountMediasBySelectorsFunc struct {
	mock.Mock
}

type CountMediasBySelectorsFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *CountMediasBySelectorsFunc) EXPECT() *CountMediasBySelectorsFunc_Expecter {
	return &CountMediasBySelectorsFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, owner, selectors
func (_m *CountMediasBySelectorsFunc) Execute(ctx context.Context, owner catalog.Owner, selectors []catalog.MediaSelector) (int, error) {
	ret := _m.Called(ctx, owner, selectors)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.Owner, []catalog.MediaSelector) (int, error)); ok {
		return rf(ctx, owner, selectors)
	}
	if rf, ok := ret.Get(0).(func(context.Context, catalog.Owner, []catalog.MediaSelector) int); ok {
		r0 = rf(ctx, owner, selectors)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, catalog.Owner, []catalog.MediaSelector) error); ok {
		r1 = rf(ctx, owner, selectors)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountMediasBySelectorsFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type CountMediasBySelectorsFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - owner catalog.Owner
//   - selectors []catalog.MediaSelector
func (_e *CountMediasBySelectorsFunc_Expecter) Execute(ctx interface{}, owner interface{}, selectors interface{}) *CountMediasBySelectorsFunc_Execute_Call {
	return &CountMediasBySelectorsFunc_Execute_Call{Call: _e.mock.On("Execute", ctx, owner, selectors)}
}

func (_c *CountMediasBySelectorsFunc_Execute_Call) Run(run func(ctx context.Context, owner catalog.Owner, selectors []catalog.MediaSelector)) *CountMediasBySelectorsFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(catalog.Owner), args[2].([]catalog.MediaSelector))
	})
	return _c
}

func (_c *CountMediasBySelectorsFunc_Execute_Call) Return(_a0 int, _a1 error) *CountMediasBySelectorsFunc_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CountMediasBySelectorsFunc_Execute_Call) RunAndReturn(run func(context.Context, catalog.Owner, []catalog.MediaSelector) (int, error)) *CountMediasBySelectorsFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewCountMediasBySelectorsFunc creates a new instance of CountMediasBySelectorsFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCountMediasBySelectorsFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *CountMediasBySelectorsFunc {
	mock := &CountMediasBySelectorsFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
