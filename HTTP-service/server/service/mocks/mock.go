// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	"serverdb/models"
)

// Mockrepository is a mock of repository interface.
type Mockrepository struct {
	ctrl     *gomock.Controller
	recorder *MockrepositoryMockRecorder
}

// MockrepositoryMockRecorder is the mock recorder for Mockrepository.
type MockrepositoryMockRecorder struct {
	mock *Mockrepository
}

// NewMockrepository creates a new mock instance.
func NewMockrepository(ctrl *gomock.Controller) *Mockrepository {
	mock := &Mockrepository{ctrl: ctrl}
	mock.recorder = &MockrepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrepository) EXPECT() *MockrepositoryMockRecorder {
	return m.recorder
}

// FindUser mocks base method.
func (m *Mockrepository) FindUser(userID string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUser", userID)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUser indicates an expected call of FindUser.
func (mr *MockrepositoryMockRecorder) FindUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*Mockrepository)(nil).FindUser), userID)
}

// Insert mocks base method.
func (m *Mockrepository) Insert(user models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockrepositoryMockRecorder) Insert(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*Mockrepository)(nil).Insert), user)
}

// RemoveAll mocks base method.
func (m *Mockrepository) RemoveAll(userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAll indicates an expected call of RemoveAll.
func (mr *MockrepositoryMockRecorder) RemoveAll(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*Mockrepository)(nil).RemoveAll), userID)
}

// Update mocks base method.
func (m *Mockrepository) Update(user, order, object string, objectChange interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", user, order, object, objectChange)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockrepositoryMockRecorder) Update(user, order, object, objectChange interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*Mockrepository)(nil).Update), user, order, object, objectChange)
}
