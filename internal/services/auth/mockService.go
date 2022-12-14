// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/calebtracey/phish-stats-api/internal/services/auth (interfaces: ServiceI)

// Package auth is a generated GoMock package.
package auth

import (
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

// GenerateAllTokens mocks base method.
func (m *MockServiceI) GenerateAllTokens(arg0 models.User) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAllTokens", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateAllTokens indicates an expected call of GenerateAllTokens.
func (mr *MockServiceIMockRecorder) GenerateAllTokens(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAllTokens", reflect.TypeOf((*MockServiceI)(nil).GenerateAllTokens), arg0)
}

// HashPassword mocks base method.
func (m *MockServiceI) HashPassword(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashPassword", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// HashPassword indicates an expected call of HashPassword.
func (mr *MockServiceIMockRecorder) HashPassword(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashPassword", reflect.TypeOf((*MockServiceI)(nil).HashPassword), arg0)
}

// ValidateToken mocks base method.
func (m *MockServiceI) ValidateToken(arg0 string) (*SignedDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", arg0)
	ret0, _ := ret[0].(*SignedDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockServiceIMockRecorder) ValidateToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockServiceI)(nil).ValidateToken), arg0)
}

// VerifyPassword mocks base method.
func (m *MockServiceI) VerifyPassword(arg0, arg1 string) (bool, string) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyPassword", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	return ret0, ret1
}

// VerifyPassword indicates an expected call of VerifyPassword.
func (mr *MockServiceIMockRecorder) VerifyPassword(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyPassword", reflect.TypeOf((*MockServiceI)(nil).VerifyPassword), arg0, arg1)
}
