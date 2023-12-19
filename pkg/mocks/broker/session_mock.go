// Code generated by MockGen. DO NOT EDIT.
// Source: .//broker/session.go

// Package broker is a generated GoMock package.
package broker

import (
	broker "github.com/BAN1ce/skyTree/pkg/broker/session"
	reflect "reflect"

	proto "github.com/BAN1ce/Tree/proto"
	gomock "github.com/golang/mock/gomock"
)

// MockSession is a mock of Session interface.
type MockSession struct {
	ctrl     *gomock.Controller
	recorder *MockSessionMockRecorder
}

// MockSessionMockRecorder is the mock recorder for MockSession.
type MockSessionMockRecorder struct {
	mock *MockSession
}

// NewMockSession creates a new mock instance.
func NewMockSession(ctrl *gomock.Controller) *MockSession {
	mock := &MockSession{ctrl: ctrl}
	mock.recorder = &MockSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSession) EXPECT() *MockSessionMockRecorder {
	return m.recorder
}

// CreateSubTopic mocks base method.
func (m *MockSession) CreateSubTopic(topic string, option *proto.SubOption) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateSubTopic", topic, option)
}

// CreateSubTopic indicates an expected call of CreateSubTopic.
func (mr *MockSessionMockRecorder) CreateSubTopic(topic, option interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubTopic", reflect.TypeOf((*MockSession)(nil).CreateSubTopic), topic, option)
}

// CreateTopicUnFinishedMessage mocks base method.
func (m *MockSession) CreateTopicUnFinishedMessage(topic string, message []broker.UnFinishedMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateTopicUnFinishedMessage", topic, message)
}

// CreateTopicUnFinishedMessage indicates an expected call of CreateTopicUnFinishedMessage.
func (mr *MockSessionMockRecorder) CreateTopicUnFinishedMessage(topic, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTopicUnFinishedMessage", reflect.TypeOf((*MockSession)(nil).CreateTopicUnFinishedMessage), topic, message)
}

// DeleteSubTopic mocks base method.
func (m *MockSession) DeleteSubTopic(topic string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteSubTopic", topic)
}

// DeleteSubTopic indicates an expected call of DeleteSubTopic.
func (mr *MockSessionMockRecorder) DeleteSubTopic(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubTopic", reflect.TypeOf((*MockSession)(nil).DeleteSubTopic), topic)
}

// DeleteTopicLatestPushedMessageID mocks base method.
func (m *MockSession) DeleteTopicLatestPushedMessageID(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteTopicLatestPushedMessageID", topic, messageID)
}

// DeleteTopicLatestPushedMessageID indicates an expected call of DeleteTopicLatestPushedMessageID.
func (mr *MockSessionMockRecorder) DeleteTopicLatestPushedMessageID(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTopicLatestPushedMessageID", reflect.TypeOf((*MockSession)(nil).DeleteTopicLatestPushedMessageID), topic, messageID)
}

// DeleteTopicUnFinishedMessage mocks base method.
func (m *MockSession) DeleteTopicUnFinishedMessage(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteTopicUnFinishedMessage", topic, messageID)
}

// DeleteTopicUnFinishedMessage indicates an expected call of DeleteTopicUnFinishedMessage.
func (mr *MockSessionMockRecorder) DeleteTopicUnFinishedMessage(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTopicUnFinishedMessage", reflect.TypeOf((*MockSession)(nil).DeleteTopicUnFinishedMessage), topic, messageID)
}

// GetConnectProperties mocks base method.
func (m *MockSession) GetConnectProperties() (*broker.ConnectProperties, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectProperties")
	ret0, _ := ret[0].(*broker.ConnectProperties)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectProperties indicates an expected call of GetConnectProperties.
func (mr *MockSessionMockRecorder) GetConnectProperties() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectProperties", reflect.TypeOf((*MockSession)(nil).GetConnectProperties))
}

// GetWillMessage mocks base method.
func (m *MockSession) GetWillMessage() (*broker.WillMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWillMessage")
	ret0, _ := ret[0].(*broker.WillMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWillMessage indicates an expected call of GetWillMessage.
func (mr *MockSessionMockRecorder) GetWillMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWillMessage", reflect.TypeOf((*MockSession)(nil).GetWillMessage))
}

// ReadSubTopics mocks base method.
func (m *MockSession) ReadSubTopics() map[string]*proto.SubOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSubTopics")
	ret0, _ := ret[0].(map[string]*proto.SubOption)
	return ret0
}

