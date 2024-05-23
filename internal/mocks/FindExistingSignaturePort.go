// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"

	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// FindExistingSignaturePort is an autogenerated mock type for the FindExistingSignaturePort type
type FindExistingSignaturePort struct {
	mock.Mock
}

type FindExistingSignaturePort_Expecter struct {
	mock *mock.Mock
}

func (_m *FindExistingSignaturePort) EXPECT() *FindExistingSignaturePort_Expecter {
	return &FindExistingSignaturePort_Expecter{mock: &_m.Mock}
}

// FindSignatures provides a mock function with given fields: ctx, owner, signatures
func (_m *FindExistingSignaturePort) FindSignatures(ctx context.Context, owner ownermodel.Owner, signatures []catalog.MediaSignature) (map[catalog.MediaSignature]catalog.MediaId, error) {
	ret := _m.Called(ctx, owner, signatures)

	if len(ret) == 0 {
		panic("no return value specified for FindSignatures")
	}

	var r0 map[catalog.MediaSignature]catalog.MediaId
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner, []catalog.MediaSignature) (map[catalog.MediaSignature]catalog.MediaId, error)); ok {
		return rf(ctx, owner, signatures)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner, []catalog.MediaSignature) map[catalog.MediaSignature]catalog.MediaId); ok {
		r0 = rf(ctx, owner, signatures)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[catalog.MediaSignature]catalog.MediaId)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ownermodel.Owner, []catalog.MediaSignature) error); ok {
		r1 = rf(ctx, owner, signatures)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindExistingSignaturePort_FindSignatures_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindSignatures'
type FindExistingSignaturePort_FindSignatures_Call struct {
	*mock.Call
}

// FindSignatures is a helper method to define mock.On call
//   - ctx context.Context
//   - owner ownermodel.Owner
//   - signatures []catalog.MediaSignature
func (_e *FindExistingSignaturePort_Expecter) FindSignatures(ctx interface{}, owner interface{}, signatures interface{}) *FindExistingSignaturePort_FindSignatures_Call {
	return &FindExistingSignaturePort_FindSignatures_Call{Call: _e.mock.On("FindSignatures", ctx, owner, signatures)}
}

func (_c *FindExistingSignaturePort_FindSignatures_Call) Run(run func(ctx context.Context, owner ownermodel.Owner, signatures []catalog.MediaSignature)) *FindExistingSignaturePort_FindSignatures_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ownermodel.Owner), args[2].([]catalog.MediaSignature))
	})
	return _c
}

func (_c *FindExistingSignaturePort_FindSignatures_Call) Return(_a0 map[catalog.MediaSignature]catalog.MediaId, _a1 error) *FindExistingSignaturePort_FindSignatures_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FindExistingSignaturePort_FindSignatures_Call) RunAndReturn(run func(context.Context, ownermodel.Owner, []catalog.MediaSignature) (map[catalog.MediaSignature]catalog.MediaId, error)) *FindExistingSignaturePort_FindSignatures_Call {
	_c.Call.Return(run)
	return _c
}

// NewFindExistingSignaturePort creates a new instance of FindExistingSignaturePort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFindExistingSignaturePort(t interface {
	mock.TestingT
	Cleanup(func())
}) *FindExistingSignaturePort {
	mock := &FindExistingSignaturePort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}