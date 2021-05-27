// Code generated by mockery v2.3.0. DO NOT EDIT.

package backup

import mock "github.com/stretchr/testify/mock"

// MockOnlineStorageAdapter is an autogenerated mock type for the OnlineStorageAdapter type
type MockOnlineStorageAdapter struct {
	mock.Mock
}

// UploadFile provides a mock function with given fields: media, folderName, filename
func (_m *MockOnlineStorageAdapter) UploadFile(media ReadableMedia, folderName string, filename string) (string, error) {
	ret := _m.Called(media, folderName, filename)

	var r0 string
	if rf, ok := ret.Get(0).(func(ReadableMedia, string, string) string); ok {
		r0 = rf(media, folderName, filename)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ReadableMedia, string, string) error); ok {
		r1 = rf(media, folderName, filename)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}