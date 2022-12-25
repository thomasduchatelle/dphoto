// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
)

// Listener is an autogenerated mock type for the Listener type
type Listener struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *Listener) Execute(_a0 config.Config) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewListener interface {
	mock.TestingT
	Cleanup(func())
}

// NewListener creates a new instance of Listener. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewListener(t mockConstructorTestingTNewListener) *Listener {
	mock := &Listener{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
