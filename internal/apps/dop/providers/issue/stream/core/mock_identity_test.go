// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/core/user/identity/provider.go

// Package core is a generated GoMock package.
package core

import (
	reflect "reflect"

	apistructs "github.com/erda-project/erda/apistructs"
	identity "github.com/erda-project/erda/internal/core/user/common"
	gomock "github.com/golang/mock/gomock"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// FindUsers mocks base method.
func (m *MockInterface) FindUsers(ids []string) ([]identity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUsers", ids)
	ret0, _ := ret[0].([]identity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUsers indicates an expected call of FindUsers.
func (mr *MockInterfaceMockRecorder) FindUsers(ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUsers", reflect.TypeOf((*MockInterface)(nil).FindUsers), ids)
}

// FindUsersByKey mocks base method.
func (m *MockInterface) FindUsersByKey(key string) ([]identity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUsersByKey", key)
	ret0, _ := ret[0].([]identity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUsersByKey indicates an expected call of FindUsersByKey.
func (mr *MockInterfaceMockRecorder) FindUsersByKey(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUsersByKey", reflect.TypeOf((*MockInterface)(nil).FindUsersByKey), key)
}

// GetUser mocks base method.
func (m *MockInterface) GetUser(userID string) (*identity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", userID)
	ret0, _ := ret[0].(*identity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockInterfaceMockRecorder) GetUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockInterface)(nil).GetUser), userID)
}

// GetUsers mocks base method.
func (m *MockInterface) GetUsers(IDs []string, needDesensitize bool) (map[string]apistructs.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", IDs, needDesensitize)
	ret0, _ := ret[0].(map[string]apistructs.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockInterfaceMockRecorder) GetUsers(IDs, needDesensitize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockInterface)(nil).GetUsers), IDs, needDesensitize)
}
