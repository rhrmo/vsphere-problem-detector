// Code generated by MockGen. DO NOT EDIT.
// Source: ./permissions.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	types "github.com/vmware/govmomi/vim25/types"
	reflect "reflect"
)

// MockAuthManager is a mock of AuthManager interface
type MockAuthManager struct {
	ctrl     *gomock.Controller
	recorder *MockAuthManagerMockRecorder
}

// MockAuthManagerMockRecorder is the mock recorder for MockAuthManager
type MockAuthManagerMockRecorder struct {
	mock *MockAuthManager
}

// NewMockAuthManager creates a new mock instance
func NewMockAuthManager(ctrl *gomock.Controller) *MockAuthManager {
	mock := &MockAuthManager{ctrl: ctrl}
	mock.recorder = &MockAuthManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthManager) EXPECT() *MockAuthManagerMockRecorder {
	return m.recorder
}

// FetchUserPrivilegeOnEntities mocks base method
func (m *MockAuthManager) FetchUserPrivilegeOnEntities(ctx context.Context, entities []types.ManagedObjectReference, userName string) ([]types.UserPrivilegeResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchUserPrivilegeOnEntities", ctx, entities, userName)
	ret0, _ := ret[0].([]types.UserPrivilegeResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchUserPrivilegeOnEntities indicates an expected call of FetchUserPrivilegeOnEntities
func (mr *MockAuthManagerMockRecorder) FetchUserPrivilegeOnEntities(ctx, entities, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchUserPrivilegeOnEntities", reflect.TypeOf((*MockAuthManager)(nil).FetchUserPrivilegeOnEntities), ctx, entities, userName)
}
