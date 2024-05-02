// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// TimelineMutationObserver is an autogenerated mock type for the TimelineMutationObserver type
type TimelineMutationObserver struct {
	mock.Mock
}

// Observe provides a mock function with given fields: ctx, transfers
func (_m *TimelineMutationObserver) Observe(ctx context.Context, transfers catalog.TransferredMedias) error {
	ret := _m.Called(ctx, transfers)

	if len(ret) == 0 {
		panic("no return value specified for Observe")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.TransferredMedias) error); ok {
		r0 = rf(ctx, transfers)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTimelineMutationObserver creates a new instance of TimelineMutationObserver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTimelineMutationObserver(t interface {
	mock.TestingT
	Cleanup(func())
}) *TimelineMutationObserver {
	mock := &TimelineMutationObserver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}