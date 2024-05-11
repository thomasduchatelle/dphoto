// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// AmendAlbumDatesObserver is an autogenerated mock type for the AmendAlbumDatesObserver type
type AmendAlbumDatesObserver struct {
	mock.Mock
}

type AmendAlbumDatesObserver_Expecter struct {
	mock *mock.Mock
}

func (_m *AmendAlbumDatesObserver) EXPECT() *AmendAlbumDatesObserver_Expecter {
	return &AmendAlbumDatesObserver_Expecter{mock: &_m.Mock}
}

// OnAlbumDatesAmended provides a mock function with given fields: ctx, existingTimeline, updatedAlbum
func (_m *AmendAlbumDatesObserver) OnAlbumDatesAmended(ctx context.Context, existingTimeline []*catalog.Album, updatedAlbum catalog.Album) error {
	ret := _m.Called(ctx, existingTimeline, updatedAlbum)

	if len(ret) == 0 {
		panic("no return value specified for OnAlbumDatesAmended")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*catalog.Album, catalog.Album) error); ok {
		r0 = rf(ctx, existingTimeline, updatedAlbum)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AmendAlbumDatesObserver_OnAlbumDatesAmended_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnAlbumDatesAmended'
type AmendAlbumDatesObserver_OnAlbumDatesAmended_Call struct {
	*mock.Call
}

// OnAlbumDatesAmended is a helper method to define mock.On call
//   - ctx context.Context
//   - existingTimeline []*catalog.Album
//   - updatedAlbum catalog.Album
func (_e *AmendAlbumDatesObserver_Expecter) OnAlbumDatesAmended(ctx interface{}, existingTimeline interface{}, updatedAlbum interface{}) *AmendAlbumDatesObserver_OnAlbumDatesAmended_Call {
	return &AmendAlbumDatesObserver_OnAlbumDatesAmended_Call{Call: _e.mock.On("OnAlbumDatesAmended", ctx, existingTimeline, updatedAlbum)}
}

func (_c *AmendAlbumDatesObserver_OnAlbumDatesAmended_Call) Run(run func(ctx context.Context, existingTimeline []*catalog.Album, updatedAlbum catalog.Album)) *AmendAlbumDatesObserver_OnAlbumDatesAmended_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*catalog.Album), args[2].(catalog.Album))
	})
	return _c
}

func (_c *AmendAlbumDatesObserver_OnAlbumDatesAmended_Call) Return(_a0 error) *AmendAlbumDatesObserver_OnAlbumDatesAmended_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AmendAlbumDatesObserver_OnAlbumDatesAmended_Call) RunAndReturn(run func(context.Context, []*catalog.Album, catalog.Album) error) *AmendAlbumDatesObserver_OnAlbumDatesAmended_Call {
	_c.Call.Return(run)
	return _c
}

// NewAmendAlbumDatesObserver creates a new instance of AmendAlbumDatesObserver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAmendAlbumDatesObserver(t interface {
	mock.TestingT
	Cleanup(func())
}) *AmendAlbumDatesObserver {
	mock := &AmendAlbumDatesObserver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
