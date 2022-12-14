// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/calebtracey/phish-stats-api/internal/services/psql (interfaces: ServiceI)

// Package psql is a generated GoMock package.
package psql

import (
	context "context"
	reflect "reflect"

	models "github.com/calebtracey/phish-stats-api/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockServiceI is a mock of ServiceI interface.
type MockServiceI struct {
	ctrl     *gomock.Controller
	recorder *MockServiceIMockRecorder
}

// MockServiceIMockRecorder is the mock recorder for MockServiceI.
type MockServiceIMockRecorder struct {
	mock *MockServiceI
}

// NewMockServiceI creates a new mock instance.
func NewMockServiceI(ctrl *gomock.Controller) *MockServiceI {
	mock := &MockServiceI{ctrl: ctrl}
	mock.recorder = &MockServiceIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceI) EXPECT() *MockServiceIMockRecorder {
	return m.recorder
}

// FindUser mocks base method.
func (m *MockServiceI) FindUser(arg0 context.Context, arg1 string) (*models.UserParsedResponse, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUser", arg0, arg1)
	ret0, _ := ret[0].(*models.UserParsedResponse)
	ret1, _ := ret[1].([]error)
	return ret0, ret1
}

// FindUser indicates an expected call of FindUser.
func (mr *MockServiceIMockRecorder) FindUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*MockServiceI)(nil).FindUser), arg0, arg1)
}

// InsertNewUser mocks base method.
func (m *MockServiceI) InsertNewUser(arg0 context.Context, arg1 string) (*models.NewUserResponse, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNewUser", arg0, arg1)
	ret0, _ := ret[0].(*models.NewUserResponse)
	ret1, _ := ret[1].([]error)
	return ret0, ret1
}

// InsertNewUser indicates an expected call of InsertNewUser.
func (mr *MockServiceIMockRecorder) InsertNewUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNewUser", reflect.TypeOf((*MockServiceI)(nil).InsertNewUser), arg0, arg1)
}

// InsertOne mocks base method.
func (m *MockServiceI) InsertOne(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOne", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOne indicates an expected call of InsertOne.
func (mr *MockServiceIMockRecorder) InsertOne(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOne", reflect.TypeOf((*MockServiceI)(nil).InsertOne), arg0, arg1)
}
