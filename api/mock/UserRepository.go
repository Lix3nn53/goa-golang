// Code generated by MockGen. DO NOT EDIT.
// Source: ./app/repository/userRepository/userRepository.go

// Package mock is a generated GoMock package.
package mock

import (
	userModel "goa-golang/app/model/userModel"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserRepositoryInterface is a mock of UserRepositoryInterface interface.
type MockUserRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryInterfaceMockRecorder
}

// MockUserRepositoryInterfaceMockRecorder is the mock recorder for MockUserRepositoryInterface.
type MockUserRepositoryInterfaceMockRecorder struct {
	mock *MockUserRepositoryInterface
}

// NewMockUserRepositoryInterface creates a new mock instance.
func NewMockUserRepositoryInterface(ctrl *gomock.Controller) *MockUserRepositoryInterface {
	mock := &MockUserRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepositoryInterface) EXPECT() *MockUserRepositoryInterfaceMockRecorder {
	return m.recorder
}

// CreateUUID mocks base method.
func (m *MockUserRepositoryInterface) CreateUUID(uuid string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUUID", uuid)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUUID indicates an expected call of CreateUUID.
func (mr *MockUserRepositoryInterfaceMockRecorder) CreateUUID(uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUUID", reflect.TypeOf((*MockUserRepositoryInterface)(nil).CreateUUID), uuid)
}

// CreateWebData mocks base method.
func (m *MockUserRepositoryInterface) CreateWebData(uuid string, create userModel.CreateUser) (*userModel.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWebData", uuid, create)
	ret0, _ := ret[0].(*userModel.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWebData indicates an expected call of CreateWebData.
func (mr *MockUserRepositoryInterfaceMockRecorder) CreateWebData(uuid, create interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWebData", reflect.TypeOf((*MockUserRepositoryInterface)(nil).CreateWebData), uuid, create)
}

// FindByID mocks base method.
func (m *MockUserRepositoryInterface) FindByID(uuid string) (*userModel.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", uuid)
	ret0, _ := ret[0].(*userModel.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockUserRepositoryInterfaceMockRecorder) FindByID(uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockUserRepositoryInterface)(nil).FindByID), uuid)
}

// RemoveByID mocks base method.
func (m *MockUserRepositoryInterface) RemoveByID(uuid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveByID", uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveByID indicates an expected call of RemoveByID.
func (mr *MockUserRepositoryInterfaceMockRecorder) RemoveByID(uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveByID", reflect.TypeOf((*MockUserRepositoryInterface)(nil).RemoveByID), uuid)
}

// UpdateByID mocks base method.
func (m *MockUserRepositoryInterface) UpdateByID(uuid string, user userModel.UpdateUser) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateByID", uuid, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateByID indicates an expected call of UpdateByID.
func (mr *MockUserRepositoryInterfaceMockRecorder) UpdateByID(uuid, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateByID", reflect.TypeOf((*MockUserRepositoryInterface)(nil).UpdateByID), uuid, user)
}
