// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	dnsdomain "github.com/thomasduchatelle/dphoto/pkg/dnsdomain"
)

// CertificateAuthority is an autogenerated mock type for the CertificateAuthority type
type CertificateAuthority struct {
	mock.Mock
}

// RequestCertificate provides a mock function with given fields: email, domain
func (_m *CertificateAuthority) RequestCertificate(email string, domain string) (*dnsdomain.CompleteCertificate, error) {
	ret := _m.Called(email, domain)

	if len(ret) == 0 {
		panic("no return value specified for RequestCertificate")
	}

	var r0 *dnsdomain.CompleteCertificate
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*dnsdomain.CompleteCertificate, error)); ok {
		return rf(email, domain)
	}
	if rf, ok := ret.Get(0).(func(string, string) *dnsdomain.CompleteCertificate); ok {
		r0 = rf(email, domain)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dnsdomain.CompleteCertificate)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(email, domain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCertificateAuthority creates a new instance of CertificateAuthority. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCertificateAuthority(t interface {
	mock.TestingT
	Cleanup(func())
}) *CertificateAuthority {
	mock := &CertificateAuthority{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
