// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	catalog "github.com/thomasduchatelle/dphoto/domain/catalog"

	testing "testing"
)

// RepositoryPort is an autogenerated mock type for the RepositoryPort type
type RepositoryPort struct {
	mock.Mock
}

// CountMedias provides a mock function with given fields: owner, folderName
func (_m *RepositoryPort) CountMedias(owner string, folderName string) (int, error) {
	ret := _m.Called(owner, folderName)

	var r0 int
	if rf, ok := ret.Get(0).(func(string, string) int); ok {
		r0 = rf(owner, folderName)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(owner, folderName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteEmptyAlbum provides a mock function with given fields: owner, folderName
func (_m *RepositoryPort) DeleteEmptyAlbum(owner string, folderName string) error {
	ret := _m.Called(owner, folderName)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(owner, folderName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteEmptyMoveTransaction provides a mock function with given fields: transactionId
func (_m *RepositoryPort) DeleteEmptyMoveTransaction(transactionId string) error {
	ret := _m.Called(transactionId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(transactionId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAlbum provides a mock function with given fields: owner, folderName
func (_m *RepositoryPort) FindAlbum(owner string, folderName string) (*catalog.Album, error) {
	ret := _m.Called(owner, folderName)

	var r0 *catalog.Album
	if rf, ok := ret.Get(0).(func(string, string) *catalog.Album); ok {
		r0 = rf(owner, folderName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*catalog.Album)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(owner, folderName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllAlbums provides a mock function with given fields: owner
func (_m *RepositoryPort) FindAllAlbums(owner string) ([]*catalog.Album, error) {
	ret := _m.Called(owner)

	var r0 []*catalog.Album
	if rf, ok := ret.Get(0).(func(string) []*catalog.Album); ok {
		r0 = rf(owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.Album)
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

// FindExistingSignatures provides a mock function with given fields: owner, signatures
func (_m *RepositoryPort) FindExistingSignatures(owner string, signatures []*catalog.MediaSignature) ([]*catalog.MediaSignature, error) {
	ret := _m.Called(owner, signatures)

	var r0 []*catalog.MediaSignature
	if rf, ok := ret.Get(0).(func(string, []*catalog.MediaSignature) []*catalog.MediaSignature); ok {
		r0 = rf(owner, signatures)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.MediaSignature)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []*catalog.MediaSignature) error); ok {
		r1 = rf(owner, signatures)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindFilesToMove provides a mock function with given fields: transactionId, pageToken
func (_m *RepositoryPort) FindFilesToMove(transactionId string, pageToken string) ([]*catalog.MovedMedia, string, error) {
	ret := _m.Called(transactionId, pageToken)

	var r0 []*catalog.MovedMedia
	if rf, ok := ret.Get(0).(func(string, string) []*catalog.MovedMedia); ok {
		r0 = rf(transactionId, pageToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.MovedMedia)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string, string) string); ok {
		r1 = rf(transactionId, pageToken)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(transactionId, pageToken)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FindMediaLocations provides a mock function with given fields: owner, signature
func (_m *RepositoryPort) FindMediaLocations(owner string, signature catalog.MediaSignature) ([]*catalog.MediaLocation, error) {
	ret := _m.Called(owner, signature)

	var r0 []*catalog.MediaLocation
	if rf, ok := ret.Get(0).(func(string, catalog.MediaSignature) []*catalog.MediaLocation); ok {
		r0 = rf(owner, signature)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.MediaLocation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, catalog.MediaSignature) error); ok {
		r1 = rf(owner, signature)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindMedias provides a mock function with given fields: owner, folderName, filter
func (_m *RepositoryPort) FindMedias(owner string, folderName string, filter catalog.FindMediaFilter) (*catalog.MediaPage, error) {
	ret := _m.Called(owner, folderName, filter)

	var r0 *catalog.MediaPage
	if rf, ok := ret.Get(0).(func(string, string, catalog.FindMediaFilter) *catalog.MediaPage); ok {
		r0 = rf(owner, folderName, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*catalog.MediaPage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, catalog.FindMediaFilter) error); ok {
		r1 = rf(owner, folderName, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindReadyMoveTransactions provides a mock function with given fields: owner
func (_m *RepositoryPort) FindReadyMoveTransactions(owner string) ([]*catalog.MoveTransaction, error) {
	ret := _m.Called(owner)

	var r0 []*catalog.MoveTransaction
	if rf, ok := ret.Get(0).(func(string) []*catalog.MoveTransaction); ok {
		r0 = rf(owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.MoveTransaction)
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

// InsertAlbum provides a mock function with given fields: album
func (_m *RepositoryPort) InsertAlbum(album catalog.Album) error {
	ret := _m.Called(album)

	var r0 error
	if rf, ok := ret.Get(0).(func(catalog.Album) error); ok {
		r0 = rf(album)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertMedias provides a mock function with given fields: owner, media
func (_m *RepositoryPort) InsertMedias(owner string, media []catalog.CreateMediaRequest) error {
	ret := _m.Called(owner, media)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []catalog.CreateMediaRequest) error); ok {
		r0 = rf(owner, media)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAlbum provides a mock function with given fields: album
func (_m *RepositoryPort) UpdateAlbum(album catalog.Album) error {
	ret := _m.Called(album)

	var r0 error
	if rf, ok := ret.Get(0).(func(catalog.Album) error); ok {
		r0 = rf(album)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateMedias provides a mock function with given fields: filter, newFolderName
func (_m *RepositoryPort) UpdateMedias(filter *catalog.UpdateMediaFilter, newFolderName string) (string, int, error) {
	ret := _m.Called(filter, newFolderName)

	var r0 string
	if rf, ok := ret.Get(0).(func(*catalog.UpdateMediaFilter, string) string); ok {
		r0 = rf(filter, newFolderName)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(*catalog.UpdateMediaFilter, string) int); ok {
		r1 = rf(filter, newFolderName)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*catalog.UpdateMediaFilter, string) error); ok {
		r2 = rf(filter, newFolderName)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UpdateMediasLocation provides a mock function with given fields: owner, transactionId, moves
func (_m *RepositoryPort) UpdateMediasLocation(owner string, transactionId string, moves []*catalog.MovedMedia) error {
	ret := _m.Called(owner, transactionId, moves)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, []*catalog.MovedMedia) error); ok {
		r0 = rf(owner, transactionId, moves)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepositoryPort creates a new instance of RepositoryPort. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepositoryPort(t testing.TB) *RepositoryPort {
	mock := &RepositoryPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
