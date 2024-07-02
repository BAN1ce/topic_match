// Code generated by MockGen. DO NOT EDIT.
// Source: .//broker/monitor/monitor.go

// Package monitor is a generated GoMock package.
package monitor

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockClientKeepAliveMonitor is a mock of ClientKeepAliveMonitor interface.
type MockClientKeepAliveMonitor struct {
	ctrl     *gomock.Controller
	recorder *MockClientKeepAliveMonitorMockRecorder
}

// MockClientKeepAliveMonitorMockRecorder is the mock recorder for MockClientKeepAliveMonitor.
type MockClientKeepAliveMonitorMockRecorder struct {
	mock *MockClientKeepAliveMonitor
}

// NewMockClientKeepAliveMonitor creates a new mock instance.
func NewMockClientKeepAliveMonitor(ctrl *gomock.Controller) *MockClientKeepAliveMonitor {
	mock := &MockClientKeepAliveMonitor{ctrl: ctrl}
	mock.recorder = &MockClientKeepAliveMonitorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientKeepAliveMonitor) EXPECT() *MockClientKeepAliveMonitorMockRecorder {
	return m.recorder
}

// SetClientAliveTime mocks base method.
func (m *MockClientKeepAliveMonitor) SetClientAliveTime(clientID string, t *time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetClientAliveTime", clientID, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetClientAliveTime indicates an expected call of SetClientAliveTime.
func (mr *MockClientKeepAliveMonitorMockRecorder) SetClientAliveTime(clientID, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetClientAliveTime", reflect.TypeOf((*MockClientKeepAliveMonitor)(nil).SetClientAliveTime), clientID, t)
}
