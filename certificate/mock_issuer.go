// Code generated by mockery v1.0.0. DO NOT EDIT.

package certificate

import (
	tls "crypto/tls"

	mock "github.com/stretchr/testify/mock"
)

// MockIssuer is an autogenerated mock type for the Issuer type
type MockIssuer struct {
	mock.Mock
}

// Issue provides a mock function with given fields: request, receiver
func (_m *MockIssuer) Issue(request Request, receiver func(*tls.Certificate, error)) {
	_m.Called(request, receiver)
}
