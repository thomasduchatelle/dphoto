// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// InsertAlbumPortFunc is an autogenerated mock type for the InsertAlbumPortFunc type
type InsertAlbumPortFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, album
func (_m *InsertAlbumPortFunc) Execute(ctx context.Context, album catalog.Album) error {
	ret := _m.Called(ctx, album)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.Album) error); ok {
		r0 = rf(ctx, album)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewInsertAlbumPortFunc creates a new instance of InsertAlbumPortFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInsertAlbumPortFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *InsertAlbumPortFunc {
	mock := &InsertAlbumPortFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}