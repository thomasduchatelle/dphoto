// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// RenameAlbumPort is an autogenerated mock type for the RenameAlbumPort type
type RenameAlbumPort struct {
	mock.Mock
}

// RenameAlbum provides a mock function with given fields: folderName, newName, renameFolder
func (_m *RenameAlbumPort) RenameAlbum(folderName string, newName string, renameFolder bool) error {
	ret := _m.Called(folderName, newName, renameFolder)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, bool) error); ok {
		r0 = rf(folderName, newName, renameFolder)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRenameAlbumPort interface {
	mock.TestingT
	Cleanup(func())
}

// NewRenameAlbumPort creates a new instance of RenameAlbumPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRenameAlbumPort(t mockConstructorTestingTNewRenameAlbumPort) *RenameAlbumPort {
	mock := &RenameAlbumPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
