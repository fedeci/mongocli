// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mongodb/mongocli/internal/store (interfaces: ProcessLister)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	mongodbatlas "go.mongodb.org/atlas/mongodbatlas"
)

// MockProcessLister is a mock of ProcessLister interface.
type MockProcessLister struct {
	ctrl     *gomock.Controller
	recorder *MockProcessListerMockRecorder
}

// MockProcessListerMockRecorder is the mock recorder for MockProcessLister.
type MockProcessListerMockRecorder struct {
	mock *MockProcessLister
}

// NewMockProcessLister creates a new mock instance.
func NewMockProcessLister(ctrl *gomock.Controller) *MockProcessLister {
	mock := &MockProcessLister{ctrl: ctrl}
	mock.recorder = &MockProcessListerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProcessLister) EXPECT() *MockProcessListerMockRecorder {
	return m.recorder
}

// Processes mocks base method.
func (m *MockProcessLister) Processes(arg0 string, arg1 *mongodbatlas.ProcessesListOptions) ([]*mongodbatlas.Process, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Processes", arg0, arg1)
	ret0, _ := ret[0].([]*mongodbatlas.Process)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Processes indicates an expected call of Processes.
func (mr *MockProcessListerMockRecorder) Processes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Processes", reflect.TypeOf((*MockProcessLister)(nil).Processes), arg0, arg1)
}
