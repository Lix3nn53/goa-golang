// Code generated by MockGen. DO NOT EDIT.
// Source: UserService.go

// Package mock is a generated GoMock package.
package mock

import "goa-golang/app/model/billingModel"

/*
// MockAuthServiceUseCase is a mock of UserUseCase interface
type MockPaymentAdapaterUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentAdapterCaseMockRecorder
}
// MockAuthUseCaseMockRecorder is the mock recorder for MockAuthServiceUseCase
type MockPaymentAdapterCaseMockRecorder struct {
	mock *MockPaymentAdapaterUseCase
}
// NewMockUserServiceCase creates a new mock instance
func NewMockPaymentAdapterCase(ctrl *gomock.Controller) *MockPaymentAdapaterUseCase {
	mock := &MockPaymentAdapaterUseCase{ctrl: ctrl}
	mock.recorder = &MockPaymentAdapterCaseMockRecorder{mock}
	return mock
}
// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPaymentAdapaterUseCase) EXPECT() *MockPaymentAdapterCaseMockRecorder {
	return m.recorder
}
// Register mocks base method
func (m *MockPaymentAdapaterUseCase) CreateCustomer(params billingModel.CustomerParams) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCustomer", params)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}
// Register indicates an expected call of Register
func (mr *MockPaymentAdapterCaseMockRecorder) CreateCustomer(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCustomer", reflect.TypeOf((*MockPaymentAdapaterUseCase)(nil).CreateCustomer), params)
}
*/
// StripeAdapter adapts Stripe API
type FakeAdapter struct {
	Payment *FakePayment
}

type FakePayment struct{}

func (r *FakePayment) CreateCustomer(params billingModel.CustomerParams) (string, error) {
	return "fake", nil
}

// Pay from email to email this amount
func (b *FakeAdapter) CreateCustomer(params billingModel.CustomerParams) (string, error) {
	return b.Payment.CreateCustomer(params)
}
