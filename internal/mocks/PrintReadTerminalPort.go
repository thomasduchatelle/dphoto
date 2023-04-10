// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PrintReadTerminalPort is an autogenerated mock type for the PrintReadTerminalPort type
type PrintReadTerminalPort struct {
	mock.Mock
}

// Print provides a mock function with given fields: question
func (_m *PrintReadTerminalPort) Print(question string) {
	_m.Called(question)
}

// ReadAnswer provides a mock function with given fields:
func (_m *PrintReadTerminalPort) ReadAnswer() (string, error) {
	ret := _m.Called()

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPrintReadTerminalPort interface {
	mock.TestingT
	Cleanup(func())
}

// NewPrintReadTerminalPort creates a new instance of PrintReadTerminalPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPrintReadTerminalPort(t mockConstructorTestingTNewPrintReadTerminalPort) *PrintReadTerminalPort {
	mock := &PrintReadTerminalPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
