// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// AdminStoreInterface is an autogenerated mock type for the AdminStoreInterface type
type AdminStoreInterface struct {
	mock.Mock
}

type AdminStoreInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *AdminStoreInterface) EXPECT() *AdminStoreInterface_Expecter {
	return &AdminStoreInterface_Expecter{mock: &_m.Mock}
}

// GetBestRTPProxy provides a mock function with given fields:
func (_m *AdminStoreInterface) GetBestRTPProxy() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]byte, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AdminStoreInterface_GetBestRTPProxy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBestRTPProxy'
type AdminStoreInterface_GetBestRTPProxy_Call struct {
	*mock.Call
}

// GetBestRTPProxy is a helper method to define mock.On call
func (_e *AdminStoreInterface_Expecter) GetBestRTPProxy() *AdminStoreInterface_GetBestRTPProxy_Call {
	return &AdminStoreInterface_GetBestRTPProxy_Call{Call: _e.mock.On("GetBestRTPProxy")}
}

func (_c *AdminStoreInterface_GetBestRTPProxy_Call) Run(run func()) *AdminStoreInterface_GetBestRTPProxy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AdminStoreInterface_GetBestRTPProxy_Call) Return(_a0 []byte, _a1 error) *AdminStoreInterface_GetBestRTPProxy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AdminStoreInterface_GetBestRTPProxy_Call) RunAndReturn(run func() ([]byte, error)) *AdminStoreInterface_GetBestRTPProxy_Call {
	_c.Call.Return(run)
	return _c
}

// Healthz provides a mock function with given fields:
func (_m *AdminStoreInterface) Healthz() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AdminStoreInterface_Healthz_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Healthz'
type AdminStoreInterface_Healthz_Call struct {
	*mock.Call
}

// Healthz is a helper method to define mock.On call
func (_e *AdminStoreInterface_Expecter) Healthz() *AdminStoreInterface_Healthz_Call {
	return &AdminStoreInterface_Healthz_Call{Call: _e.mock.On("Healthz")}
}

func (_c *AdminStoreInterface_Healthz_Call) Run(run func()) *AdminStoreInterface_Healthz_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AdminStoreInterface_Healthz_Call) Return(_a0 error) *AdminStoreInterface_Healthz_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AdminStoreInterface_Healthz_Call) RunAndReturn(run func() error) *AdminStoreInterface_Healthz_Call {
	_c.Call.Return(run)
	return _c
}

// NewAdminStoreInterface creates a new instance of AdminStoreInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAdminStoreInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *AdminStoreInterface {
	mock := &AdminStoreInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
