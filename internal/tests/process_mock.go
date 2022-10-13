// Code generated by mockery v2.14.0. DO NOT EDIT.

package tests

import mock "github.com/stretchr/testify/mock"

// Process is an autogenerated mock type for the Process type
type Process struct {
	mock.Mock
}

// Execute provides a mock function with given fields: name, args
func (_m *Process) Execute(name string, args ...[]string) ([]byte, error) {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, ...[]string) []byte); ok {
		r0 = rf(name, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...[]string) error); ok {
		r1 = rf(name, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewProcess interface {
	mock.TestingT
	Cleanup(func())
}

// NewProcess creates a new instance of Process. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProcess(t mockConstructorTestingTNewProcess) *Process {
	mock := &Process{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
