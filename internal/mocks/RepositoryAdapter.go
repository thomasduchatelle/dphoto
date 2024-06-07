// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"

	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// RepositoryAdapter is an autogenerated mock type for the RepositoryAdapter type
type RepositoryAdapter struct {
	mock.Mock
}

type RepositoryAdapter_Expecter struct {
	mock *mock.Mock
}

func (_m *RepositoryAdapter) EXPECT() *RepositoryAdapter_Expecter {
	return &RepositoryAdapter_Expecter{mock: &_m.Mock}
}

// CountMedia provides a mock function with given fields: ctx, album
func (_m *RepositoryAdapter) CountMedia(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	_va := make([]interface{}, len(album))
	for _i := range album {
		_va[_i] = album[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CountMedia")
	}

	var r0 map[catalog.AlbumId]int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...catalog.AlbumId) (map[catalog.AlbumId]int, error)); ok {
		return rf(ctx, album...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...catalog.AlbumId) map[catalog.AlbumId]int); ok {
		r0 = rf(ctx, album...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[catalog.AlbumId]int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...catalog.AlbumId) error); ok {
		r1 = rf(ctx, album...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RepositoryAdapter_CountMedia_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountMedia'
type RepositoryAdapter_CountMedia_Call struct {
	*mock.Call
}

// CountMedia is a helper method to define mock.On call
//   - ctx context.Context
//   - album ...catalog.AlbumId
func (_e *RepositoryAdapter_Expecter) CountMedia(ctx interface{}, album ...interface{}) *RepositoryAdapter_CountMedia_Call {
	return &RepositoryAdapter_CountMedia_Call{Call: _e.mock.On("CountMedia",
		append([]interface{}{ctx}, album...)...)}
}

func (_c *RepositoryAdapter_CountMedia_Call) Run(run func(ctx context.Context, album ...catalog.AlbumId)) *RepositoryAdapter_CountMedia_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]catalog.AlbumId, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(catalog.AlbumId)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *RepositoryAdapter_CountMedia_Call) Return(_a0 map[catalog.AlbumId]int, _a1 error) *RepositoryAdapter_CountMedia_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RepositoryAdapter_CountMedia_Call) RunAndReturn(run func(context.Context, ...catalog.AlbumId) (map[catalog.AlbumId]int, error)) *RepositoryAdapter_CountMedia_Call {
	_c.Call.Return(run)
	return _c
}

// FindAlbumByIds provides a mock function with given fields: ctx, ids
func (_m *RepositoryAdapter) FindAlbumByIds(ctx context.Context, ids ...catalog.AlbumId) ([]*catalog.Album, error) {
	_va := make([]interface{}, len(ids))
	for _i := range ids {
		_va[_i] = ids[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindAlbumByIds")
	}

	var r0 []*catalog.Album
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...catalog.AlbumId) ([]*catalog.Album, error)); ok {
		return rf(ctx, ids...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...catalog.AlbumId) []*catalog.Album); ok {
		r0 = rf(ctx, ids...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.Album)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...catalog.AlbumId) error); ok {
		r1 = rf(ctx, ids...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RepositoryAdapter_FindAlbumByIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAlbumByIds'
type RepositoryAdapter_FindAlbumByIds_Call struct {
	*mock.Call
}

// FindAlbumByIds is a helper method to define mock.On call
//   - ctx context.Context
//   - ids ...catalog.AlbumId
func (_e *RepositoryAdapter_Expecter) FindAlbumByIds(ctx interface{}, ids ...interface{}) *RepositoryAdapter_FindAlbumByIds_Call {
	return &RepositoryAdapter_FindAlbumByIds_Call{Call: _e.mock.On("FindAlbumByIds",
		append([]interface{}{ctx}, ids...)...)}
}

func (_c *RepositoryAdapter_FindAlbumByIds_Call) Run(run func(ctx context.Context, ids ...catalog.AlbumId)) *RepositoryAdapter_FindAlbumByIds_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]catalog.AlbumId, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(catalog.AlbumId)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *RepositoryAdapter_FindAlbumByIds_Call) Return(_a0 []*catalog.Album, _a1 error) *RepositoryAdapter_FindAlbumByIds_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RepositoryAdapter_FindAlbumByIds_Call) RunAndReturn(run func(context.Context, ...catalog.AlbumId) ([]*catalog.Album, error)) *RepositoryAdapter_FindAlbumByIds_Call {
	_c.Call.Return(run)
	return _c
}

// FindAlbumsByOwner provides a mock function with given fields: ctx, owner
func (_m *RepositoryAdapter) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
	ret := _m.Called(ctx, owner)

	if len(ret) == 0 {
		panic("no return value specified for FindAlbumsByOwner")
	}

	var r0 []*catalog.Album
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner) ([]*catalog.Album, error)); ok {
		return rf(ctx, owner)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner) []*catalog.Album); ok {
		r0 = rf(ctx, owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.Album)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ownermodel.Owner) error); ok {
		r1 = rf(ctx, owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RepositoryAdapter_FindAlbumsByOwner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAlbumsByOwner'
type RepositoryAdapter_FindAlbumsByOwner_Call struct {
	*mock.Call
}

// FindAlbumsByOwner is a helper method to define mock.On call
//   - ctx context.Context
//   - owner ownermodel.Owner
func (_e *RepositoryAdapter_Expecter) FindAlbumsByOwner(ctx interface{}, owner interface{}) *RepositoryAdapter_FindAlbumsByOwner_Call {
	return &RepositoryAdapter_FindAlbumsByOwner_Call{Call: _e.mock.On("FindAlbumsByOwner", ctx, owner)}
}

func (_c *RepositoryAdapter_FindAlbumsByOwner_Call) Run(run func(ctx context.Context, owner ownermodel.Owner)) *RepositoryAdapter_FindAlbumsByOwner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ownermodel.Owner))
	})
	return _c
}

