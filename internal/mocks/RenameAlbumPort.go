// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// RenameAlbumPort is an autogenerated mock type for the RenameAlbumPort type
type RenameAlbumPort struct {
	mock.Mock
}

// RenameAlbum provides a mock function with given fields: folderName, newName, renameFolder
func (_m *RenameAlbumPort) RenameAlbum(folderName string, newName string, renameFolder bool) error {
	ret := _m.Called(folderName, newName, renameFolder)

	if len(ret) == 0 {
		panic("no return value specified for RenameAlbum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, bool) error); ok {
		r0 = rf(folderName, newName, renameFolder)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRenameAlbumPort creates a new instance of RenameAlbumPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRenameAlbumPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *RenameAlbumPort {
	mock := &RenameAlbumPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
