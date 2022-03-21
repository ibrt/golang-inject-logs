// Code generated by MockGen. DO NOT EDIT.
// Source: ../logs.go

// Package mocklogz is a generated GoMock package.
package mocklogz

import (
	context "context"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	logz "github.com/ibrt/golang-inject-logs/logz"
)

// MockLogs is a mock of Logs interface.
type MockLogs struct {
	ctrl     *gomock.Controller
	recorder *MockLogsMockRecorder
}

// MockLogsMockRecorder is the mock recorder for MockLogs.
type MockLogsMockRecorder struct {
	mock *MockLogs
}

// NewMockLogs creates a new mock instance.
func NewMockLogs(ctrl *gomock.Controller) *MockLogs {
	mock := &MockLogs{ctrl: ctrl}
	mock.recorder = &MockLogsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogs) EXPECT() *MockLogsMockRecorder {
	return m.recorder
}

// AddMetadata mocks base method.
func (m *MockLogs) AddMetadata(ctx context.Context, k string, v interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddMetadata", ctx, k, v)
}

// AddMetadata indicates an expected call of AddMetadata.
func (mr *MockLogsMockRecorder) AddMetadata(ctx, k, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetadata", reflect.TypeOf((*MockLogs)(nil).AddMetadata), ctx, k, v)
}

// Debug mocks base method.
func (m *MockLogs) Debug(ctx context.Context, skipCallers int, format string, options ...logz.Option) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, skipCallers, format}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockLogsMockRecorder) Debug(ctx, skipCallers, format interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, skipCallers, format}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogs)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *MockLogs) Error(ctx context.Context, err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", ctx, err)
}

// Error indicates an expected call of Error.
func (mr *MockLogsMockRecorder) Error(ctx, err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogs)(nil).Error), ctx, err)
}

// Info mocks base method.
func (m *MockLogs) Info(ctx context.Context, skipCallers int, format string, options ...logz.Option) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, skipCallers, format}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockLogsMockRecorder) Info(ctx, skipCallers, format interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, skipCallers, format}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogs)(nil).Info), varargs...)
}

// SetUser mocks base method.
func (m *MockLogs) SetUser(ctx context.Context, user *logz.User) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetUser", ctx, user)
}

// SetUser indicates an expected call of SetUser.
func (mr *MockLogsMockRecorder) SetUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUser", reflect.TypeOf((*MockLogs)(nil).SetUser), ctx, user)
}

// TraceHTTPRequestServer mocks base method.
func (m *MockLogs) TraceHTTPRequestServer(ctx context.Context, req *http.Request, reqBody []byte) (context.Context, func()) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TraceHTTPRequestServer", ctx, req, reqBody)
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(func())
	return ret0, ret1
}

// TraceHTTPRequestServer indicates an expected call of TraceHTTPRequestServer.
func (mr *MockLogsMockRecorder) TraceHTTPRequestServer(ctx, req, reqBody interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TraceHTTPRequestServer", reflect.TypeOf((*MockLogs)(nil).TraceHTTPRequestServer), ctx, req, reqBody)
}

// TraceSpan mocks base method.
func (m *MockLogs) TraceSpan(ctx context.Context, op, desc string) (context.Context, func()) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TraceSpan", ctx, op, desc)
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(func())
	return ret0, ret1
}

// TraceSpan indicates an expected call of TraceSpan.
func (mr *MockLogsMockRecorder) TraceSpan(ctx, op, desc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TraceSpan", reflect.TypeOf((*MockLogs)(nil).TraceSpan), ctx, op, desc)
}

// Warning mocks base method.
func (m *MockLogs) Warning(ctx context.Context, err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Warning", ctx, err)
}

// Warning indicates an expected call of Warning.
func (mr *MockLogsMockRecorder) Warning(ctx, err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warning", reflect.TypeOf((*MockLogs)(nil).Warning), ctx, err)
}

// MockContextLogs is a mock of ContextLogs interface.
type MockContextLogs struct {
	ctrl     *gomock.Controller
	recorder *MockContextLogsMockRecorder
}

// MockContextLogsMockRecorder is the mock recorder for MockContextLogs.
type MockContextLogsMockRecorder struct {
	mock *MockContextLogs
}

// NewMockContextLogs creates a new mock instance.
func NewMockContextLogs(ctrl *gomock.Controller) *MockContextLogs {
	mock := &MockContextLogs{ctrl: ctrl}
	mock.recorder = &MockContextLogsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContextLogs) EXPECT() *MockContextLogsMockRecorder {
	return m.recorder
}

// AddMetadata mocks base method.
func (m *MockContextLogs) AddMetadata(k string, v interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddMetadata", k, v)
}

// AddMetadata indicates an expected call of AddMetadata.
func (mr *MockContextLogsMockRecorder) AddMetadata(k, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetadata", reflect.TypeOf((*MockContextLogs)(nil).AddMetadata), k, v)
}

// Debug mocks base method.
func (m *MockContextLogs) Debug(format string, options ...logz.Option) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockContextLogsMockRecorder) Debug(format interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockContextLogs)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *MockContextLogs) Error(err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", err)
}

// Error indicates an expected call of Error.
func (mr *MockContextLogsMockRecorder) Error(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockContextLogs)(nil).Error), err)
}

// Info mocks base method.
func (m *MockContextLogs) Info(format string, options ...logz.Option) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockContextLogsMockRecorder) Info(format interface{}, options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, options...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockContextLogs)(nil).Info), varargs...)
}

// SetUser mocks base method.
func (m *MockContextLogs) SetUser(user *logz.User) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetUser", user)
}

// SetUser indicates an expected call of SetUser.
func (mr *MockContextLogsMockRecorder) SetUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUser", reflect.TypeOf((*MockContextLogs)(nil).SetUser), user)
}

// TraceHTTPRequestServer mocks base method.
func (m *MockContextLogs) TraceHTTPRequestServer(req *http.Request, reqBody []byte) (context.Context, func()) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TraceHTTPRequestServer", req, reqBody)
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(func())
	return ret0, ret1
}

// TraceHTTPRequestServer indicates an expected call of TraceHTTPRequestServer.
func (mr *MockContextLogsMockRecorder) TraceHTTPRequestServer(req, reqBody interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TraceHTTPRequestServer", reflect.TypeOf((*MockContextLogs)(nil).TraceHTTPRequestServer), req, reqBody)
}

// TraceSpan mocks base method.
func (m *MockContextLogs) TraceSpan(op, desc string) (context.Context, func()) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TraceSpan", op, desc)
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(func())
	return ret0, ret1
}

// TraceSpan indicates an expected call of TraceSpan.
func (mr *MockContextLogsMockRecorder) TraceSpan(op, desc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TraceSpan", reflect.TypeOf((*MockContextLogs)(nil).TraceSpan), op, desc)
}

// Warning mocks base method.
func (m *MockContextLogs) Warning(err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Warning", err)
}

// Warning indicates an expected call of Warning.
func (mr *MockContextLogsMockRecorder) Warning(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warning", reflect.TypeOf((*MockContextLogs)(nil).Warning), err)
}
