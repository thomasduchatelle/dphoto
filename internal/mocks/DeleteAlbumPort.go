// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DeleteAlbumPort is an autogenerated mock type for the DeleteAlbumPort type
type DeleteAlbumPort struct {
	mock.Mock
}

type DeleteAlbumPort_Expecter struct {
	mock *mock.Mock
}

func (_m *DeleteAlbumPort) EXPECT() *DeleteAlbumPort_Expecter {
	return &DeleteAlbumPort_Expecter{mock: &_m.Mock}
}

// DeleteAlbum provides a mock function with given fields: folderName
func (_m *DeleteAlbumPort) DeleteAlbum(folderName string) error {
	ret := _m.Called(folderName)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAlbum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(folderName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAlbumPort_DeleteAlbum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAlbum'
type DeleteAlbumPort_DeleteAlbum_Call struct {
	*mock.Call
}

// DeleteAlbum is a helper method to define mock.On call
//   - folderName string
func (_e *DeleteAlbumPort_Expecter) DeleteAlbum(folderName interface{}) *DeleteAlbumPort_DeleteAlbum_Call {
	return &DeleteAlbumPort_DeleteAlbum_Call{Call: _e.mock.On("DeleteAlbum", folderName)}
}

func (_c *DeleteAlbumPort_DeleteAlbum_Call) Run(run func(folderName string)) *DeleteAlbumPort_DeleteAlbum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *DeleteAlbumPort_DeleteAlbum_Call) Return(_a0 error) *DeleteAlbumPort_DeleteAlbum_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DeleteAlbumPort_DeleteAlbum_Call) RunAndReturn(run func(string) error) *DeleteAlbumPort_DeleteAlbum_Call {
	_c.Call.Return(run)
	return _c
}

// NewDeleteAlbumPort creates a new instance of DeleteAlbumPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDeleteAlbumPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *DeleteAlbumPort {
	mock := &DeleteAlbumPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
