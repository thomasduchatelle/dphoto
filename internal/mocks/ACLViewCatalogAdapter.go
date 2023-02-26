// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// ACLViewCatalogAdapter is an autogenerated mock type for the ACLViewCatalogAdapter type
type ACLViewCatalogAdapter struct {
	mock.Mock
}

// FindAlbums provides a mock function with given fields: keys
func (_m *ACLViewCatalogAdapter) FindAlbums(keys []catalog.AlbumId) ([]*catalog.Album, error) {
	ret := _m.Called(keys)

	var r0 []*catalog.Album
	var r1 error
	if rf, ok := ret.Get(0).(func([]catalog.AlbumId) ([]*catalog.Album, error)); ok {
		return rf(keys)
	}
	if rf, ok := ret.Get(0).(func([]catalog.AlbumId) []*catalog.Album); ok {
		r0 = rf(keys)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.Album)
		}
	}

	if rf, ok := ret.Get(1).(func([]catalog.AlbumId) error); ok {
		r1 = rf(keys)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllAlbums provides a mock function with given fields: owner
func (_m *ACLViewCatalogAdapter) FindAllAlbums(owner string) ([]*catalog.Album, error) {
	ret := _m.Called(owner)

	var r0 []*catalog.Album
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]*catalog.Album, error)); ok {
		return rf(owner)
	}
	if rf, ok := ret.Get(0).(func(string) []*catalog.Album); ok {
		r0 = rf(owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.Album)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListMedias provides a mock function with given fields: owner, folderName, request
func (_m *ACLViewCatalogAdapter) ListMedias(owner string, folderName string, request catalog.PageRequest) (*catalog.MediaPage, error) {
	ret := _m.Called(owner, folderName, request)

	var r0 *catalog.MediaPage
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, catalog.PageRequest) (*catalog.MediaPage, error)); ok {
		return rf(owner, folderName, request)
	}
	if rf, ok := ret.Get(0).(func(string, string, catalog.PageRequest) *catalog.MediaPage); ok {
		r0 = rf(owner, folderName, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*catalog.MediaPage)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, catalog.PageRequest) error); ok {
		r1 = rf(owner, folderName, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewACLViewCatalogAdapter interface {
	mock.TestingT
	Cleanup(func())
}

// NewACLViewCatalogAdapter creates a new instance of ACLViewCatalogAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewACLViewCatalogAdapter(t mockConstructorTestingTNewACLViewCatalogAdapter) *ACLViewCatalogAdapter {
	mock := &ACLViewCatalogAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
