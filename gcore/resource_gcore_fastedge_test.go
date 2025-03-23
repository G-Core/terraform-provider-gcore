package gcore

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

type mockParams struct {
	expectId      int64
	expectPayload any
	retStatus     int
	retBody       string
}

type mockSdk struct {
	sdk.ClientSDKmock
	t           *testing.T
	addApp      []mockParams
	addAppCount int
	getApp      []mockParams
	getAppCount int
	updApp      []mockParams
	updAppCount int
	delApp      []mockParams
	delAppCount int
	addBin      []mockParams
	addBinCount int
	getBin      []mockParams
	getBinCount int
	delBin      []mockParams
	delBinCount int
}

func (m *mockSdk) GetApp(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	if m.getAppCount >= len(m.getApp) {
		panic("unexpected call to GetApp")
	}
	defer func() { m.getAppCount++ }()
	compare(m.t, id, m.getApp[m.getAppCount].expectId)
	return sdk.NewJsonHttpResponse(m.getApp[m.getAppCount].retStatus, m.getApp[m.getAppCount].retBody), nil
}

func (m *mockSdk) AddApp(ctx context.Context, body sdk.AddAppJSONRequestBody, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	if m.addAppCount >= len(m.addApp) {
		panic("unexpected call to AddApp")
	}
	defer func() { m.addAppCount++ }()
	compare(m.t, body, m.addApp[m.addAppCount].expectPayload.(sdk.App))
	return sdk.NewJsonHttpResponse(m.addApp[m.addAppCount].retStatus, m.addApp[m.addAppCount].retBody), nil
}

func (m *mockSdk) UpdateApp(ctx context.Context, id int64, body sdk.UpdateAppJSONRequestBody, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	if m.updAppCount >= len(m.updApp) {
		panic("unexpected call to UpdateApp")
	}
	defer func() { m.updAppCount++ }()
	compare(m.t, id, m.updApp[m.updAppCount].expectId)
	compare(m.t, body.App, m.updApp[m.updAppCount].expectPayload.(sdk.App))
	return sdk.NewJsonHttpResponse(m.updApp[m.updAppCount].retStatus, m.updApp[m.updAppCount].retBody), nil
}

func (m *mockSdk) DelApp(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	if m.delAppCount >= len(m.delApp) {
		panic("unexpected call to DelApp")
	}
	defer func() { m.delAppCount++ }()
	compare(m.t, id, m.delApp[m.delAppCount].expectId)
	return sdk.NewJsonHttpResponse(m.delApp[m.delAppCount].retStatus, m.delApp[m.delAppCount].retBody), nil
}

func (m *mockSdk) GetBinary(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	if m.getBinCount >= len(m.getBin) {
		panic("unexpected call to GetBinary")
	}
	defer func() { m.getBinCount++ }()
	compare(m.t, id, m.getBin[m.getBinCount].expectId)
	return sdk.NewJsonHttpResponse(m.getBin[m.getBinCount].retStatus, m.getBin[m.getBinCount].retBody), nil
}

func (m *mockSdk) StoreBinaryWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	if m.addBinCount >= len(m.addBin) {
		panic("unexpected call to StoreBinary")
	}
	defer func() { m.addBinCount++ }()
	return sdk.NewJsonHttpResponse(m.addBin[m.addAppCount].retStatus, m.addBin[m.addBinCount].retBody), nil
}

func (m *mockSdk) DelBinary(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	if m.delBinCount >= len(m.delBin) {
		panic("unexpected call to DelApp")
	}
	defer func() { m.delBinCount++ }()
	compare(m.t, id, m.delBin[m.delBinCount].expectId)
	return sdk.NewJsonHttpResponse(m.delBin[m.delBinCount].retStatus, m.delBin[m.delBinCount].retBody), nil
}

func ptr[T any](val T) *T {
	return &val
}

func compare[T any](t *testing.T, actual, expected T) {
	t.Helper()
	if !cmp.Equal(actual, expected) {
		t.Errorf("unexpected value, got: %#v, expected: %#v", actual, expected)
	}
}
