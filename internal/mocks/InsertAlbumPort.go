// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// InsertAlbumPort is an autogenerated mock type for the InsertAlbumPort type
type InsertAlbumPort struct {
	mock.Mock
}

// InsertAlbum provides a mock function with given fields: ctx, album
func (_m *InsertAlbumPort) InsertAlbum(ctx context.Context, album catalog.Album) error {
	ret := _m.Called(ctx, album)

	if len(ret) == 0 {
		panic("no return value specified for InsertAlbum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.Album) error); ok {
		r0 = rf(ctx, album)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewInsertAlbumPort creates a new instance of InsertAlbumPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInsertAlbumPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *InsertAlbumPort {
	mock := &InsertAlbumPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
