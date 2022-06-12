// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	catalog "github.com/thomasduchatelle/dphoto/domain/catalog"

	testing "testing"
)

// MoveMediaOperator is an autogenerated mock type for the MoveMediaOperator type
type MoveMediaOperator struct {
	mock.Mock
}

// Continue provides a mock function with given fields:
func (_m *MoveMediaOperator) Continue() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Move provides a mock function with given fields: source, dest
func (_m *MoveMediaOperator) Move(source catalog.MediaLocation, dest catalog.MediaLocation) (string, error) {
	ret := _m.Called(source, dest)

	var r0 string
	if rf, ok := ret.Get(0).(func(catalog.MediaLocation, catalog.MediaLocation) string); ok {
		r0 = rf(source, dest)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(catalog.MediaLocation, catalog.MediaLocation) error); ok {
		r1 = rf(source, dest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStatus provides a mock function with given fields: done, total
func (_m *MoveMediaOperator) UpdateStatus(done int, total int) error {
	ret := _m.Called(done, total)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, int) error); ok {
		r0 = rf(done, total)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMoveMediaOperator creates a new instance of MoveMediaOperator. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewMoveMediaOperator(t testing.TB) *MoveMediaOperator {
	mock := &MoveMediaOperator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
