// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/devices/devicemanagerregistrator.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	devices "github.com/yndd/ztp-dhcp/pkg/devices"
)

// MockDeviceManagerRegistrator is a mock of DeviceManagerRegistrator interface.
type MockDeviceManagerRegistrator struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceManagerRegistratorMockRecorder
}

// MockDeviceManagerRegistratorMockRecorder is the mock recorder for MockDeviceManagerRegistrator.
type MockDeviceManagerRegistratorMockRecorder struct {
	mock *MockDeviceManagerRegistrator
}

// NewMockDeviceManagerRegistrator creates a new mock instance.
func NewMockDeviceManagerRegistrator(ctrl *gomock.Controller) *MockDeviceManagerRegistrator {
	mock := &MockDeviceManagerRegistrator{ctrl: ctrl}
	mock.recorder = &MockDeviceManagerRegistratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeviceManagerRegistrator) EXPECT() *MockDeviceManagerRegistratorMockRecorder {
	return m.recorder
}

// RegisterDevice mocks base method.
func (m *MockDeviceManagerRegistrator) RegisterDevice(arg0 []string, arg1 devices.Device) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterDevice", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterDevice indicates an expected call of RegisterDevice.
func (mr *MockDeviceManagerRegistratorMockRecorder) RegisterDevice(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterDevice", reflect.TypeOf((*MockDeviceManagerRegistrator)(nil).RegisterDevice), arg0, arg1)
}
