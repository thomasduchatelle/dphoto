// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// StorageAdapter is an autogenerated mock type for the StorageAdapter type
type StorageAdapter struct {
	mock.Mock
}

// ContentSignedUrl provides a mock function with given fields: owner, folderName, filename, expires
func (_m *StorageAdapter) ContentSignedUrl(owner string, folderName string, filename string, expires time.Duration) (string, error) {
	ret := _m.Called(owner, folderName, filename, expires)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string, time.Duration) string); ok {
		r0 = rf(owner, folderName, filename, expires)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, time.Duration) error); ok {
		r1 = rf(owner, folderName, filename, expires)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchFile provides a mock function with given fields: owner, folderName, filename
func (_m *StorageAdapter) FetchFile(owner string, folderName string, filename string) ([]byte, error) {
	ret := _m.Called(owner, folderName, filename)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, string, string) []byte); ok {
		r0 = rf(owner, folderName, filename)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(owner, folderName, filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStorageAdapter creates a new instance of StorageAdapter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorageAdapter(t testing.TB) *StorageAdapter {
	mock := &StorageAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
