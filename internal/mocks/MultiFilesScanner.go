// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	backup "github.com/thomasduchatelle/dphoto/pkg/backup"

	mock "github.com/stretchr/testify/mock"
)

// MultiFilesScanner is an autogenerated mock type for the MultiFilesScanner type
type MultiFilesScanner struct {
	mock.Mock
}

type MultiFilesScanner_Expecter struct {
	mock *mock.Mock
}

func (_m *MultiFilesScanner) EXPECT() *MultiFilesScanner_Expecter {
	return &MultiFilesScanner_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, owner, volume, optionSlice
func (_m *MultiFilesScanner) Execute(ctx context.Context, owner string, volume backup.SourceVolume, optionSlice ...backup.Options) ([]*backup.ScannedFolder, error) {
	_va := make([]interface{}, len(optionSlice))
	for _i := range optionSlice {
		_va[_i] = optionSlice[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, owner, volume)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 []*backup.ScannedFolder
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, backup.SourceVolume, ...backup.Options) ([]*backup.ScannedFolder, error)); ok {
		return rf(ctx, owner, volume, optionSlice...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, backup.SourceVolume, ...backup.Options) []*backup.ScannedFolder); ok {
		r0 = rf(ctx, owner, volume, optionSlice...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*backup.ScannedFolder)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, backup.SourceVolume, ...backup.Options) error); ok {
		r1 = rf(ctx, owner, volume, optionSlice...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MultiFilesScanner_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MultiFilesScanner_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - owner string
//   - volume backup.SourceVolume
//   - optionSlice ...backup.Options
func (_e *MultiFilesScanner_Expecter) Execute(ctx interface{}, owner interface{}, volume interface{}, optionSlice ...interface{}) *MultiFilesScanner_Execute_Call {
	return &MultiFilesScanner_Execute_Call{Call: _e.mock.On("Execute",
		append([]interface{}{ctx, owner, volume}, optionSlice...)...)}
}

func (_c *MultiFilesScanner_Execute_Call) Run(run func(ctx context.Context, owner string, volume backup.SourceVolume, optionSlice ...backup.Options)) *MultiFilesScanner_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]backup.Options, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(backup.Options)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(backup.SourceVolume), variadicArgs...)
	})
	return _c
}

func (_c *MultiFilesScanner_Execute_Call) Return(_a0 []*backup.ScannedFolder, _a1 error) *MultiFilesScanner_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MultiFilesScanner_Execute_Call) RunAndReturn(run func(context.Context, string, backup.SourceVolume, ...backup.Options) ([]*backup.ScannedFolder, error)) *MultiFilesScanner_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMultiFilesScanner creates a new instance of MultiFilesScanner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMultiFilesScanner(t interface {
	mock.TestingT
	Cleanup(func())
}) *MultiFilesScanner {
	mock := &MultiFilesScanner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
