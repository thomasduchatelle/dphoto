// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"

	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// InsertMediaSimulator is an autogenerated mock type for the InsertMediaSimulator type
type InsertMediaSimulator struct {
	mock.Mock
}

type InsertMediaSimulator_Expecter struct {
	mock *mock.Mock
}

func (_m *InsertMediaSimulator) EXPECT() *InsertMediaSimulator_Expecter {
	return &InsertMediaSimulator_Expecter{mock: &_m.Mock}
}

// SimulateInsertingMedia provides a mock function with given fields: ctx, owner, signatures
func (_m *InsertMediaSimulator) SimulateInsertingMedia(ctx context.Context, owner ownermodel.Owner, signatures []catalog.MediaSignature) ([]catalog.MediaFutureReference, error) {
	ret := _m.Called(ctx, owner, signatures)

	if len(ret) == 0 {
		panic("no return value specified for SimulateInsertingMedia")
	}

	var r0 []catalog.MediaFutureReference
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner, []catalog.MediaSignature) ([]catalog.MediaFutureReference, error)); ok {
		return rf(ctx, owner, signatures)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner, []catalog.MediaSignature) []catalog.MediaFutureReference); ok {
		r0 = rf(ctx, owner, signatures)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]catalog.MediaFutureReference)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ownermodel.Owner, []catalog.MediaSignature) error); ok {
		r1 = rf(ctx, owner, signatures)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertMediaSimulator_SimulateInsertingMedia_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SimulateInsertingMedia'
type InsertMediaSimulator_SimulateInsertingMedia_Call struct {
	*mock.Call
}

// SimulateInsertingMedia is a helper method to define mock.On call
//   - ctx context.Context
//   - owner ownermodel.Owner
//   - signatures []catalog.MediaSignature
func (_e *InsertMediaSimulator_Expecter) SimulateInsertingMedia(ctx interface{}, owner interface{}, signatures interface{}) *InsertMediaSimulator_SimulateInsertingMedia_Call {
	return &InsertMediaSimulator_SimulateInsertingMedia_Call{Call: _e.mock.On("SimulateInsertingMedia", ctx, owner, signatures)}
}

func (_c *InsertMediaSimulator_SimulateInsertingMedia_Call) Run(run func(ctx context.Context, owner ownermodel.Owner, signatures []catalog.MediaSignature)) *InsertMediaSimulator_SimulateInsertingMedia_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ownermodel.Owner), args[2].([]catalog.MediaSignature))
	})
	return _c
}

func (_c *InsertMediaSimulator_SimulateInsertingMedia_Call) Return(_a0 []catalog.MediaFutureReference, _a1 error) *InsertMediaSimulator_SimulateInsertingMedia_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *InsertMediaSimulator_SimulateInsertingMedia_Call) RunAndReturn(run func(context.Context, ownermodel.Owner, []catalog.MediaSignature) ([]catalog.MediaFutureReference, error)) *InsertMediaSimulator_SimulateInsertingMedia_Call {
	_c.Call.Return(run)
	return _c
}

// NewInsertMediaSimulator creates a new instance of InsertMediaSimulator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInsertMediaSimulator(t interface {
	mock.TestingT
	Cleanup(func())
}) *InsertMediaSimulator {
	mock := &InsertMediaSimulator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
