// Code generated by mockery v2.43.0. DO NOT EDIT.

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

// NewUserInputPort creates a new instance of UserInputPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserInputPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserInputPort {
	mock := &UserInputPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