// ReadSubTopics indicates an expected call of ReadSubTopics.
func (mr *MockSessionMockRecorder) ReadSubTopics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSubTopics", reflect.TypeOf((*MockSession)(nil).ReadSubTopics))
}

// ReadTopicLatestPushedMessageID mocks base method.
func (m *MockSession) ReadTopicLatestPushedMessageID(topic string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTopicLatestPushedMessageID", topic)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ReadTopicLatestPushedMessageID indicates an expected call of ReadTopicLatestPushedMessageID.
func (mr *MockSessionMockRecorder) ReadTopicLatestPushedMessageID(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTopicLatestPushedMessageID", reflect.TypeOf((*MockSession)(nil).ReadTopicLatestPushedMessageID), topic)
}

// ReadTopicUnFinishedMessage mocks base method.
func (m *MockSession) ReadTopicUnFinishedMessage(topic string) []broker.UnFinishedMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTopicUnFinishedMessage", topic)
	ret0, _ := ret[0].([]broker.UnFinishedMessage)
	return ret0
}

// ReadTopicUnFinishedMessage indicates an expected call of ReadTopicUnFinishedMessage.
func (mr *MockSessionMockRecorder) ReadTopicUnFinishedMessage(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTopicUnFinishedMessage", reflect.TypeOf((*MockSession)(nil).ReadTopicUnFinishedMessage), topic)
}

// Release mocks base method.
func (m *MockSession) Release() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Release")
}

// Release indicates an expected call of Release.
func (mr *MockSessionMockRecorder) Release() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Release", reflect.TypeOf((*MockSession)(nil).Release))
}

// SetConnectProperties mocks base method.
func (m *MockSession) SetConnectProperties(properties *broker.ConnectProperties) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetConnectProperties", properties)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetConnectProperties indicates an expected call of SetConnectProperties.
func (mr *MockSessionMockRecorder) SetConnectProperties(properties interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConnectProperties", reflect.TypeOf((*MockSession)(nil).SetConnectProperties), properties)
}

// SetTopicLatestPushedMessageID mocks base method.
func (m *MockSession) SetTopicLatestPushedMessageID(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTopicLatestPushedMessageID", topic, messageID)
}

// SetTopicLatestPushedMessageID indicates an expected call of SetTopicLatestPushedMessageID.
func (mr *MockSessionMockRecorder) SetTopicLatestPushedMessageID(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTopicLatestPushedMessageID", reflect.TypeOf((*MockSession)(nil).SetTopicLatestPushedMessageID), topic, messageID)
}

// SetWillMessage mocks base method.
func (m *MockSession) SetWillMessage(message *broker.WillMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWillMessage", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWillMessage indicates an expected call of SetWillMessage.
func (mr *MockSessionMockRecorder) SetWillMessage(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWillMessage", reflect.TypeOf((*MockSession)(nil).SetWillMessage), message)
}

// MockSessionTopic is a mock of TopicManager interface.
type MockSessionTopic struct {
	ctrl     *gomock.Controller
	recorder *MockSessionTopicMockRecorder
}

// MockSessionTopicMockRecorder is the mock recorder for MockSessionTopic.
type MockSessionTopicMockRecorder struct {
	mock *MockSessionTopic
}

// NewMockSessionTopic creates a new mock instance.
func NewMockSessionTopic(ctrl *gomock.Controller) *MockSessionTopic {
	mock := &MockSessionTopic{ctrl: ctrl}
	mock.recorder = &MockSessionTopicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionTopic) EXPECT() *MockSessionTopicMockRecorder {
	return m.recorder
}

// CreateSubTopic mocks base method.
func (m *MockSessionTopic) CreateSubTopic(topic string, option *proto.SubOption) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateSubTopic", topic, option)
}

// CreateSubTopic indicates an expected call of CreateSubTopic.
func (mr *MockSessionTopicMockRecorder) CreateSubTopic(topic, option interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubTopic", reflect.TypeOf((*MockSessionTopic)(nil).CreateSubTopic), topic, option)
}

// CreateTopicUnFinishedMessage mocks base method.
func (m *MockSessionTopic) CreateTopicUnFinishedMessage(topic string, message []broker.UnFinishedMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateTopicUnFinishedMessage", topic, message)
}

// CreateTopicUnFinishedMessage indicates an expected call of CreateTopicUnFinishedMessage.
func (mr *MockSessionTopicMockRecorder) CreateTopicUnFinishedMessage(topic, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTopicUnFinishedMessage", reflect.TypeOf((*MockSessionTopic)(nil).CreateTopicUnFinishedMessage), topic, message)
}

