// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"

	mock "github.com/stretchr/testify/mock"

	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// IdentityQueriesScopeRepository is an autogenerated mock type for the IdentityQueriesScopeRepository type
type IdentityQueriesScopeRepository struct {
	mock.Mock
}

type IdentityQueriesScopeRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *IdentityQueriesScopeRepository) EXPECT() *IdentityQueriesScopeRepository_Expecter {
	return &IdentityQueriesScopeRepository_Expecter{mock: &_m.Mock}
}

// ListScopesByOwners provides a mock function with given fields: ctx, owners, types
func (_m *IdentityQueriesScopeRepository) ListScopesByOwners(ctx context.Context, owners []ownermodel.Owner, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	_va := make([]interface{}, len(types))
	for _i := range types {
		_va[_i] = types[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, owners)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListScopesByOwners")
	}

	var r0 []*aclcore.Scope
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []ownermodel.Owner, ...aclcore.ScopeType) ([]*aclcore.Scope, error)); ok {
		return rf(ctx, owners, types...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []ownermodel.Owner, ...aclcore.ScopeType) []*aclcore.Scope); ok {
		r0 = rf(ctx, owners, types...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*aclcore.Scope)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []ownermodel.Owner, ...aclcore.ScopeType) error); ok {
		r1 = rf(ctx, owners, types...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IdentityQueriesScopeRepository_ListScopesByOwners_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListScopesByOwners'
type IdentityQueriesScopeRepository_ListScopesByOwners_Call struct {
	*mock.Call
}

// ListScopesByOwners is a helper method to define mock.On call
//   - ctx context.Context
//   - owners []ownermodel.Owner
//   - types ...aclcore.ScopeType
func (_e *IdentityQueriesScopeRepository_Expecter) ListScopesByOwners(ctx interface{}, owners interface{}, types ...interface{}) *IdentityQueriesScopeRepository_ListScopesByOwners_Call {
	return &IdentityQueriesScopeRepository_ListScopesByOwners_Call{Call: _e.mock.On("ListScopesByOwners",
		append([]interface{}{ctx, owners}, types...)...)}
}

func (_c *IdentityQueriesScopeRepository_ListScopesByOwners_Call) Run(run func(ctx context.Context, owners []ownermodel.Owner, types ...aclcore.ScopeType)) *IdentityQueriesScopeRepository_ListScopesByOwners_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]aclcore.ScopeType, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(aclcore.ScopeType)
			}
		}
		run(args[0].(context.Context), args[1].([]ownermodel.Owner), variadicArgs...)
	})
	return _c
}

func (_c *IdentityQueriesScopeRepository_ListScopesByOwners_Call) Return(_a0 []*aclcore.Scope, _a1 error) *IdentityQueriesScopeRepository_ListScopesByOwners_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IdentityQueriesScopeRepository_ListScopesByOwners_Call) RunAndReturn(run func(context.Context, []ownermodel.Owner, ...aclcore.ScopeType) ([]*aclcore.Scope, error)) *IdentityQueriesScopeRepository_ListScopesByOwners_Call {
	_c.Call.Return(run)
	return _c
}

// NewIdentityQueriesScopeRepository creates a new instance of IdentityQueriesScopeRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIdentityQueriesScopeRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IdentityQueriesScopeRepository {
	mock := &IdentityQueriesScopeRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
