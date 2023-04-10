// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UserInputPort is an autogenerated mock type for the UserInputPort type
type UserInputPort struct {
	mock.Mock
}

// StartListening provides a mock function with given fields:
func (_m *UserInputPort) StartListening() {
	_m.Called()
}

type mockConstructorTestingTNewUserInputPort interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserInputPort creates a new instance of UserInputPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserInputPort(t mockConstructorTestingTNewUserInputPort) *UserInputPort {
	mock := &UserInputPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
