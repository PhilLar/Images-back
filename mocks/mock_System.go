// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/PhilLar/Images-back/models (interfaces: System)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSystem is a mock of System interface
type MockSystem struct {
	ctrl     *gomock.Controller
	recorder *MockSystemMockRecorder
}

// MockSystemMockRecorder is the mock recorder for MockSystem
type MockSystemMockRecorder struct {
	mock *MockSystem
}

// NewMockSystem creates a new mock instance
func NewMockSystem(ctrl *gomock.Controller) *MockSystem {
	mock := &MockSystem{ctrl: ctrl}
	mock.recorder = &MockSystemMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSystem) EXPECT() *MockSystemMockRecorder {
	return m.recorder
}

// Remove mocks base method
func (m *MockSystem) Remove(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockSystemMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockSystem)(nil).Remove), arg0)
}
