// Code generated by MockGen. DO NOT EDIT.
// Source: .//broker/sub_center.go

// Package broker is a generated GoMock package.
package broker

import (
	reflect "reflect"

	packets "github.com/eclipse/paho.golang/packets"
	gomock "github.com/golang/mock/gomock"
)

// MockSubClient is a mock of SubClient interface.
type MockSubClient struct {
	ctrl     *gomock.Controller
	recorder *MockSubClientMockRecorder
}

// MockSubClientMockRecorder is the mock recorder for MockSubClient.
type MockSubClientMockRecorder struct {
	mock *MockSubClient
}

// NewMockSubClient creates a new mock instance.
func NewMockSubClient(ctrl *gomock.Controller) *MockSubClient {
	mock := &MockSubClient{ctrl: ctrl}
	mock.recorder = &MockSubClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubClient) EXPECT() *MockSubClientMockRecorder {
	return m.recorder
}

// GetClientID mocks base method.
func (m *MockSubClient) GetClientID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClientID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetClientID indicates an expected call of GetClientID.
func (mr *MockSubClientMockRecorder) GetClientID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClientID", reflect.TypeOf((*MockSubClient)(nil).GetClientID))
}

// GetQoS mocks base method.
func (m *MockSubClient) GetQoS() int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQoS")
	ret0, _ := ret[0].(int32)
	return ret0
}

// GetQoS indicates an expected call of GetQoS.
func (mr *MockSubClientMockRecorder) GetQoS() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQoS", reflect.TypeOf((*MockSubClient)(nil).GetQoS))
}

// MockSubCenter is a mock of SubCenter interface.
type MockSubCenter struct {
	ctrl     *gomock.Controller
	recorder *MockSubCenterMockRecorder
}

// MockSubCenterMockRecorder is the mock recorder for MockSubCenter.
type MockSubCenterMockRecorder struct {
	mock *MockSubCenter
}

// NewMockSubCenter creates a new mock instance.
func NewMockSubCenter(ctrl *gomock.Controller) *MockSubCenter {
	mock := &MockSubCenter{ctrl: ctrl}
	mock.recorder = &MockSubCenterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubCenter) EXPECT() *MockSubCenterMockRecorder {
	return m.recorder
}

// CreateSub mocks base method.
func (m *MockSubCenter) CreateSub(clientID string, topics []packets.SubOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSub", clientID, topics)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSub indicates an expected call of CreateSub.
func (mr *MockSubCenterMockRecorder) CreateSub(clientID, topics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSub", reflect.TypeOf((*MockSubCenter)(nil).CreateSub), clientID, topics)
}

// DeleteClient mocks base method.
func (m *MockSubCenter) DeleteClient(clientID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteClient", clientID)
}

// DeleteClient indicates an expected call of DeleteClient.
func (mr *MockSubCenterMockRecorder) DeleteClient(clientID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteClient", reflect.TypeOf((*MockSubCenter)(nil).DeleteClient), clientID)
}

// DeleteSub mocks base method.
func (m *MockSubCenter) DeleteSub(clientID string, topics []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSub", clientID, topics)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSub indicates an expected call of DeleteSub.
func (mr *MockSubCenterMockRecorder) DeleteSub(clientID, topics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSub", reflect.TypeOf((*MockSubCenter)(nil).DeleteSub), clientID, topics)
}

// Match mocks base method.
func (m *MockSubCenter) Match(topic string) map[string]int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Match", topic)
	ret0, _ := ret[0].(map[string]int32)
	return ret0
}

// Match indicates an expected call of Match.
func (mr *MockSubCenterMockRecorder) Match(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Match", reflect.TypeOf((*MockSubCenter)(nil).Match), topic)
}

// MatchTopic mocks base method.
func (m *MockSubCenter) MatchTopic(topic string) map[string]int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MatchTopic", topic)
	ret0, _ := ret[0].(map[string]int32)
	return ret0
}

// MatchTopic indicates an expected call of MatchTopic.
func (mr *MockSubCenterMockRecorder) MatchTopic(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MatchTopic", reflect.TypeOf((*MockSubCenter)(nil).MatchTopic), topic)
}
