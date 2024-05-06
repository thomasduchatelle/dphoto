// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

// IAccessTokenGenerator is an autogenerated mock type for the IAccessTokenGenerator type
type IAccessTokenGenerator struct {
	mock.Mock
}

type IAccessTokenGenerator_Expecter struct {
	mock *mock.Mock
}

func (_m *IAccessTokenGenerator) EXPECT() *IAccessTokenGenerator_Expecter {
	return &IAccessTokenGenerator_Expecter{mock: &_m.Mock}
}

// GenerateAccessToken provides a mock function with given fields: email
func (_m *IAccessTokenGenerator) GenerateAccessToken(email string) (*aclcore.Authentication, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for GenerateAccessToken")
	}

	var r0 *aclcore.Authentication
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*aclcore.Authentication, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *aclcore.Authentication); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*aclcore.Authentication)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IAccessTokenGenerator_GenerateAccessToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateAccessToken'
type IAccessTokenGenerator_GenerateAccessToken_Call struct {
	*mock.Call
}

// GenerateAccessToken is a helper method to define mock.On call
//   - email string
func (_e *IAccessTokenGenerator_Expecter) GenerateAccessToken(email interface{}) *IAccessTokenGenerator_GenerateAccessToken_Call {
	return &IAccessTokenGenerator_GenerateAccessToken_Call{Call: _e.mock.On("GenerateAccessToken", email)}
}

func (_c *IAccessTokenGenerator_GenerateAccessToken_Call) Run(run func(email string)) *IAccessTokenGenerator_GenerateAccessToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *IAccessTokenGenerator_GenerateAccessToken_Call) Return(_a0 *aclcore.Authentication, _a1 error) *IAccessTokenGenerator_GenerateAccessToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IAccessTokenGenerator_GenerateAccessToken_Call) RunAndReturn(run func(string) (*aclcore.Authentication, error)) *IAccessTokenGenerator_GenerateAccessToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewIAccessTokenGenerator creates a new instance of IAccessTokenGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIAccessTokenGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *IAccessTokenGenerator {
	mock := &IAccessTokenGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
