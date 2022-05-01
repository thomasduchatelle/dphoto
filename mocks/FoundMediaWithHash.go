// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// FoundMediaWithHash is an autogenerated mock type for the FoundMediaWithHash type
type FoundMediaWithHash struct {
	mock.Mock
}

// Sha256Hash provides a mock function with given fields:
func (_m *FoundMediaWithHash) Sha256Hash() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewFoundMediaWithHash creates a new instance of FoundMediaWithHash. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewFoundMediaWithHash(t testing.TB) *FoundMediaWithHash {
	mock := &FoundMediaWithHash{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
