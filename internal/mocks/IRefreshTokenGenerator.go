// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

// IRefreshTokenGenerator is an autogenerated mock type for the IRefreshTokenGenerator type
type IRefreshTokenGenerator struct {
	mock.Mock
}

type IRefreshTokenGenerator_Expecter struct {
	mock *mock.Mock
}

func (_m *IRefreshTokenGenerator) EXPECT() *IRefreshTokenGenerator_Expecter {
	return &IRefreshTokenGenerator_Expecter{mock: &_m.Mock}
}

// GenerateRefreshToken provides a mock function with given fields: spec
func (_m *IRefreshTokenGenerator) GenerateRefreshToken(spec aclcore.RefreshTokenSpec) (string, error) {
	ret := _m.Called(spec)

	if len(ret) == 0 {
		panic("no return value specified for GenerateRefreshToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(aclcore.RefreshTokenSpec) (string, error)); ok {
		return rf(spec)
	}
	if rf, ok := ret.Get(0).(func(aclcore.RefreshTokenSpec) string); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(aclcore.RefreshTokenSpec) error); ok {
		r1 = rf(spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IRefreshTokenGenerator_GenerateRefreshToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateRefreshToken'
type IRefreshTokenGenerator_GenerateRefreshToken_Call struct {
	*mock.Call
}

// GenerateRefreshToken is a helper method to define mock.On call
//   - spec aclcore.RefreshTokenSpec
func (_e *IRefreshTokenGenerator_Expecter) GenerateRefreshToken(spec interface{}) *IRefreshTokenGenerator_GenerateRefreshToken_Call {
	return &IRefreshTokenGenerator_GenerateRefreshToken_Call{Call: _e.mock.On("GenerateRefreshToken", spec)}
}

func (_c *IRefreshTokenGenerator_GenerateRefreshToken_Call) Run(run func(spec aclcore.RefreshTokenSpec)) *IRefreshTokenGenerator_GenerateRefreshToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(aclcore.RefreshTokenSpec))
	})
	return _c
}

func (_c *IRefreshTokenGenerator_GenerateRefreshToken_Call) Return(_a0 string, _a1 error) *IRefreshTokenGenerator_GenerateRefreshToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IRefreshTokenGenerator_GenerateRefreshToken_Call) RunAndReturn(run func(aclcore.RefreshTokenSpec) (string, error)) *IRefreshTokenGenerator_GenerateRefreshToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewIRefreshTokenGenerator creates a new instance of IRefreshTokenGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIRefreshTokenGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *IRefreshTokenGenerator {
	mock := &IRefreshTokenGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
