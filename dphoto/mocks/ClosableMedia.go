// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ClosableMedia is an autogenerated mock type for the ClosableMedia type
type ClosableMedia struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *ClosableMedia) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}