// DeleteSubTopic mocks base method.
func (m *MockSessionTopic) DeleteSubTopic(topic string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteSubTopic", topic)
}

// DeleteSubTopic indicates an expected call of DeleteSubTopic.
func (mr *MockSessionTopicMockRecorder) DeleteSubTopic(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubTopic", reflect.TypeOf((*MockSessionTopic)(nil).DeleteSubTopic), topic)
}

// DeleteTopicLatestPushedMessageID mocks base method.
func (m *MockSessionTopic) DeleteTopicLatestPushedMessageID(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteTopicLatestPushedMessageID", topic, messageID)
}

// DeleteTopicLatestPushedMessageID indicates an expected call of DeleteTopicLatestPushedMessageID.
func (mr *MockSessionTopicMockRecorder) DeleteTopicLatestPushedMessageID(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTopicLatestPushedMessageID", reflect.TypeOf((*MockSessionTopic)(nil).DeleteTopicLatestPushedMessageID), topic, messageID)
}

// DeleteTopicUnFinishedMessage mocks base method.
func (m *MockSessionTopic) DeleteTopicUnFinishedMessage(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteTopicUnFinishedMessage", topic, messageID)
}

// DeleteTopicUnFinishedMessage indicates an expected call of DeleteTopicUnFinishedMessage.
func (mr *MockSessionTopicMockRecorder) DeleteTopicUnFinishedMessage(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTopicUnFinishedMessage", reflect.TypeOf((*MockSessionTopic)(nil).DeleteTopicUnFinishedMessage), topic, messageID)
}

// ReadSubTopics mocks base method.
func (m *MockSessionTopic) ReadSubTopics() map[string]*proto.SubOption {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSubTopics")
	ret0, _ := ret[0].(map[string]*proto.SubOption)
	return ret0
}

// ReadSubTopics indicates an expected call of ReadSubTopics.
func (mr *MockSessionTopicMockRecorder) ReadSubTopics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSubTopics", reflect.TypeOf((*MockSessionTopic)(nil).ReadSubTopics))
}

// ReadTopicLatestPushedMessageID mocks base method.
func (m *MockSessionTopic) ReadTopicLatestPushedMessageID(topic string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTopicLatestPushedMessageID", topic)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ReadTopicLatestPushedMessageID indicates an expected call of ReadTopicLatestPushedMessageID.
func (mr *MockSessionTopicMockRecorder) ReadTopicLatestPushedMessageID(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTopicLatestPushedMessageID", reflect.TypeOf((*MockSessionTopic)(nil).ReadTopicLatestPushedMessageID), topic)
}

// ReadTopicUnFinishedMessage mocks base method.
func (m *MockSessionTopic) ReadTopicUnFinishedMessage(topic string) []broker.UnFinishedMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTopicUnFinishedMessage", topic)
	ret0, _ := ret[0].([]broker.UnFinishedMessage)
	return ret0
}

// ReadTopicUnFinishedMessage indicates an expected call of ReadTopicUnFinishedMessage.
func (mr *MockSessionTopicMockRecorder) ReadTopicUnFinishedMessage(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTopicUnFinishedMessage", reflect.TypeOf((*MockSessionTopic)(nil).ReadTopicUnFinishedMessage), topic)
}

// SetTopicLatestPushedMessageID mocks base method.
func (m *MockSessionTopic) SetTopicLatestPushedMessageID(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTopicLatestPushedMessageID", topic, messageID)
}

// SetTopicLatestPushedMessageID indicates an expected call of SetTopicLatestPushedMessageID.
func (mr *MockSessionTopicMockRecorder) SetTopicLatestPushedMessageID(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTopicLatestPushedMessageID", reflect.TypeOf((*MockSessionTopic)(nil).SetTopicLatestPushedMessageID), topic, messageID)
}

// MockSessionTopicMessage is a mock of SessionTopicMessage interface.
type MockSessionTopicMessage struct {
	ctrl     *gomock.Controller
	recorder *MockSessionTopicMessageMockRecorder
}

// MockSessionTopicMessageMockRecorder is the mock recorder for MockSessionTopicMessage.
type MockSessionTopicMessageMockRecorder struct {
	mock *MockSessionTopicMessage
}

