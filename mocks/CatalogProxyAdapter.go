// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	catalog "github.com/thomasduchatelle/dphoto/domain/catalog"

	testing "testing"
)

// CatalogProxyAdapter is an autogenerated mock type for the CatalogProxyAdapter type
type CatalogProxyAdapter struct {
	mock.Mock
}

// Create provides a mock function with given fields: createRequest
func (_m *CatalogProxyAdapter) Create(createRequest catalog.CreateAlbum) error {
	ret := _m.Called(createRequest)

	var r0 error
	if rf, ok := ret.Get(0).(func(catalog.CreateAlbum) error); ok {
		r0 = rf(createRequest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllAlbums provides a mock function with given fields: owner
func (_m *CatalogProxyAdapter) FindAllAlbums(owner string) ([]*catalog.Album, error) {
	ret := _m.Called(owner)

	var r0 []*catalog.Album
	if rf, ok := ret.Get(0).(func(string) []*catalog.Album); ok {
		r0 = rf(owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.Album)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindSignatures provides a mock function with given fields: owner, signatures
func (_m *CatalogProxyAdapter) FindSignatures(owner string, signatures []*catalog.MediaSignature) ([]*catalog.MediaSignature, error) {
	ret := _m.Called(owner, signatures)

	var r0 []*catalog.MediaSignature
	if rf, ok := ret.Get(0).(func(string, []*catalog.MediaSignature) []*catalog.MediaSignature); ok {
		r0 = rf(owner, signatures)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.MediaSignature)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []*catalog.MediaSignature) error); ok {
		r1 = rf(owner, signatures)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertMedias provides a mock function with given fields: owner, medias
func (_m *CatalogProxyAdapter) InsertMedias(owner string, medias []catalog.CreateMediaRequest) error {
	ret := _m.Called(owner, medias)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []catalog.CreateMediaRequest) error); ok {
		r0 = rf(owner, medias)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCatalogProxyAdapter creates a new instance of CatalogProxyAdapter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewCatalogProxyAdapter(t testing.TB) *CatalogProxyAdapter {
	mock := &CatalogProxyAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
