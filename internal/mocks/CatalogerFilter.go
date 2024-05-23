// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// CatalogerFilter is an autogenerated mock type for the CatalogerFilter type
type CatalogerFilter struct {
	mock.Mock
}

type CatalogerFilter_Expecter struct {
	mock *mock.Mock
}

func (_m *CatalogerFilter) EXPECT() *CatalogerFilter_Expecter {
	return &CatalogerFilter_Expecter{mock: &_m.Mock}
}

// FilterOut provides a mock function with given fields: media, reference
func (_m *CatalogerFilter) FilterOut(media backup.AnalysedMedia, reference backup.CatalogReference) (backup.ProgressEventType, bool) {
	ret := _m.Called(media, reference)

	if len(ret) == 0 {
		panic("no return value specified for FilterOut")
	}

	var r0 backup.ProgressEventType
	var r1 bool
	if rf, ok := ret.Get(0).(func(backup.AnalysedMedia, backup.CatalogReference) (backup.ProgressEventType, bool)); ok {
		return rf(media, reference)
	}
	if rf, ok := ret.Get(0).(func(backup.AnalysedMedia, backup.CatalogReference) backup.ProgressEventType); ok {
		r0 = rf(media, reference)
	} else {
		r0 = ret.Get(0).(backup.ProgressEventType)
	}

	if rf, ok := ret.Get(1).(func(backup.AnalysedMedia, backup.CatalogReference) bool); ok {
		r1 = rf(media, reference)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// CatalogerFilter_FilterOut_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FilterOut'
type CatalogerFilter_FilterOut_Call struct {
	*mock.Call
}

// FilterOut is a helper method to define mock.On call
//   - media backup.AnalysedMedia
//   - reference backup.CatalogReference
func (_e *CatalogerFilter_Expecter) FilterOut(media interface{}, reference interface{}) *CatalogerFilter_FilterOut_Call {
	return &CatalogerFilter_FilterOut_Call{Call: _e.mock.On("FilterOut", media, reference)}
}

func (_c *CatalogerFilter_FilterOut_Call) Run(run func(media backup.AnalysedMedia, reference backup.CatalogReference)) *CatalogerFilter_FilterOut_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(backup.AnalysedMedia), args[1].(backup.CatalogReference))
	})
	return _c
}

func (_c *CatalogerFilter_FilterOut_Call) Return(_a0 backup.ProgressEventType, _a1 bool) *CatalogerFilter_FilterOut_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CatalogerFilter_FilterOut_Call) RunAndReturn(run func(backup.AnalysedMedia, backup.CatalogReference) (backup.ProgressEventType, bool)) *CatalogerFilter_FilterOut_Call {
	_c.Call.Return(run)
	return _c
}

// NewCatalogerFilter creates a new instance of CatalogerFilter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCatalogerFilter(t interface {
	mock.TestingT
	Cleanup(func())
}) *CatalogerFilter {
	mock := &CatalogerFilter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}