// NewMockSessionTopicMessage creates a new mock instance.
func NewMockSessionTopicMessage(ctrl *gomock.Controller) *MockSessionTopicMessage {
	mock := &MockSessionTopicMessage{ctrl: ctrl}
	mock.recorder = &MockSessionTopicMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionTopicMessage) EXPECT() *MockSessionTopicMessageMockRecorder {
	return m.recorder
}

// CreateTopicUnFinishedMessage mocks base method.
func (m *MockSessionTopicMessage) CreateTopicUnFinishedMessage(topic string, message []broker.UnFinishedMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateTopicUnFinishedMessage", topic, message)
}

// CreateTopicUnFinishedMessage indicates an expected call of CreateTopicUnFinishedMessage.
func (mr *MockSessionTopicMessageMockRecorder) CreateTopicUnFinishedMessage(topic, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTopicUnFinishedMessage", reflect.TypeOf((*MockSessionTopicMessage)(nil).CreateTopicUnFinishedMessage), topic, message)
}

// DeleteTopicLatestPushedMessageID mocks base method.
func (m *MockSessionTopicMessage) DeleteTopicLatestPushedMessageID(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteTopicLatestPushedMessageID", topic, messageID)
}

// DeleteTopicLatestPushedMessageID indicates an expected call of DeleteTopicLatestPushedMessageID.
func (mr *MockSessionTopicMessageMockRecorder) DeleteTopicLatestPushedMessageID(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTopicLatestPushedMessageID", reflect.TypeOf((*MockSessionTopicMessage)(nil).DeleteTopicLatestPushedMessageID), topic, messageID)
}

// DeleteTopicUnFinishedMessage mocks base method.
func (m *MockSessionTopicMessage) DeleteTopicUnFinishedMessage(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteTopicUnFinishedMessage", topic, messageID)
}

// DeleteTopicUnFinishedMessage indicates an expected call of DeleteTopicUnFinishedMessage.
func (mr *MockSessionTopicMessageMockRecorder) DeleteTopicUnFinishedMessage(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTopicUnFinishedMessage", reflect.TypeOf((*MockSessionTopicMessage)(nil).DeleteTopicUnFinishedMessage), topic, messageID)
}

// ReadTopicLatestPushedMessageID mocks base method.
func (m *MockSessionTopicMessage) ReadTopicLatestPushedMessageID(topic string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTopicLatestPushedMessageID", topic)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ReadTopicLatestPushedMessageID indicates an expected call of ReadTopicLatestPushedMessageID.
func (mr *MockSessionTopicMessageMockRecorder) ReadTopicLatestPushedMessageID(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTopicLatestPushedMessageID", reflect.TypeOf((*MockSessionTopicMessage)(nil).ReadTopicLatestPushedMessageID), topic)
}

// ReadTopicUnFinishedMessage mocks base method.
func (m *MockSessionTopicMessage) ReadTopicUnFinishedMessage(topic string) []broker.UnFinishedMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTopicUnFinishedMessage", topic)
	ret0, _ := ret[0].([]broker.UnFinishedMessage)
	return ret0
}

// ReadTopicUnFinishedMessage indicates an expected call of ReadTopicUnFinishedMessage.
func (mr *MockSessionTopicMessageMockRecorder) ReadTopicUnFinishedMessage(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTopicUnFinishedMessage", reflect.TypeOf((*MockSessionTopicMessage)(nil).ReadTopicUnFinishedMessage), topic)
}

// SetTopicLatestPushedMessageID mocks base method.
func (m *MockSessionTopicMessage) SetTopicLatestPushedMessageID(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTopicLatestPushedMessageID", topic, messageID)
}

// SetTopicLatestPushedMessageID indicates an expected call of SetTopicLatestPushedMessageID.
func (mr *MockSessionTopicMessageMockRecorder) SetTopicLatestPushedMessageID(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTopicLatestPushedMessageID", reflect.TypeOf((*MockSessionTopicMessage)(nil).SetTopicLatestPushedMessageID), topic, messageID)
}

// MockSessionTopicUnFinishedMessage is a mock of TopicUnFinishedMessage interface.
type MockSessionTopicUnFinishedMessage struct {
	ctrl     *gomock.Controller
	recorder *MockSessionTopicUnFinishedMessageMockRecorder
}

// MockSessionTopicUnFinishedMessageMockRecorder is the mock recorder for MockSessionTopicUnFinishedMessage.
type MockSessionTopicUnFinishedMessageMockRecorder struct {
	mock *MockSessionTopicUnFinishedMessage
}

