// Code generated by MockGen. DO NOT EDIT.
// Source: server.go

// Package mock_controller is a generated GoMock package.
package mock_controller

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	cvmsgspb "pixielabs.ai/pixielabs/src/shared/cvmsgspb"
)

// MockVzUpdater is a mock of VzUpdater interface.
type MockVzUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockVzUpdaterMockRecorder
}

// MockVzUpdaterMockRecorder is the mock recorder for MockVzUpdater.
type MockVzUpdaterMockRecorder struct {
	mock *MockVzUpdater
}

// NewMockVzUpdater creates a new mock instance.
func NewMockVzUpdater(ctrl *gomock.Controller) *MockVzUpdater {
	mock := &MockVzUpdater{ctrl: ctrl}
	mock.recorder = &MockVzUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVzUpdater) EXPECT() *MockVzUpdaterMockRecorder {
	return m.recorder
}

// AddToUpdateQueue mocks base method.
func (m *MockVzUpdater) AddToUpdateQueue(vizierID uuid.UUID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToUpdateQueue", vizierID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// AddToUpdateQueue indicates an expected call of AddToUpdateQueue.
func (mr *MockVzUpdaterMockRecorder) AddToUpdateQueue(vizierID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToUpdateQueue", reflect.TypeOf((*MockVzUpdater)(nil).AddToUpdateQueue), vizierID)
}

// UpdateOrInstallVizier mocks base method.
func (m *MockVzUpdater) UpdateOrInstallVizier(vizierID uuid.UUID, version string, redeployEtcd bool) (*cvmsgspb.V2CMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrInstallVizier", vizierID, version, redeployEtcd)
	ret0, _ := ret[0].(*cvmsgspb.V2CMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateOrInstallVizier indicates an expected call of UpdateOrInstallVizier.
func (mr *MockVzUpdaterMockRecorder) UpdateOrInstallVizier(vizierID, version, redeployEtcd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrInstallVizier", reflect.TypeOf((*MockVzUpdater)(nil).UpdateOrInstallVizier), vizierID, version, redeployEtcd)
}

// VersionUpToDate mocks base method.
func (m *MockVzUpdater) VersionUpToDate(version string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VersionUpToDate", version)
	ret0, _ := ret[0].(bool)
	return ret0
}

// VersionUpToDate indicates an expected call of VersionUpToDate.
func (mr *MockVzUpdaterMockRecorder) VersionUpToDate(version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VersionUpToDate", reflect.TypeOf((*MockVzUpdater)(nil).VersionUpToDate), version)
}
