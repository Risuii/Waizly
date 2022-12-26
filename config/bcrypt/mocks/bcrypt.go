package mocks

import mock "github.com/stretchr/testify/mock"

// Bcrypt is an autogenerated mock type for the Bcrypt type
type Bcrypt struct {
	mock.Mock
}

// ComparePasswordHash provides a mock function with given fields: plain, hash
func (_m *Bcrypt) ComparePasswordHash(plain string, hash string) bool {
	ret := _m.Called(plain, hash)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(plain, hash)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// HashPassword provides a mock function with given fields: plain
func (_m *Bcrypt) HashPassword(plain string) (string, error) {
	ret := _m.Called(plain)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(plain)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(plain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
