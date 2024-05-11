// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// TransferMediasRepositoryPort is an autogenerated mock type for the TransferMediasRepositoryPort type
type TransferMediasRepositoryPort struct {
	mock.Mock
}

type TransferMediasRepositoryPort_Expecter struct {
	mock *mock.Mock
}

func (_m *TransferMediasRepositoryPort) EXPECT() *TransferMediasRepositoryPort_Expecter {
	return &TransferMediasRepositoryPort_Expecter{mock: &_m.Mock}
}

// TransferMediasFromRecords provides a mock function with given fields: ctx, records
func (_m *TransferMediasRepositoryPort) TransferMediasFromRecords(ctx context.Context, records catalog.MediaTransferRecords) (catalog.TransferredMedias, error) {
	ret := _m.Called(ctx, records)

	if len(ret) == 0 {
		panic("no return value specified for TransferMediasFromRecords")
	}

	var r0 catalog.TransferredMedias
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.MediaTransferRecords) (catalog.TransferredMedias, error)); ok {
		return rf(ctx, records)
	}
	if rf, ok := ret.Get(0).(func(context.Context, catalog.MediaTransferRecords) catalog.TransferredMedias); ok {
		r0 = rf(ctx, records)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(catalog.TransferredMedias)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, catalog.MediaTransferRecords) error); ok {
		r1 = rf(ctx, records)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransferMediasRepositoryPort_TransferMediasFromRecords_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TransferMediasFromRecords'
type TransferMediasRepositoryPort_TransferMediasFromRecords_Call struct {
	*mock.Call
}

// TransferMediasFromRecords is a helper method to define mock.On call
//   - ctx context.Context
//   - records catalog.MediaTransferRecords
func (_e *TransferMediasRepositoryPort_Expecter) TransferMediasFromRecords(ctx interface{}, records interface{}) *TransferMediasRepositoryPort_TransferMediasFromRecords_Call {
	return &TransferMediasRepositoryPort_TransferMediasFromRecords_Call{Call: _e.mock.On("TransferMediasFromRecords", ctx, records)}
}

func (_c *TransferMediasRepositoryPort_TransferMediasFromRecords_Call) Run(run func(ctx context.Context, records catalog.MediaTransferRecords)) *TransferMediasRepositoryPort_TransferMediasFromRecords_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(catalog.MediaTransferRecords))
	})
	return _c
}

func (_c *TransferMediasRepositoryPort_TransferMediasFromRecords_Call) Return(_a0 catalog.TransferredMedias, _a1 error) *TransferMediasRepositoryPort_TransferMediasFromRecords_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TransferMediasRepositoryPort_TransferMediasFromRecords_Call) RunAndReturn(run func(context.Context, catalog.MediaTransferRecords) (catalog.TransferredMedias, error)) *TransferMediasRepositoryPort_TransferMediasFromRecords_Call {
	_c.Call.Return(run)
	return _c
}

// NewTransferMediasRepositoryPort creates a new instance of TransferMediasRepositoryPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransferMediasRepositoryPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransferMediasRepositoryPort {
	mock := &TransferMediasRepositoryPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
