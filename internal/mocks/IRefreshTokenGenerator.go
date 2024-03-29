// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

// IRefreshTokenGenerator is an autogenerated mock type for the IRefreshTokenGenerator type
type IRefreshTokenGenerator struct {
	mock.Mock
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
