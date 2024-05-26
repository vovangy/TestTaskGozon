// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_users is a generated GoMock package.
package mock_users

import (
	context "context"
	models "myHabr/internal/models"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockUserUsecase is a mock of UserUsecase interface.
type MockUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseMockRecorder
}

// MockUserUsecaseMockRecorder is the mock recorder for MockUserUsecase.
type MockUserUsecaseMockRecorder struct {
	mock *MockUserUsecase
}

// NewMockUserUsecase creates a new mock instance.
func NewMockUserUsecase(ctrl *gomock.Controller) *MockUserUsecase {
	mock := &MockUserUsecase{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUsecase) EXPECT() *MockUserUsecaseMockRecorder {
	return m.recorder
}

// GetUsernameById mocks base method.
func (m *MockUserUsecase) GetUsernameById(ctx context.Context, id int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsernameById", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsernameById indicates an expected call of GetUsernameById.
func (mr *MockUserUsecaseMockRecorder) GetUsernameById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsernameById", reflect.TypeOf((*MockUserUsecase)(nil).GetUsernameById), ctx, id)
}

// Login mocks base method.
func (m *MockUserUsecase) Login(ctx context.Context, data *models.UserSignInUp) (*models.UserCreatedInfo, string, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, data)
	ret0, _ := ret[0].(*models.UserCreatedInfo)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(time.Time)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// Login indicates an expected call of Login.
func (mr *MockUserUsecaseMockRecorder) Login(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserUsecase)(nil).Login), ctx, data)
}

// SignUp mocks base method.
func (m *MockUserUsecase) SignUp(ctx context.Context, data *models.UserSignInUp) (*models.UserCreatedInfo, string, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", ctx, data)
	ret0, _ := ret[0].(*models.UserCreatedInfo)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(time.Time)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUserUsecaseMockRecorder) SignUp(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUserUsecase)(nil).SignUp), ctx, data)
}

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// CheckUser mocks base method.
func (m *MockUserRepo) CheckUser(ctx context.Context, data *models.UserSignInUp) (*models.UserCreatedInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUser", ctx, data)
	ret0, _ := ret[0].(*models.UserCreatedInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockUserRepoMockRecorder) CheckUser(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockUserRepo)(nil).CheckUser), ctx, data)
}

// CreateUser mocks base method.
func (m *MockUserRepo) CreateUser(ctx context.Context, user *models.UserSignInUp) (*models.UserCreatedInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(*models.UserCreatedInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepoMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepo)(nil).CreateUser), ctx, user)
}

// GetUserByLogin mocks base method.
func (m *MockUserRepo) GetUserByLogin(ctx context.Context, username string) (*models.UserCreatedInfo, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", ctx, username)
	ret0, _ := ret[0].(*models.UserCreatedInfo)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockUserRepoMockRecorder) GetUserByLogin(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockUserRepo)(nil).GetUserByLogin), ctx, username)
}

// GetUsernameById mocks base method.
func (m *MockUserRepo) GetUsernameById(ctx context.Context, id int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsernameById", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsernameById indicates an expected call of GetUsernameById.
func (mr *MockUserRepoMockRecorder) GetUsernameById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsernameById", reflect.TypeOf((*MockUserRepo)(nil).GetUsernameById), ctx, id)
}
