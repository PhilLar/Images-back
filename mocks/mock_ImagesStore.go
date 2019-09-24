// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/PhilLar/Images-back/handlers (interfaces: ImagesStore)

// Package mocks is a generated GoMock package.
package mocks

import (
	models "github.com/PhilLar/Images-back/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockImagesStore is a mock of ImagesStore interface
type MockImagesStore struct {
	ctrl     *gomock.Controller
	recorder *MockImagesStoreMockRecorder
}

// MockImagesStoreMockRecorder is the mock recorder for MockImagesStore
type MockImagesStoreMockRecorder struct {
	mock *MockImagesStore
}

// NewMockImagesStore creates a new mock instance
func NewMockImagesStore(ctrl *gomock.Controller) *MockImagesStore {
	mock := &MockImagesStore{ctrl: ctrl}
	mock.recorder = &MockImagesStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockImagesStore) EXPECT() *MockImagesStoreMockRecorder {
	return m.recorder
}

// AllImages mocks base method
func (m *MockImagesStore) AllImages() ([]*models.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllImages")
	ret0, _ := ret[0].([]*models.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllImages indicates an expected call of AllImages
func (mr *MockImagesStoreMockRecorder) AllImages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllImages", reflect.TypeOf((*MockImagesStore)(nil).AllImages))
}

// DeleteImage mocks base method
func (m *MockImagesStore) DeleteImage(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteImage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteImage indicates an expected call of DeleteImage
func (mr *MockImagesStoreMockRecorder) DeleteImage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteImage", reflect.TypeOf((*MockImagesStore)(nil).DeleteImage), arg0)
}

// InsertImage mocks base method
func (m *MockImagesStore) InsertImage(arg0, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertImage", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertImage indicates an expected call of InsertImage
func (mr *MockImagesStoreMockRecorder) InsertImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertImage", reflect.TypeOf((*MockImagesStore)(nil).InsertImage), arg0, arg1)
}