// NewMockSessionTopicUnFinishedMessage creates a new mock instance.
func NewMockSessionTopicUnFinishedMessage(ctrl *gomock.Controller) *MockSessionTopicUnFinishedMessage {
	mock := &MockSessionTopicUnFinishedMessage{ctrl: ctrl}
	mock.recorder = &MockSessionTopicUnFinishedMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionTopicUnFinishedMessage) EXPECT() *MockSessionTopicUnFinishedMessageMockRecorder {
	return m.recorder
}

// CreateTopicUnFinishedMessage mocks base method.
func (m *MockSessionTopicUnFinishedMessage) CreateTopicUnFinishedMessage(topic string, message []broker.UnFinishedMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateTopicUnFinishedMessage", topic, message)
}

// CreateTopicUnFinishedMessage indicates an expected call of CreateTopicUnFinishedMessage.
func (mr *MockSessionTopicUnFinishedMessageMockRecorder) CreateTopicUnFinishedMessage(topic, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTopicUnFinishedMessage", reflect.TypeOf((*MockSessionTopicUnFinishedMessage)(nil).CreateTopicUnFinishedMessage), topic, message)
}

// DeleteTopicUnFinishedMessage mocks base method.
func (m *MockSessionTopicUnFinishedMessage) DeleteTopicUnFinishedMessage(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteTopicUnFinishedMessage", topic, messageID)
}

// DeleteTopicUnFinishedMessage indicates an expected call of DeleteTopicUnFinishedMessage.
func (mr *MockSessionTopicUnFinishedMessageMockRecorder) DeleteTopicUnFinishedMessage(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTopicUnFinishedMessage", reflect.TypeOf((*MockSessionTopicUnFinishedMessage)(nil).DeleteTopicUnFinishedMessage), topic, messageID)
}

// ReadTopicUnFinishedMessage mocks base method.
func (m *MockSessionTopicUnFinishedMessage) ReadTopicUnFinishedMessage(topic string) []broker.UnFinishedMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTopicUnFinishedMessage", topic)
	ret0, _ := ret[0].([]broker.UnFinishedMessage)
	return ret0
}

// ReadTopicUnFinishedMessage indicates an expected call of ReadTopicUnFinishedMessage.
func (mr *MockSessionTopicUnFinishedMessageMockRecorder) ReadTopicUnFinishedMessage(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTopicUnFinishedMessage", reflect.TypeOf((*MockSessionTopicUnFinishedMessage)(nil).ReadTopicUnFinishedMessage), topic)
}

// MockSessionTopicLatestPushedMessage is a mock of TopicLatestPushedMessage interface.
type MockSessionTopicLatestPushedMessage struct {
	ctrl     *gomock.Controller
	recorder *MockSessionTopicLatestPushedMessageMockRecorder
}

// MockSessionTopicLatestPushedMessageMockRecorder is the mock recorder for MockSessionTopicLatestPushedMessage.
type MockSessionTopicLatestPushedMessageMockRecorder struct {
	mock *MockSessionTopicLatestPushedMessage
}

// NewMockSessionTopicLatestPushedMessage creates a new mock instance.
func NewMockSessionTopicLatestPushedMessage(ctrl *gomock.Controller) *MockSessionTopicLatestPushedMessage {
	mock := &MockSessionTopicLatestPushedMessage{ctrl: ctrl}
	mock.recorder = &MockSessionTopicLatestPushedMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionTopicLatestPushedMessage) EXPECT() *MockSessionTopicLatestPushedMessageMockRecorder {
	return m.recorder
}

// DeleteTopicLatestPushedMessageID mocks base method.
func (m *MockSessionTopicLatestPushedMessage) DeleteTopicLatestPushedMessageID(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteTopicLatestPushedMessageID", topic, messageID)
}

// DeleteTopicLatestPushedMessageID indicates an expected call of DeleteTopicLatestPushedMessageID.
func (mr *MockSessionTopicLatestPushedMessageMockRecorder) DeleteTopicLatestPushedMessageID(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTopicLatestPushedMessageID", reflect.TypeOf((*MockSessionTopicLatestPushedMessage)(nil).DeleteTopicLatestPushedMessageID), topic, messageID)
}

