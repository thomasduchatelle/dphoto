// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// PostAnalyseFilter is an autogenerated mock type for the PostAnalyseFilter type
type PostAnalyseFilter struct {
	mock.Mock
}

// AcceptAnalysedMedia provides a mock function with given fields: media, folderName
func (_m *PostAnalyseFilter) AcceptAnalysedMedia(media *backup.AnalysedMedia, folderName string) bool {
	ret := _m.Called(media, folderName)

	if len(ret) == 0 {
		panic("no return value specified for AcceptAnalysedMedia")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(*backup.AnalysedMedia, string) bool); ok {
		r0 = rf(media, folderName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewPostAnalyseFilter creates a new instance of PostAnalyseFilter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPostAnalyseFilter(t interface {
	mock.TestingT
	Cleanup(func())
}) *PostAnalyseFilter {
	mock := &PostAnalyseFilter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
