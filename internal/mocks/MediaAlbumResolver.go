// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MediaAlbumResolver is an autogenerated mock type for the MediaAlbumResolver type
type MediaAlbumResolver struct {
	mock.Mock
}

// FindAlbumOfMedia provides a mock function with given fields: owner, mediaId
func (_m *MediaAlbumResolver) FindAlbumOfMedia(owner string, mediaId string) (string, error) {
	ret := _m.Called(owner, mediaId)

	if len(ret) == 0 {
		panic("no return value specified for FindAlbumOfMedia")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (string, error)); ok {
		return rf(owner, mediaId)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(owner, mediaId)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(owner, mediaId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMediaAlbumResolver creates a new instance of MediaAlbumResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMediaAlbumResolver(t interface {
	mock.TestingT
	Cleanup(func())
}) *MediaAlbumResolver {
	mock := &MediaAlbumResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
