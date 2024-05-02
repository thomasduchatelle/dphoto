// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// MoveMediaPortFunc is an autogenerated mock type for the MoveMediaPortFunc type
type MoveMediaPortFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, albumId, mediaIds
func (_m *MoveMediaPortFunc) Execute(ctx context.Context, albumId catalog.AlbumId, mediaIds []catalog.MediaId) error {
	ret := _m.Called(ctx, albumId, mediaIds)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.AlbumId, []catalog.MediaId) error); ok {
		r0 = rf(ctx, albumId, mediaIds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMoveMediaPortFunc creates a new instance of MoveMediaPortFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMoveMediaPortFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *MoveMediaPortFunc {
	mock := &MoveMediaPortFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}