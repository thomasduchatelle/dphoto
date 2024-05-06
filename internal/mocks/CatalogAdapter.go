// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// CatalogAdapter is an autogenerated mock type for the CatalogAdapter type
type CatalogAdapter struct {
	mock.Mock
}

// AssignIdsToNewMedias provides a mock function with given fields: owner, medias
func (_m *CatalogAdapter) AssignIdsToNewMedias(owner string, medias []*backup.AnalysedMedia) (map[*backup.AnalysedMedia]string, error) {
	ret := _m.Called(owner, medias)

	if len(ret) == 0 {
		panic("no return value specified for AssignIdsToNewMedias")
	}

	var r0 map[*backup.AnalysedMedia]string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, []*backup.AnalysedMedia) (map[*backup.AnalysedMedia]string, error)); ok {
		return rf(owner, medias)
	}
	if rf, ok := ret.Get(0).(func(string, []*backup.AnalysedMedia) map[*backup.AnalysedMedia]string); ok {
		r0 = rf(owner, medias)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[*backup.AnalysedMedia]string)
		}
	}

	if rf, ok := ret.Get(1).(func(string, []*backup.AnalysedMedia) error); ok {
		r1 = rf(owner, medias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlbumsTimeline provides a mock function with given fields: owner
func (_m *CatalogAdapter) GetAlbumsTimeline(owner string) (backup.TimelineAdapter, error) {
	ret := _m.Called(owner)

	if len(ret) == 0 {
		panic("no return value specified for GetAlbumsTimeline")
	}

	var r0 backup.TimelineAdapter
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (backup.TimelineAdapter, error)); ok {
		return rf(owner)
	}
	if rf, ok := ret.Get(0).(func(string) backup.TimelineAdapter); ok {
		r0 = rf(owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(backup.TimelineAdapter)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IndexMedias provides a mock function with given fields: owner, requests
func (_m *CatalogAdapter) IndexMedias(owner string, requests []*backup.CatalogMediaRequest) error {
	ret := _m.Called(owner, requests)

	if len(ret) == 0 {
		panic("no return value specified for IndexMedias")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []*backup.CatalogMediaRequest) error); ok {
		r0 = rf(owner, requests)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCatalogAdapter creates a new instance of CatalogAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCatalogAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *CatalogAdapter {
	mock := &CatalogAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