func (_c *RepositoryAdapter_FindAlbumsByOwner_Call) Return(_a0 []*catalog.Album, _a1 error) *RepositoryAdapter_FindAlbumsByOwner_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RepositoryAdapter_FindAlbumsByOwner_Call) RunAndReturn(run func(context.Context, ownermodel.Owner) ([]*catalog.Album, error)) *RepositoryAdapter_FindAlbumsByOwner_Call {
	_c.Call.Return(run)
	return _c
}

// FindMediaCurrentAlbum provides a mock function with given fields: ctx, owner, mediaId
func (_m *RepositoryAdapter) FindMediaCurrentAlbum(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId) (*catalog.AlbumId, error) {
	ret := _m.Called(ctx, owner, mediaId)

	if len(ret) == 0 {
		panic("no return value specified for FindMediaCurrentAlbum")
	}

	var r0 *catalog.AlbumId
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner, catalog.MediaId) (*catalog.AlbumId, error)); ok {
		return rf(ctx, owner, mediaId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner, catalog.MediaId) *catalog.AlbumId); ok {
		r0 = rf(ctx, owner, mediaId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*catalog.AlbumId)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ownermodel.Owner, catalog.MediaId) error); ok {
		r1 = rf(ctx, owner, mediaId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RepositoryAdapter_FindMediaCurrentAlbum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMediaCurrentAlbum'
type RepositoryAdapter_FindMediaCurrentAlbum_Call struct {
	*mock.Call
}

// FindMediaCurrentAlbum is a helper method to define mock.On call
//   - ctx context.Context
//   - owner ownermodel.Owner
//   - mediaId catalog.MediaId
func (_e *RepositoryAdapter_Expecter) FindMediaCurrentAlbum(ctx interface{}, owner interface{}, mediaId interface{}) *RepositoryAdapter_FindMediaCurrentAlbum_Call {
	return &RepositoryAdapter_FindMediaCurrentAlbum_Call{Call: _e.mock.On("FindMediaCurrentAlbum", ctx, owner, mediaId)}
}

func (_c *RepositoryAdapter_FindMediaCurrentAlbum_Call) Run(run func(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId)) *RepositoryAdapter_FindMediaCurrentAlbum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ownermodel.Owner), args[2].(catalog.MediaId))
	})
	return _c
}

func (_c *RepositoryAdapter_FindMediaCurrentAlbum_Call) Return(id *catalog.AlbumId, err error) *RepositoryAdapter_FindMediaCurrentAlbum_Call {
	_c.Call.Return(id, err)
	return _c
}

func (_c *RepositoryAdapter_FindMediaCurrentAlbum_Call) RunAndReturn(run func(context.Context, ownermodel.Owner, catalog.MediaId) (*catalog.AlbumId, error)) *RepositoryAdapter_FindMediaCurrentAlbum_Call {
	_c.Call.Return(run)
	return _c
}

// FindMedias provides a mock function with given fields: ctx, request
func (_m *RepositoryAdapter) FindMedias(ctx context.Context, request *catalog.FindMediaRequest) ([]*catalog.MediaMeta, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for FindMedias")
	}

	var r0 []*catalog.MediaMeta
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *catalog.FindMediaRequest) ([]*catalog.MediaMeta, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *catalog.FindMediaRequest) []*catalog.MediaMeta); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.MediaMeta)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *catalog.FindMediaRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RepositoryAdapter_FindMedias_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMedias'
type RepositoryAdapter_FindMedias_Call struct {
	*mock.Call
}

// FindMedias is a helper method to define mock.On call
//   - ctx context.Context
//   - request *catalog.FindMediaRequest
func (_e *RepositoryAdapter_Expecter) FindMedias(ctx interface{}, request interface{}) *RepositoryAdapter_FindMedias_Call {
	return &RepositoryAdapter_FindMedias_Call{Call: _e.mock.On("FindMedias", ctx, request)}
}

func (_c *RepositoryAdapter_FindMedias_Call) Run(run func(ctx context.Context, request *catalog.FindMediaRequest)) *RepositoryAdapter_FindMedias_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*catalog.FindMediaRequest))
	})
	return _c
}

func (_c *RepositoryAdapter_FindMedias_Call) Return(medias []*catalog.MediaMeta, err error) *RepositoryAdapter_FindMedias_Call {
	_c.Call.Return(medias, err)
	return _c
}

func (_c *RepositoryAdapter_FindMedias_Call) RunAndReturn(run func(context.Context, *catalog.FindMediaRequest) ([]*catalog.MediaMeta, error)) *RepositoryAdapter_FindMedias_Call {
	_c.Call.Return(run)
	return _c
}

// NewRepositoryAdapter creates a new instance of RepositoryAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepositoryAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *RepositoryAdapter {
	mock := &RepositoryAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
