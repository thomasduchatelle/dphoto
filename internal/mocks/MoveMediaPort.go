// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// MoveMediaPort is an autogenerated mock type for the MoveMediaPort type
type MoveMediaPort struct {
	mock.Mock
}

// MoveMedia provides a mock function with given fields: ctx, albumId, mediaIds
func (_m *MoveMediaPort) MoveMedia(ctx context.Context, albumId catalog.AlbumId, mediaIds []catalog.MediaId) error {
	ret := _m.Called(ctx, albumId, mediaIds)

	if len(ret) == 0 {
		panic("no return value specified for MoveMedia")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.AlbumId, []catalog.MediaId) error); ok {
		r0 = rf(ctx, albumId, mediaIds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMoveMediaPort creates a new instance of MoveMediaPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMoveMediaPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *MoveMediaPort {
	mock := &MoveMediaPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
