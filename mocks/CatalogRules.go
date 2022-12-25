// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// CatalogRules is an autogenerated mock type for the CatalogRules type
type CatalogRules struct {
	mock.Mock
}

// CanListMediasFromAlbum provides a mock function with given fields: owner, folderName
func (_m *CatalogRules) CanListMediasFromAlbum(owner string, folderName string) error {
	ret := _m.Called(owner, folderName)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(owner, folderName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CanReadMedia provides a mock function with given fields: owner, id
func (_m *CatalogRules) CanReadMedia(owner string, id string) error {
	ret := _m.Called(owner, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(owner, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Owner provides a mock function with given fields:
func (_m *CatalogRules) Owner() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SharedByUserGrid provides a mock function with given fields: owner
func (_m *CatalogRules) SharedByUserGrid(owner string) (map[string][]string, error) {
	ret := _m.Called(owner)

	var r0 map[string][]string
	if rf, ok := ret.Get(0).(func(string) map[string][]string); ok {
		r0 = rf(owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]string)
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

// SharedWithUserAlbum provides a mock function with given fields:
func (_m *CatalogRules) SharedWithUserAlbum() ([]catalog.AlbumId, error) {
	ret := _m.Called()

	var r0 []catalog.AlbumId
	if rf, ok := ret.Get(0).(func() []catalog.AlbumId); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]catalog.AlbumId)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCatalogRules interface {
	mock.TestingT
	Cleanup(func())
}

// NewCatalogRules creates a new instance of CatalogRules. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCatalogRules(t mockConstructorTestingTNewCatalogRules) *CatalogRules {
	mock := &CatalogRules{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
