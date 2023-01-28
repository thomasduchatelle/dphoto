// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// ShareAlbumCatalogPort is an autogenerated mock type for the ShareAlbumCatalogPort type
type ShareAlbumCatalogPort struct {
	mock.Mock
}

// FindAlbum provides a mock function with given fields: owner, folderName
func (_m *ShareAlbumCatalogPort) FindAlbum(owner string, folderName string) (*catalog.Album, error) {
	ret := _m.Called(owner, folderName)

	var r0 *catalog.Album
	if rf, ok := ret.Get(0).(func(string, string) *catalog.Album); ok {
		r0 = rf(owner, folderName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*catalog.Album)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(owner, folderName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewShareAlbumCatalogPort interface {
	mock.TestingT
	Cleanup(func())
}

// NewShareAlbumCatalogPort creates a new instance of ShareAlbumCatalogPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewShareAlbumCatalogPort(t mockConstructorTestingTNewShareAlbumCatalogPort) *ShareAlbumCatalogPort {
	mock := &ShareAlbumCatalogPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}