// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"

	mock "github.com/stretchr/testify/mock"
)

// DynamoBatchGetItem is an autogenerated mock type for the DynamoBatchGetItem type
type DynamoBatchGetItem struct {
	mock.Mock
}

// BatchGetItem provides a mock function with given fields: _a0
func (_m *DynamoBatchGetItem) BatchGetItem(_a0 *dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for BatchGetItem")
	}

	var r0 *dynamodb.BatchGetItemOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(*dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*dynamodb.BatchGetItemInput) *dynamodb.BatchGetItemOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.BatchGetItemOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(*dynamodb.BatchGetItemInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
