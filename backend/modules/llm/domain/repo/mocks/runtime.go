package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/kiosk404/airi-go/backend/modules/llm/domain/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockIRuntimeRepo is a mock of IRuntimeRepo interface.
type MockIRuntimeRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIRuntimeRepoMockRecorder
	isgomock struct{}
}

// MockIRuntimeRepoMockRecorder is the mock recorder for MockIRuntimeRepo.
type MockIRuntimeRepoMockRecorder struct {
	mock *MockIRuntimeRepo
}

// NewMockIRuntimeRepo creates a new mock instance.
func NewMockIRuntimeRepo(ctrl *gomock.Controller) *MockIRuntimeRepo {
	mock := &MockIRuntimeRepo{ctrl: ctrl}
	mock.recorder = &MockIRuntimeRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRuntimeRepo) EXPECT() *MockIRuntimeRepoMockRecorder {
	return m.recorder
}

// CreateModelRequestRecord mocks base method.
func (m *MockIRuntimeRepo) CreateModelRequestRecord(ctx context.Context, record *entity.ModelRequestRecord) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateModelRequestRecord", ctx, record)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateModelRequestRecord indicates an expected call of CreateModelRequestRecord.
func (mr *MockIRuntimeRepoMockRecorder) CreateModelRequestRecord(ctx, record any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateModelRequestRecord", reflect.TypeOf((*MockIRuntimeRepo)(nil).CreateModelRequestRecord), ctx, record)
}
