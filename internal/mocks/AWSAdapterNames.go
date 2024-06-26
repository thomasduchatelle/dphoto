// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// AWSAdapterNames is an autogenerated mock type for the AWSAdapterNames type
type AWSAdapterNames struct {
	mock.Mock
}

type AWSAdapterNames_Expecter struct {
	mock *mock.Mock
}

func (_m *AWSAdapterNames) EXPECT() *AWSAdapterNames_Expecter {
	return &AWSAdapterNames_Expecter{mock: &_m.Mock}
}

// ArchiveCacheBucketName provides a mock function with given fields:
func (_m *AWSAdapterNames) ArchiveCacheBucketName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ArchiveCacheBucketName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// AWSAdapterNames_ArchiveCacheBucketName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ArchiveCacheBucketName'
type AWSAdapterNames_ArchiveCacheBucketName_Call struct {
	*mock.Call
}

// ArchiveCacheBucketName is a helper method to define mock.On call
func (_e *AWSAdapterNames_Expecter) ArchiveCacheBucketName() *AWSAdapterNames_ArchiveCacheBucketName_Call {
	return &AWSAdapterNames_ArchiveCacheBucketName_Call{Call: _e.mock.On("ArchiveCacheBucketName")}
}

func (_c *AWSAdapterNames_ArchiveCacheBucketName_Call) Run(run func()) *AWSAdapterNames_ArchiveCacheBucketName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AWSAdapterNames_ArchiveCacheBucketName_Call) Return(_a0 string) *AWSAdapterNames_ArchiveCacheBucketName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AWSAdapterNames_ArchiveCacheBucketName_Call) RunAndReturn(run func() string) *AWSAdapterNames_ArchiveCacheBucketName_Call {
	_c.Call.Return(run)
	return _c
}

// ArchiveJobsSNSARN provides a mock function with given fields:
func (_m *AWSAdapterNames) ArchiveJobsSNSARN() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ArchiveJobsSNSARN")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// AWSAdapterNames_ArchiveJobsSNSARN_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ArchiveJobsSNSARN'
type AWSAdapterNames_ArchiveJobsSNSARN_Call struct {
	*mock.Call
}

// ArchiveJobsSNSARN is a helper method to define mock.On call
func (_e *AWSAdapterNames_Expecter) ArchiveJobsSNSARN() *AWSAdapterNames_ArchiveJobsSNSARN_Call {
	return &AWSAdapterNames_ArchiveJobsSNSARN_Call{Call: _e.mock.On("ArchiveJobsSNSARN")}
}

func (_c *AWSAdapterNames_ArchiveJobsSNSARN_Call) Run(run func()) *AWSAdapterNames_ArchiveJobsSNSARN_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AWSAdapterNames_ArchiveJobsSNSARN_Call) Return(_a0 string) *AWSAdapterNames_ArchiveJobsSNSARN_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AWSAdapterNames_ArchiveJobsSNSARN_Call) RunAndReturn(run func() string) *AWSAdapterNames_ArchiveJobsSNSARN_Call {
	_c.Call.Return(run)
	return _c
}

// ArchiveJobsSQSURL provides a mock function with given fields:
func (_m *AWSAdapterNames) ArchiveJobsSQSURL() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ArchiveJobsSQSURL")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// AWSAdapterNames_ArchiveJobsSQSURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ArchiveJobsSQSURL'
type AWSAdapterNames_ArchiveJobsSQSURL_Call struct {
	*mock.Call
}

// ArchiveJobsSQSURL is a helper method to define mock.On call
func (_e *AWSAdapterNames_Expecter) ArchiveJobsSQSURL() *AWSAdapterNames_ArchiveJobsSQSURL_Call {
	return &AWSAdapterNames_ArchiveJobsSQSURL_Call{Call: _e.mock.On("ArchiveJobsSQSURL")}
}

func (_c *AWSAdapterNames_ArchiveJobsSQSURL_Call) Run(run func()) *AWSAdapterNames_ArchiveJobsSQSURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AWSAdapterNames_ArchiveJobsSQSURL_Call) Return(_a0 string) *AWSAdapterNames_ArchiveJobsSQSURL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AWSAdapterNames_ArchiveJobsSQSURL_Call) RunAndReturn(run func() string) *AWSAdapterNames_ArchiveJobsSQSURL_Call {
	_c.Call.Return(run)
	return _c
}

// ArchiveMainBucketName provides a mock function with given fields:
func (_m *AWSAdapterNames) ArchiveMainBucketName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ArchiveMainBucketName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// AWSAdapterNames_ArchiveMainBucketName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ArchiveMainBucketName'
type AWSAdapterNames_ArchiveMainBucketName_Call struct {
	*mock.Call
}

// ArchiveMainBucketName is a helper method to define mock.On call
func (_e *AWSAdapterNames_Expecter) ArchiveMainBucketName() *AWSAdapterNames_ArchiveMainBucketName_Call {
	return &AWSAdapterNames_ArchiveMainBucketName_Call{Call: _e.mock.On("ArchiveMainBucketName")}
}

func (_c *AWSAdapterNames_ArchiveMainBucketName_Call) Run(run func()) *AWSAdapterNames_ArchiveMainBucketName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AWSAdapterNames_ArchiveMainBucketName_Call) Return(_a0 string) *AWSAdapterNames_ArchiveMainBucketName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AWSAdapterNames_ArchiveMainBucketName_Call) RunAndReturn(run func() string) *AWSAdapterNames_ArchiveMainBucketName_Call {
	_c.Call.Return(run)
	return _c
}

// DynamoDBName provides a mock function with given fields:
func (_m *AWSAdapterNames) DynamoDBName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DynamoDBName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// AWSAdapterNames_DynamoDBName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DynamoDBName'
type AWSAdapterNames_DynamoDBName_Call struct {
	*mock.Call
}

// DynamoDBName is a helper method to define mock.On call
func (_e *AWSAdapterNames_Expecter) DynamoDBName() *AWSAdapterNames_DynamoDBName_Call {
	return &AWSAdapterNames_DynamoDBName_Call{Call: _e.mock.On("DynamoDBName")}
}

func (_c *AWSAdapterNames_DynamoDBName_Call) Run(run func()) *AWSAdapterNames_DynamoDBName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AWSAdapterNames_DynamoDBName_Call) Return(_a0 string) *AWSAdapterNames_DynamoDBName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AWSAdapterNames_DynamoDBName_Call) RunAndReturn(run func() string) *AWSAdapterNames_DynamoDBName_Call {
	_c.Call.Return(run)
	return _c
}

// NewAWSAdapterNames creates a new instance of AWSAdapterNames. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAWSAdapterNames(t interface {
	mock.TestingT
	Cleanup(func())
}) *AWSAdapterNames {
	mock := &AWSAdapterNames{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
