// Code generated by MockGen. DO NOT EDIT.
// Source: artifact_tracker.pb.go

// Package mock_artifacttrackerpb is a generated GoMock package.
package mock_artifacttrackerpb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	artifacttrackerpb "pixielabs.ai/pixielabs/src/cloud/artifact_tracker/artifacttrackerpb"
	versionspb "pixielabs.ai/pixielabs/src/shared/artifacts/versionspb"
)

// MockArtifactTrackerClient is a mock of ArtifactTrackerClient interface.
type MockArtifactTrackerClient struct {
	ctrl     *gomock.Controller
	recorder *MockArtifactTrackerClientMockRecorder
}

// MockArtifactTrackerClientMockRecorder is the mock recorder for MockArtifactTrackerClient.
type MockArtifactTrackerClientMockRecorder struct {
	mock *MockArtifactTrackerClient
}

// NewMockArtifactTrackerClient creates a new mock instance.
func NewMockArtifactTrackerClient(ctrl *gomock.Controller) *MockArtifactTrackerClient {
	mock := &MockArtifactTrackerClient{ctrl: ctrl}
	mock.recorder = &MockArtifactTrackerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArtifactTrackerClient) EXPECT() *MockArtifactTrackerClientMockRecorder {
	return m.recorder
}

// GetArtifactList mocks base method.
func (m *MockArtifactTrackerClient) GetArtifactList(ctx context.Context, in *artifacttrackerpb.GetArtifactListRequest, opts ...grpc.CallOption) (*versionspb.ArtifactSet, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetArtifactList", varargs...)
	ret0, _ := ret[0].(*versionspb.ArtifactSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArtifactList indicates an expected call of GetArtifactList.
func (mr *MockArtifactTrackerClientMockRecorder) GetArtifactList(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArtifactList", reflect.TypeOf((*MockArtifactTrackerClient)(nil).GetArtifactList), varargs...)
}

// GetDownloadLink mocks base method.
func (m *MockArtifactTrackerClient) GetDownloadLink(ctx context.Context, in *artifacttrackerpb.GetDownloadLinkRequest, opts ...grpc.CallOption) (*artifacttrackerpb.GetDownloadLinkResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetDownloadLink", varargs...)
	ret0, _ := ret[0].(*artifacttrackerpb.GetDownloadLinkResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDownloadLink indicates an expected call of GetDownloadLink.
func (mr *MockArtifactTrackerClientMockRecorder) GetDownloadLink(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDownloadLink", reflect.TypeOf((*MockArtifactTrackerClient)(nil).GetDownloadLink), varargs...)
}

// MockArtifactTrackerServer is a mock of ArtifactTrackerServer interface.
type MockArtifactTrackerServer struct {
	ctrl     *gomock.Controller
	recorder *MockArtifactTrackerServerMockRecorder
}

// MockArtifactTrackerServerMockRecorder is the mock recorder for MockArtifactTrackerServer.
type MockArtifactTrackerServerMockRecorder struct {
	mock *MockArtifactTrackerServer
}

// NewMockArtifactTrackerServer creates a new mock instance.
func NewMockArtifactTrackerServer(ctrl *gomock.Controller) *MockArtifactTrackerServer {
	mock := &MockArtifactTrackerServer{ctrl: ctrl}
	mock.recorder = &MockArtifactTrackerServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArtifactTrackerServer) EXPECT() *MockArtifactTrackerServerMockRecorder {
	return m.recorder
}

// GetArtifactList mocks base method.
func (m *MockArtifactTrackerServer) GetArtifactList(arg0 context.Context, arg1 *artifacttrackerpb.GetArtifactListRequest) (*versionspb.ArtifactSet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArtifactList", arg0, arg1)
	ret0, _ := ret[0].(*versionspb.ArtifactSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArtifactList indicates an expected call of GetArtifactList.
func (mr *MockArtifactTrackerServerMockRecorder) GetArtifactList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArtifactList", reflect.TypeOf((*MockArtifactTrackerServer)(nil).GetArtifactList), arg0, arg1)
}

// GetDownloadLink mocks base method.
func (m *MockArtifactTrackerServer) GetDownloadLink(arg0 context.Context, arg1 *artifacttrackerpb.GetDownloadLinkRequest) (*artifacttrackerpb.GetDownloadLinkResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDownloadLink", arg0, arg1)
	ret0, _ := ret[0].(*artifacttrackerpb.GetDownloadLinkResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDownloadLink indicates an expected call of GetDownloadLink.
func (mr *MockArtifactTrackerServerMockRecorder) GetDownloadLink(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDownloadLink", reflect.TypeOf((*MockArtifactTrackerServer)(nil).GetDownloadLink), arg0, arg1)
}
