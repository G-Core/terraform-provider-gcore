package gcore

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
		panic("unexpected call to DelBinary")
	}
	defer func() { m.delBinCount++ }()
	compare(m.t, id, m.delBin[m.delBinCount].expectId)
	return sdk.NewJsonHttpResponse(m.delBin[m.delBinCount].retStatus, m.delBin[m.delBinCount].retBody), nil
}

func (m mockSdk) ExpectationsWereMet(t *testing.T) {
	t.Helper()
	if m.addAppCount != len(m.addApp) {
		t.Errorf("expected %d calls to AddApp, got %d", len(m.addApp), m.addAppCount)
	}
	if m.updAppCount != len(m.updApp) {
		t.Errorf("expected %d calls to UpdateApp, got %d", len(m.updApp), m.updAppCount)
	}
	if m.delAppCount != len(m.delApp) {
		t.Errorf("expected %d calls to DelApp, got %d", len(m.delApp), m.delAppCount)
	}
	if m.getAppCount != len(m.getApp) {
		t.Errorf("expected %d calls to GetApp, got %d", len(m.getApp), m.getAppCount)
	}
	if m.addBinCount != len(m.addBin) {
		t.Errorf("expected %d calls to StoreBinary, got %d", len(m.addBin), m.addBinCount)
	}
	if m.getBinCount != len(m.getBin) {
		t.Errorf("expected %d calls to GetBinary, got %d", len(m.getBin), m.getBinCount)
	}
	if m.delBinCount != len(m.delBin) {
		t.Errorf("expected %d calls to DelBinary, got %d", len(m.delBin), m.delBinCount)
	}
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

func fastedgeMockProvider(mock sdk.ClientInterface) map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"gcore": func() (*schema.Provider, error) {
			return &schema.Provider{
				ResourcesMap: map[string]*schema.Resource{
					"gcore_fastedge_binary": resourceFastEdgeBinary(),
					"gcore_fastedge_app":    resourceFastEdgeApp(),
				},
				ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
					config := Config{FastEdgeClient: &sdk.ClientWithResponses{ClientInterface: mock}}
					return &config, nil
				},
			}, nil
		},
	}
}
