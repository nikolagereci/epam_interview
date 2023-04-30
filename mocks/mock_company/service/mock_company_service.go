// Code generated by MockGen. DO NOT EDIT.
// Source: ../company/company_service.go

// Package mock_company_service is a generated GoMock package.
package mock_company_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	model "github.com/ngereci/xm_interview/model"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateCompany mocks base method.
func (m *MockService) CreateCompany(ctx context.Context, newCompany *model.Company) (*model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCompany", ctx, newCompany)
	ret0, _ := ret[0].(*model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCompany indicates an expected call of CreateCompany.
func (mr *MockServiceMockRecorder) CreateCompany(ctx, newCompany interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCompany", reflect.TypeOf((*MockService)(nil).CreateCompany), ctx, newCompany)
}

// DeleteCompany mocks base method.
func (m *MockService) DeleteCompany(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCompany", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCompany indicates an expected call of DeleteCompany.
func (mr *MockServiceMockRecorder) DeleteCompany(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCompany", reflect.TypeOf((*MockService)(nil).DeleteCompany), ctx, id)
}

// GetCompanyByID mocks base method.
func (m *MockService) GetCompanyByID(ctx context.Context, id uuid.UUID) (*model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompanyByID", ctx, id)
	ret0, _ := ret[0].(*model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompanyByID indicates an expected call of GetCompanyByID.
func (mr *MockServiceMockRecorder) GetCompanyByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompanyByID", reflect.TypeOf((*MockService)(nil).GetCompanyByID), ctx, id)
}

// UpdateCompany mocks base method.
func (m *MockService) UpdateCompany(ctx context.Context, id uuid.UUID, forUpdateCompany *model.Company) (*model.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCompany", ctx, id, forUpdateCompany)
	ret0, _ := ret[0].(*model.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCompany indicates an expected call of UpdateCompany.
func (mr *MockServiceMockRecorder) UpdateCompany(ctx, id, forUpdateCompany interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCompany", reflect.TypeOf((*MockService)(nil).UpdateCompany), ctx, id, forUpdateCompany)
}