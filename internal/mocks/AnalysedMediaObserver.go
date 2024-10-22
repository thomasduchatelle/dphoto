// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	backup "github.com/thomasduchatelle/dphoto/pkg/backup"

	mock "github.com/stretchr/testify/mock"
)

// AnalysedMediaObserver is an autogenerated mock type for the AnalysedMediaObserver type
type AnalysedMediaObserver struct {
	mock.Mock
}

type AnalysedMediaObserver_Expecter struct {
	mock *mock.Mock
}

func (_m *AnalysedMediaObserver) EXPECT() *AnalysedMediaObserver_Expecter {
	return &AnalysedMediaObserver_Expecter{mock: &_m.Mock}
}

// OnAnalysedMedia provides a mock function with given fields: ctx, media
func (_m *AnalysedMediaObserver) OnAnalysedMedia(ctx context.Context, media *backup.AnalysedMedia) error {
	ret := _m.Called(ctx, media)

	if len(ret) == 0 {
		panic("no return value specified for OnAnalysedMedia")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *backup.AnalysedMedia) error); ok {
		r0 = rf(ctx, media)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AnalysedMediaObserver_OnAnalysedMedia_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnAnalysedMedia'
type AnalysedMediaObserver_OnAnalysedMedia_Call struct {
	*mock.Call
}

// OnAnalysedMedia is a helper method to define mock.On call
//   - ctx context.Context
//   - media *backup.AnalysedMedia
func (_e *AnalysedMediaObserver_Expecter) OnAnalysedMedia(ctx interface{}, media interface{}) *AnalysedMediaObserver_OnAnalysedMedia_Call {
	return &AnalysedMediaObserver_OnAnalysedMedia_Call{Call: _e.mock.On("OnAnalysedMedia", ctx, media)}
}

func (_c *AnalysedMediaObserver_OnAnalysedMedia_Call) Run(run func(ctx context.Context, media *backup.AnalysedMedia)) *AnalysedMediaObserver_OnAnalysedMedia_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*backup.AnalysedMedia))
	})
	return _c
}

func (_c *AnalysedMediaObserver_OnAnalysedMedia_Call) Return(_a0 error) *AnalysedMediaObserver_OnAnalysedMedia_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AnalysedMediaObserver_OnAnalysedMedia_Call) RunAndReturn(run func(context.Context, *backup.AnalysedMedia) error) *AnalysedMediaObserver_OnAnalysedMedia_Call {
	_c.Call.Return(run)
	return _c
}

// NewAnalysedMediaObserver creates a new instance of AnalysedMediaObserver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAnalysedMediaObserver(t interface {
	mock.TestingT
	Cleanup(func())
}) *AnalysedMediaObserver {
	mock := &AnalysedMediaObserver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
