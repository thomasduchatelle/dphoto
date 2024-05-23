// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	backup "github.com/thomasduchatelle/dphoto/pkg/backup"

	mock "github.com/stretchr/testify/mock"
)

// RunnerCataloger is an autogenerated mock type for the RunnerCataloger type
type RunnerCataloger struct {
	mock.Mock
}

type RunnerCataloger_Expecter struct {
	mock *mock.Mock
}

func (_m *RunnerCataloger) EXPECT() *RunnerCataloger_Expecter {
	return &RunnerCataloger_Expecter{mock: &_m.Mock}
}

// Catalog provides a mock function with given fields: ctx, medias, progressChannel
func (_m *RunnerCataloger) Catalog(ctx context.Context, medias []*backup.AnalysedMedia, progressChannel chan *backup.ProgressEvent) ([]*backup.BackingUpMediaRequest, error) {
	ret := _m.Called(ctx, medias, progressChannel)

	if len(ret) == 0 {
		panic("no return value specified for Catalog")
	}

	var r0 []*backup.BackingUpMediaRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []*backup.AnalysedMedia, chan *backup.ProgressEvent) ([]*backup.BackingUpMediaRequest, error)); ok {
		return rf(ctx, medias, progressChannel)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []*backup.AnalysedMedia, chan *backup.ProgressEvent) []*backup.BackingUpMediaRequest); ok {
		r0 = rf(ctx, medias, progressChannel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*backup.BackingUpMediaRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []*backup.AnalysedMedia, chan *backup.ProgressEvent) error); ok {
		r1 = rf(ctx, medias, progressChannel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RunnerCataloger_Catalog_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Catalog'
type RunnerCataloger_Catalog_Call struct {
	*mock.Call
}

// Catalog is a helper method to define mock.On call
//   - ctx context.Context
//   - medias []*backup.AnalysedMedia
//   - progressChannel chan *backup.ProgressEvent
func (_e *RunnerCataloger_Expecter) Catalog(ctx interface{}, medias interface{}, progressChannel interface{}) *RunnerCataloger_Catalog_Call {
	return &RunnerCataloger_Catalog_Call{Call: _e.mock.On("Catalog", ctx, medias, progressChannel)}
}

func (_c *RunnerCataloger_Catalog_Call) Run(run func(ctx context.Context, medias []*backup.AnalysedMedia, progressChannel chan *backup.ProgressEvent)) *RunnerCataloger_Catalog_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*backup.AnalysedMedia), args[2].(chan *backup.ProgressEvent))
	})
	return _c
}

func (_c *RunnerCataloger_Catalog_Call) Return(_a0 []*backup.BackingUpMediaRequest, _a1 error) *RunnerCataloger_Catalog_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RunnerCataloger_Catalog_Call) RunAndReturn(run func(context.Context, []*backup.AnalysedMedia, chan *backup.ProgressEvent) ([]*backup.BackingUpMediaRequest, error)) *RunnerCataloger_Catalog_Call {
	_c.Call.Return(run)
	return _c
}

// NewRunnerCataloger creates a new instance of RunnerCataloger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRunnerCataloger(t interface {
	mock.TestingT
	Cleanup(func())
}) *RunnerCataloger {
	mock := &RunnerCataloger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
