// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// UpdateAlbumPort is an autogenerated mock type for the UpdateAlbumPort type
type UpdateAlbumPort struct {
	mock.Mock
}

// UpdateAlbum provides a mock function with given fields: folderName, start, end
func (_m *UpdateAlbumPort) UpdateAlbum(folderName string, start time.Time, end time.Time) error {
	ret := _m.Called(folderName, start, end)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAlbum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, time.Time, time.Time) error); ok {
		r0 = rf(folderName, start, end)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUpdateAlbumPort creates a new instance of UpdateAlbumPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUpdateAlbumPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *UpdateAlbumPort {
	mock := &UpdateAlbumPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
