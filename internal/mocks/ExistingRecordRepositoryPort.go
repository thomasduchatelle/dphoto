// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	ui "github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
)

// ExistingRecordRepositoryPort is an autogenerated mock type for the ExistingRecordRepositoryPort type
type ExistingRecordRepositoryPort struct {
	mock.Mock
}

// FindExistingRecords provides a mock function with given fields:
func (_m *ExistingRecordRepositoryPort) FindExistingRecords() ([]*ui.ExistingRecord, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FindExistingRecords")
	}

	var r0 []*ui.ExistingRecord
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*ui.ExistingRecord, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*ui.ExistingRecord); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ui.ExistingRecord)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewExistingRecordRepositoryPort creates a new instance of ExistingRecordRepositoryPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExistingRecordRepositoryPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExistingRecordRepositoryPort {
	mock := &ExistingRecordRepositoryPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
