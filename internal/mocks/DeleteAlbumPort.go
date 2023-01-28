// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DeleteAlbumPort is an autogenerated mock type for the DeleteAlbumPort type
type DeleteAlbumPort struct {
	mock.Mock
}

// DeleteAlbum provides a mock function with given fields: folderName
func (_m *DeleteAlbumPort) DeleteAlbum(folderName string) error {
	ret := _m.Called(folderName)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(folderName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewDeleteAlbumPort interface {
	mock.TestingT
	Cleanup(func())
}

// NewDeleteAlbumPort creates a new instance of DeleteAlbumPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDeleteAlbumPort(t mockConstructorTestingTNewDeleteAlbumPort) *DeleteAlbumPort {
	mock := &DeleteAlbumPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}