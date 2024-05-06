// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	ui "github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
)

// InteractiveRendererPort is an autogenerated mock type for the InteractiveRendererPort type
type InteractiveRendererPort struct {
	mock.Mock
}

// Height provides a mock function with given fields:
func (_m *InteractiveRendererPort) Height() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Height")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Print provides a mock function with given fields: question
func (_m *InteractiveRendererPort) Print(question string) {
	_m.Called(question)
}

// ReadAnswer provides a mock function with given fields:
func (_m *InteractiveRendererPort) ReadAnswer() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadAnswer")
	}

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

// Render provides a mock function with given fields: state
func (_m *InteractiveRendererPort) Render(state *ui.InteractiveViewState) error {
	ret := _m.Called(state)

	if len(ret) == 0 {
		panic("no return value specified for Render")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*ui.InteractiveViewState) error); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TakeOverScreen provides a mock function with given fields:
func (_m *InteractiveRendererPort) TakeOverScreen() {
	_m.Called()
}

// NewInteractiveRendererPort creates a new instance of InteractiveRendererPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInteractiveRendererPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *InteractiveRendererPort {
	mock := &InteractiveRendererPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
