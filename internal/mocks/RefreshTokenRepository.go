// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

// RefreshTokenRepository is an autogenerated mock type for the RefreshTokenRepository type
type RefreshTokenRepository struct {
	mock.Mock
}

// DeleteRefreshToken provides a mock function with given fields: token
func (_m *RefreshTokenRepository) DeleteRefreshToken(token string) error {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for DeleteRefreshToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindRefreshToken provides a mock function with given fields: token
func (_m *RefreshTokenRepository) FindRefreshToken(token string) (*aclcore.RefreshTokenSpec, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for FindRefreshToken")
	}

	var r0 *aclcore.RefreshTokenSpec
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*aclcore.RefreshTokenSpec, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(string) *aclcore.RefreshTokenSpec); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*aclcore.RefreshTokenSpec)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HouseKeepRefreshToken provides a mock function with given fields:
func (_m *RefreshTokenRepository) HouseKeepRefreshToken() (int, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for HouseKeepRefreshToken")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func() (int, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreRefreshToken provides a mock function with given fields: token, spec
func (_m *RefreshTokenRepository) StoreRefreshToken(token string, spec aclcore.RefreshTokenSpec) error {
	ret := _m.Called(token, spec)

	if len(ret) == 0 {
		panic("no return value specified for StoreRefreshToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, aclcore.RefreshTokenSpec) error); ok {
		r0 = rf(token, spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRefreshTokenRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *RefreshTokenRepository {
	mock := &RefreshTokenRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