// ReadTopicLatestPushedMessageID mocks base method.
func (m *MockSessionTopicLatestPushedMessage) ReadTopicLatestPushedMessageID(topic string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTopicLatestPushedMessageID", topic)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ReadTopicLatestPushedMessageID indicates an expected call of ReadTopicLatestPushedMessageID.
func (mr *MockSessionTopicLatestPushedMessageMockRecorder) ReadTopicLatestPushedMessageID(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTopicLatestPushedMessageID", reflect.TypeOf((*MockSessionTopicLatestPushedMessage)(nil).ReadTopicLatestPushedMessageID), topic)
}

// SetTopicLatestPushedMessageID mocks base method.
func (m *MockSessionTopicLatestPushedMessage) SetTopicLatestPushedMessageID(topic, messageID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTopicLatestPushedMessageID", topic, messageID)
}

// SetTopicLatestPushedMessageID indicates an expected call of SetTopicLatestPushedMessageID.
func (mr *MockSessionTopicLatestPushedMessageMockRecorder) SetTopicLatestPushedMessageID(topic, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTopicLatestPushedMessageID", reflect.TypeOf((*MockSessionTopicLatestPushedMessage)(nil).SetTopicLatestPushedMessageID), topic, messageID)
}

// MockSessionCreateConnectProperties is a mock of CreateConnectProperties interface.
type MockSessionCreateConnectProperties struct {
	ctrl     *gomock.Controller
	recorder *MockSessionCreateConnectPropertiesMockRecorder
}

// MockSessionCreateConnectPropertiesMockRecorder is the mock recorder for MockSessionCreateConnectProperties.
type MockSessionCreateConnectPropertiesMockRecorder struct {
	mock *MockSessionCreateConnectProperties
}

// NewMockSessionCreateConnectProperties creates a new mock instance.
func NewMockSessionCreateConnectProperties(ctrl *gomock.Controller) *MockSessionCreateConnectProperties {
	mock := &MockSessionCreateConnectProperties{ctrl: ctrl}
	mock.recorder = &MockSessionCreateConnectPropertiesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionCreateConnectProperties) EXPECT() *MockSessionCreateConnectPropertiesMockRecorder {
	return m.recorder
}

// GetConnectProperties mocks base method.
func (m *MockSessionCreateConnectProperties) GetConnectProperties() (*broker.ConnectProperties, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectProperties")
	ret0, _ := ret[0].(*broker.ConnectProperties)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectProperties indicates an expected call of GetConnectProperties.
func (mr *MockSessionCreateConnectPropertiesMockRecorder) GetConnectProperties() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectProperties", reflect.TypeOf((*MockSessionCreateConnectProperties)(nil).GetConnectProperties))
}

// SetConnectProperties mocks base method.
func (m *MockSessionCreateConnectProperties) SetConnectProperties(properties *broker.ConnectProperties) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetConnectProperties", properties)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetConnectProperties indicates an expected call of SetConnectProperties.
func (mr *MockSessionCreateConnectPropertiesMockRecorder) SetConnectProperties(properties interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConnectProperties", reflect.TypeOf((*MockSessionCreateConnectProperties)(nil).SetConnectProperties), properties)
}

// MockSessionWillMessage is a mock of SessionWillMessage interface.
type MockSessionWillMessage struct {
	ctrl     *gomock.Controller
	recorder *MockSessionWillMessageMockRecorder
}

// MockSessionWillMessageMockRecorder is the mock recorder for MockSessionWillMessage.
type MockSessionWillMessageMockRecorder struct {
	mock *MockSessionWillMessage
}

// NewMockSessionWillMessage creates a new mock instance.
func NewMockSessionWillMessage(ctrl *gomock.Controller) *MockSessionWillMessage {
	mock := &MockSessionWillMessage{ctrl: ctrl}
	mock.recorder = &MockSessionWillMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionWillMessage) EXPECT() *MockSessionWillMessageMockRecorder {
	return m.recorder
}

// GetWillMessage mocks base method.
func (m *MockSessionWillMessage) GetWillMessage() (*broker.WillMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWillMessage")
	ret0, _ := ret[0].(*broker.WillMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWillMessage indicates an expected call of GetWillMessage.
func (mr *MockSessionWillMessageMockRecorder) GetWillMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWillMessage", reflect.TypeOf((*MockSessionWillMessage)(nil).GetWillMessage))
}

// SetWillMessage mocks base method.
func (m *MockSessionWillMessage) SetWillMessage(message *broker.WillMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWillMessage", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWillMessage indicates an expected call of SetWillMessage.
func (mr *MockSessionWillMessageMockRecorder) SetWillMessage(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWillMessage", reflect.TypeOf((*MockSessionWillMessage)(nil).SetWillMessage), message)
}
