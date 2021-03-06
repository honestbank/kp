// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/honestbank/kp (interfaces: KPProducer)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	kp "github.com/honestbank/kp"
)

// MockKPProducer is a mock of KPProducer interface.
type MockKPProducer struct {
	ctrl     *gomock.Controller
	recorder *MockKPProducerMockRecorder
}

// MockKPProducerMockRecorder is the mock recorder for MockKPProducer.
type MockKPProducerMockRecorder struct {
	mock *MockKPProducer
}

// NewMockKPProducer creates a new mock instance.
func NewMockKPProducer(ctrl *gomock.Controller) *MockKPProducer {
	mock := &MockKPProducer{ctrl: ctrl}
	mock.recorder = &MockKPProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKPProducer) EXPECT() *MockKPProducerMockRecorder {
	return m.recorder
}

// GetProducer mocks base method.
func (m *MockKPProducer) GetProducer(arg0 kp.KafkaConfig) *kp.Producer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducer", arg0)
	ret0, _ := ret[0].(*kp.Producer)
	return ret0
}

// GetProducer indicates an expected call of GetProducer.
func (mr *MockKPProducerMockRecorder) GetProducer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducer", reflect.TypeOf((*MockKPProducer)(nil).GetProducer), arg0)
}

// ProduceMessage mocks base method.
func (m *MockKPProducer) ProduceMessage(arg0 context.Context, arg1, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProduceMessage", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProduceMessage indicates an expected call of ProduceMessage.
func (mr *MockKPProducerMockRecorder) ProduceMessage(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProduceMessage", reflect.TypeOf((*MockKPProducer)(nil).ProduceMessage), arg0, arg1, arg2, arg3)
}
