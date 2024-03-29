// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// BArchiveAdapter is an autogenerated mock type for the BArchiveAdapter type
type BArchiveAdapter struct {
	mock.Mock
}

// ArchiveMedia provides a mock function with given fields: owner, media
func (_m *BArchiveAdapter) ArchiveMedia(owner string, media *backup.BackingUpMediaRequest) (string, error) {
	ret := _m.Called(owner, media)

	if len(ret) == 0 {
		panic("no return value specified for ArchiveMedia")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *backup.BackingUpMediaRequest) (string, error)); ok {
		return rf(owner, media)
	}
	if rf, ok := ret.Get(0).(func(string, *backup.BackingUpMediaRequest) string); ok {
		r0 = rf(owner, media)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, *backup.BackingUpMediaRequest) error); ok {
		r1 = rf(owner, media)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBArchiveAdapter creates a new instance of BArchiveAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBArchiveAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *BArchiveAdapter {
	mock := &BArchiveAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
