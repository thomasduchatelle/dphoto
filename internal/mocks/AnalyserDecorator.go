// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	backup "github.com/thomasduchatelle/dphoto/pkg/backup"
)

// AnalyserDecorator is an autogenerated mock type for the AnalyserDecorator type
type AnalyserDecorator struct {
	mock.Mock
}

type AnalyserDecorator_Expecter struct {
	mock *mock.Mock
}

func (_m *AnalyserDecorator) EXPECT() *AnalyserDecorator_Expecter {
	return &AnalyserDecorator_Expecter{mock: &_m.Mock}
}

// Decorate provides a mock function with given fields: analyseFunc
func (_m *AnalyserDecorator) Decorate(analyseFunc backup.RunnerAnalyser) backup.RunnerAnalyser {
	ret := _m.Called(analyseFunc)

	if len(ret) == 0 {
		panic("no return value specified for Decorate")
	}

	var r0 backup.RunnerAnalyser
	if rf, ok := ret.Get(0).(func(backup.RunnerAnalyser) backup.RunnerAnalyser); ok {
		r0 = rf(analyseFunc)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(backup.RunnerAnalyser)
		}
	}

	return r0
}

// AnalyserDecorator_Decorate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Decorate'
type AnalyserDecorator_Decorate_Call struct {
	*mock.Call
}

// Decorate is a helper method to define mock.On call
//   - analyseFunc backup.RunnerAnalyser
func (_e *AnalyserDecorator_Expecter) Decorate(analyseFunc interface{}) *AnalyserDecorator_Decorate_Call {
	return &AnalyserDecorator_Decorate_Call{Call: _e.mock.On("Decorate", analyseFunc)}
}

func (_c *AnalyserDecorator_Decorate_Call) Run(run func(analyseFunc backup.RunnerAnalyser)) *AnalyserDecorator_Decorate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(backup.RunnerAnalyser))
	})
	return _c
}

func (_c *AnalyserDecorator_Decorate_Call) Return(_a0 backup.RunnerAnalyser) *AnalyserDecorator_Decorate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AnalyserDecorator_Decorate_Call) RunAndReturn(run func(backup.RunnerAnalyser) backup.RunnerAnalyser) *AnalyserDecorator_Decorate_Call {
	_c.Call.Return(run)
	return _c
}

// NewAnalyserDecorator creates a new instance of AnalyserDecorator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAnalyserDecorator(t interface {
	mock.TestingT
	Cleanup(func())
}) *AnalyserDecorator {
	mock := &AnalyserDecorator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
