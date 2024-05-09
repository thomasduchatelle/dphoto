// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// CreateAlbumObserverFunc is an autogenerated mock type for the CreateAlbumObserverFunc type
type CreateAlbumObserverFunc struct {
	mock.Mock
}

type CreateAlbumObserverFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *CreateAlbumObserverFunc) EXPECT() *CreateAlbumObserverFunc_Expecter {
	return &CreateAlbumObserverFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, album, records
func (_m *CreateAlbumObserverFunc) Execute(ctx context.Context, album catalog.Album, records catalog.MediaTransferRecords) error {
	ret := _m.Called(ctx, album, records)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.Album, catalog.MediaTransferRecords) error); ok {
		r0 = rf(ctx, album, records)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateAlbumObserverFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type CreateAlbumObserverFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - album catalog.Album
//   - records catalog.MediaTransferRecords
func (_e *CreateAlbumObserverFunc_Expecter) Execute(ctx interface{}, album interface{}, records interface{}) *CreateAlbumObserverFunc_Execute_Call {
	return &CreateAlbumObserverFunc_Execute_Call{Call: _e.mock.On("Execute", ctx, album, records)}
}

func (_c *CreateAlbumObserverFunc_Execute_Call) Run(run func(ctx context.Context, album catalog.Album, records catalog.MediaTransferRecords)) *CreateAlbumObserverFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(catalog.Album), args[2].(catalog.MediaTransferRecords))
	})
	return _c
}

func (_c *CreateAlbumObserverFunc_Execute_Call) Return(_a0 error) *CreateAlbumObserverFunc_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CreateAlbumObserverFunc_Execute_Call) RunAndReturn(run func(context.Context, catalog.Album, catalog.MediaTransferRecords) error) *CreateAlbumObserverFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewCreateAlbumObserverFunc creates a new instance of CreateAlbumObserverFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCreateAlbumObserverFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *CreateAlbumObserverFunc {
	mock := &CreateAlbumObserverFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
