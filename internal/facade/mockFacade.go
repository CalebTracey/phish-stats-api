// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/calebtracey/phish-stats-api/internal/facade (interfaces: ServiceI)

// Package facade is a generated GoMock package.
package facade

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

// AddUserShow mocks base method.
func (m *MockServiceI) AddUserShow(arg0 context.Context, arg1 models.AddUserShowRequest) models.AddShowResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserShow", arg0, arg1)
	ret0, _ := ret[0].(models.AddShowResponse)
	return ret0
}

// AddUserShow indicates an expected call of AddUserShow.
func (mr *MockServiceIMockRecorder) AddUserShow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserShow", reflect.TypeOf((*MockServiceI)(nil).AddUserShow), arg0, arg1)
}

// GetShow mocks base method.
func (m *MockServiceI) GetShow(arg0 context.Context, arg1 models.GetShowRequest) models.ShowResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShow", arg0, arg1)
	ret0, _ := ret[0].(models.ShowResponse)
	return ret0
}

// GetShow indicates an expected call of GetShow.
func (mr *MockServiceIMockRecorder) GetShow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShow", reflect.TypeOf((*MockServiceI)(nil).GetShow), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockServiceI) GetUser(arg0 context.Context, arg1 string) models.UserResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(models.UserResponse)
	return ret0
}

// GetUser indicates an expected call of GetUser.
func (mr *MockServiceIMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockServiceI)(nil).GetUser), arg0, arg1)
}

// LoginUser mocks base method.
func (m *MockServiceI) LoginUser(arg0 context.Context, arg1 models.User) models.UserResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", arg0, arg1)
	ret0, _ := ret[0].(models.UserResponse)
	return ret0
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockServiceIMockRecorder) LoginUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockServiceI)(nil).LoginUser), arg0, arg1)
}

// RegisterUser mocks base method.
func (m *MockServiceI) RegisterUser(arg0 context.Context, arg1 models.User) models.UserResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", arg0, arg1)
	ret0, _ := ret[0].(models.UserResponse)
	return ret0
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockServiceIMockRecorder) RegisterUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockServiceI)(nil).RegisterUser), arg0, arg1)
}
