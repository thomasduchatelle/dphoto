// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	aws "github.com/aws/aws-sdk-go-v2/aws"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ConfigFactoryFunc is an autogenerated mock type for the ConfigFactoryFunc type
type ConfigFactoryFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx
func (_m *ConfigFactoryFunc) Execute(ctx context.Context) (aws.Config, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 aws.Config
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (aws.Config, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) aws.Config); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(aws.Config)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewConfigFactoryFunc creates a new instance of ConfigFactoryFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConfigFactoryFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConfigFactoryFunc {
	mock := &ConfigFactoryFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
