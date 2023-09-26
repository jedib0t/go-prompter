// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jedib0t/go-prompter/input (interfaces: Reader)
//
// Generated by this command:
//
//	mockgen -destination mocks/input/mock_reader.go github.com/jedib0t/go-prompter/input Reader
//
// Package mock_input is a generated GoMock package.
package mock_input

import (
	context "context"
	reflect "reflect"

	tea "github.com/charmbracelet/bubbletea"
	gomock "go.uber.org/mock/gomock"
)

// MockReader is a mock of Reader interface.
type MockReader struct {
	ctrl     *gomock.Controller
	recorder *MockReaderMockRecorder
}

// MockReaderMockRecorder is the mock recorder for MockReader.
type MockReaderMockRecorder struct {
	mock *MockReader
}

// NewMockReader creates a new mock instance.
func NewMockReader(ctrl *gomock.Controller) *MockReader {
	mock := &MockReader{ctrl: ctrl}
	mock.recorder = &MockReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReader) EXPECT() *MockReaderMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockReader) Begin(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Begin", arg0)
}

// Begin indicates an expected call of Begin.
func (mr *MockReaderMockRecorder) Begin(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockReader)(nil).Begin), arg0)
}

// End mocks base method.
func (m *MockReader) End() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "End")
}

// End indicates an expected call of End.
func (mr *MockReaderMockRecorder) End() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "End", reflect.TypeOf((*MockReader)(nil).End))
}

// Errors mocks base method.
func (m *MockReader) Errors() <-chan error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Errors")
	ret0, _ := ret[0].(<-chan error)
	return ret0
}

// Errors indicates an expected call of Errors.
func (mr *MockReaderMockRecorder) Errors() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errors", reflect.TypeOf((*MockReader)(nil).Errors))
}

// KeyEvents mocks base method.
func (m *MockReader) KeyEvents() <-chan tea.KeyMsg {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeyEvents")
	ret0, _ := ret[0].(<-chan tea.KeyMsg)
	return ret0
}

// KeyEvents indicates an expected call of KeyEvents.
func (mr *MockReaderMockRecorder) KeyEvents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeyEvents", reflect.TypeOf((*MockReader)(nil).KeyEvents))
}

// MouseEvents mocks base method.
func (m *MockReader) MouseEvents() <-chan tea.MouseMsg {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MouseEvents")
	ret0, _ := ret[0].(<-chan tea.MouseMsg)
	return ret0
}

// MouseEvents indicates an expected call of MouseEvents.
func (mr *MockReaderMockRecorder) MouseEvents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MouseEvents", reflect.TypeOf((*MockReader)(nil).MouseEvents))
}

// Reset mocks base method.
func (m *MockReader) Reset() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reset")
	ret0, _ := ret[0].(error)
	return ret0
}

// Reset indicates an expected call of Reset.
func (mr *MockReaderMockRecorder) Reset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockReader)(nil).Reset))
}

// Send mocks base method.
func (m *MockReader) Send(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockReaderMockRecorder) Send(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockReader)(nil).Send), arg0)
}

// WindowSizeEvents mocks base method.
func (m *MockReader) WindowSizeEvents() <-chan tea.WindowSizeMsg {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WindowSizeEvents")
	ret0, _ := ret[0].(<-chan tea.WindowSizeMsg)
	return ret0
}

// WindowSizeEvents indicates an expected call of WindowSizeEvents.
func (mr *MockReaderMockRecorder) WindowSizeEvents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WindowSizeEvents", reflect.TypeOf((*MockReader)(nil).WindowSizeEvents))
}
