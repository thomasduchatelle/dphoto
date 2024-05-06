// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"

	mock "github.com/stretchr/testify/mock"
)

// DynamoBatchGetItem is an autogenerated mock type for the DynamoBatchGetItem type
type DynamoBatchGetItem struct {
	mock.Mock
}

type DynamoBatchGetItem_Expecter struct {
	mock *mock.Mock
}

func (_m *DynamoBatchGetItem) EXPECT() *DynamoBatchGetItem_Expecter {
	return &DynamoBatchGetItem_Expecter{mock: &_m.Mock}
}

// BatchGetItem provides a mock function with given fields: ctx, params, optFns
func (_m *DynamoBatchGetItem) BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for BatchGetItem")
	}

	var r0 *dynamodb.BatchGetItemOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dynamodb.BatchGetItemInput, ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dynamodb.BatchGetItemInput, ...func(*dynamodb.Options)) *dynamodb.BatchGetItemOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.BatchGetItemOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dynamodb.BatchGetItemInput, ...func(*dynamodb.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DynamoBatchGetItem_BatchGetItem_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BatchGetItem'
type DynamoBatchGetItem_BatchGetItem_Call struct {
	*mock.Call
}

// BatchGetItem is a helper method to define mock.On call
//   - ctx context.Context
//   - params *dynamodb.BatchGetItemInput
//   - optFns ...func(*dynamodb.Options)
func (_e *DynamoBatchGetItem_Expecter) BatchGetItem(ctx interface{}, params interface{}, optFns ...interface{}) *DynamoBatchGetItem_BatchGetItem_Call {
	return &DynamoBatchGetItem_BatchGetItem_Call{Call: _e.mock.On("BatchGetItem",
		append([]interface{}{ctx, params}, optFns...)...)}
}

func (_c *DynamoBatchGetItem_BatchGetItem_Call) Run(run func(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options))) *DynamoBatchGetItem_BatchGetItem_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]func(*dynamodb.Options), len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(func(*dynamodb.Options))
			}
		}
		run(args[0].(context.Context), args[1].(*dynamodb.BatchGetItemInput), variadicArgs...)
	})
	return _c
}

func (_c *DynamoBatchGetItem_BatchGetItem_Call) Return(_a0 *dynamodb.BatchGetItemOutput, _a1 error) *DynamoBatchGetItem_BatchGetItem_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DynamoBatchGetItem_BatchGetItem_Call) RunAndReturn(run func(context.Context, *dynamodb.BatchGetItemInput, ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)) *DynamoBatchGetItem_BatchGetItem_Call {
	_c.Call.Return(run)
	return _c
}

// NewDynamoBatchGetItem creates a new instance of DynamoBatchGetItem. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDynamoBatchGetItem(t interface {
	mock.TestingT
	Cleanup(func())
}) *DynamoBatchGetItem {
	mock := &DynamoBatchGetItem{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
