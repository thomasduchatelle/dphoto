// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// RunnerAnalyser is an autogenerated mock type for the RunnerAnalyser type
type RunnerAnalyser struct {
	mock.Mock
}

type RunnerAnalyser_Expecter struct {
	mock *mock.Mock
}

func (_m *RunnerAnalyser) EXPECT() *RunnerAnalyser_Expecter {
	return &RunnerAnalyser_Expecter{mock: &_m.Mock}
}

// Analyse provides a mock function with given fields: found, progressChannel
func (_m *RunnerAnalyser) Analyse(found backup.FoundMedia, progressChannel chan *backup.ProgressEvent) (*backup.AnalysedMedia, error) {
	ret := _m.Called(found, progressChannel)

	if len(ret) == 0 {
		panic("no return value specified for Analyse")
	}

	var r0 *backup.AnalysedMedia
	var r1 error
	if rf, ok := ret.Get(0).(func(backup.FoundMedia, chan *backup.ProgressEvent) (*backup.AnalysedMedia, error)); ok {
		return rf(found, progressChannel)
	}
	if rf, ok := ret.Get(0).(func(backup.FoundMedia, chan *backup.ProgressEvent) *backup.AnalysedMedia); ok {
		r0 = rf(found, progressChannel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*backup.AnalysedMedia)
		}
	}

	if rf, ok := ret.Get(1).(func(backup.FoundMedia, chan *backup.ProgressEvent) error); ok {
		r1 = rf(found, progressChannel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RunnerAnalyser_Analyse_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Analyse'
type RunnerAnalyser_Analyse_Call struct {
	*mock.Call
}

// Analyse is a helper method to define mock.On call
//   - found backup.FoundMedia
//   - progressChannel chan *backup.ProgressEvent
func (_e *RunnerAnalyser_Expecter) Analyse(found interface{}, progressChannel interface{}) *RunnerAnalyser_Analyse_Call {
	return &RunnerAnalyser_Analyse_Call{Call: _e.mock.On("Analyse", found, progressChannel)}
}

func (_c *RunnerAnalyser_Analyse_Call) Run(run func(found backup.FoundMedia, progressChannel chan *backup.ProgressEvent)) *RunnerAnalyser_Analyse_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(backup.FoundMedia), args[1].(chan *backup.ProgressEvent))
	})
	return _c
}

func (_c *RunnerAnalyser_Analyse_Call) Return(_a0 *backup.AnalysedMedia, _a1 error) *RunnerAnalyser_Analyse_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RunnerAnalyser_Analyse_Call) RunAndReturn(run func(backup.FoundMedia, chan *backup.ProgressEvent) (*backup.AnalysedMedia, error)) *RunnerAnalyser_Analyse_Call {
	_c.Call.Return(run)
	return _c
}

// NewRunnerAnalyser creates a new instance of RunnerAnalyser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRunnerAnalyser(t interface {
	mock.TestingT
	Cleanup(func())
}) *RunnerAnalyser {
	mock := &RunnerAnalyser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
