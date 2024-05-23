// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// AlbumLookupPort is an autogenerated mock type for the AlbumLookupPort type
type AlbumLookupPort struct {
	mock.Mock
}

type AlbumLookupPort_Expecter struct {
	mock *mock.Mock
}

func (_m *AlbumLookupPort) EXPECT() *AlbumLookupPort_Expecter {
	return &AlbumLookupPort_Expecter{mock: &_m.Mock}
}

// FindOrCreateAlbum provides a mock function with given fields: owner, mediaTime
func (_m *AlbumLookupPort) FindOrCreateAlbum(owner ownermodel.Owner, mediaTime time.Time) (string, bool, error) {
	ret := _m.Called(owner, mediaTime)

	if len(ret) == 0 {
		panic("no return value specified for FindOrCreateAlbum")
	}

	var r0 string
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(ownermodel.Owner, time.Time) (string, bool, error)); ok {
		return rf(owner, mediaTime)
	}
	if rf, ok := ret.Get(0).(func(ownermodel.Owner, time.Time) string); ok {
		r0 = rf(owner, mediaTime)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(ownermodel.Owner, time.Time) bool); ok {
		r1 = rf(owner, mediaTime)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(ownermodel.Owner, time.Time) error); ok {
		r2 = rf(owner, mediaTime)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// AlbumLookupPort_FindOrCreateAlbum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOrCreateAlbum'
type AlbumLookupPort_FindOrCreateAlbum_Call struct {
	*mock.Call
}

// FindOrCreateAlbum is a helper method to define mock.On call
//   - owner ownermodel.Owner
//   - mediaTime time.Time
func (_e *AlbumLookupPort_Expecter) FindOrCreateAlbum(owner interface{}, mediaTime interface{}) *AlbumLookupPort_FindOrCreateAlbum_Call {
	return &AlbumLookupPort_FindOrCreateAlbum_Call{Call: _e.mock.On("FindOrCreateAlbum", owner, mediaTime)}
}

func (_c *AlbumLookupPort_FindOrCreateAlbum_Call) Run(run func(owner ownermodel.Owner, mediaTime time.Time)) *AlbumLookupPort_FindOrCreateAlbum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ownermodel.Owner), args[1].(time.Time))
	})
	return _c
}

func (_c *AlbumLookupPort_FindOrCreateAlbum_Call) Return(folderName string, created bool, err error) *AlbumLookupPort_FindOrCreateAlbum_Call {
	_c.Call.Return(folderName, created, err)
	return _c
}

func (_c *AlbumLookupPort_FindOrCreateAlbum_Call) RunAndReturn(run func(ownermodel.Owner, time.Time) (string, bool, error)) *AlbumLookupPort_FindOrCreateAlbum_Call {
	_c.Call.Return(run)
	return _c
}

// NewAlbumLookupPort creates a new instance of AlbumLookupPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAlbumLookupPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *AlbumLookupPort {
	mock := &AlbumLookupPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
