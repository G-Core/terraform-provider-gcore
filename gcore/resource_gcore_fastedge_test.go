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

type funcMock struct {
	params []mockParams
	count  int
}

type mockSDK struct {
	sdk.ClientSDKmock
	t     *testing.T
	mocks map[string]*funcMock
}

func (m *mockSDK) mock(name string, id int64, body any) (*http.Response, error) {
	m.t.Helper()
	mock, ok := m.mocks[name]
	if !ok || mock.count >= len(mock.params) {
		panic("unexpected call to " + name)
	}
	defer func() { mock.count++ }()
	params := mock.params[mock.count]
	if params.expectId != 0 {
		compare(m.t, id, params.expectId, name, mock.count)
	}
	if params.expectPayload != nil {
		compare(m.t, body, params.expectPayload, name, mock.count)
	}
	return sdk.NewJsonHttpResponse(params.retStatus, params.retBody), nil
}

func (m *mockSDK) GetApp(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return m.mock("GetApp", id, nil)
}
func (m *mockSDK) AddApp(ctx context.Context, body sdk.AddAppJSONRequestBody, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return m.mock("AddApp", 0, body)
}
func (m *mockSDK) UpdateApp(ctx context.Context, id int64, body sdk.UpdateAppJSONRequestBody, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return m.mock("UpdateApp", id, body)
}
func (m *mockSDK) DelApp(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return m.mock("DelApp", id, nil)
}
func (m *mockSDK) GetBinary(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return m.mock("GetBinary", id, nil)
}
func (m *mockSDK) StoreBinaryWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return m.mock("StoreBinary", 0, nil)
}
func (m *mockSDK) DelBinary(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return m.mock("DelBinary", id, nil)
}
func (m *mockSDK) GetTemplate(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return m.mock("GetTemplate", id, nil)
}
func (m *mockSDK) AddTemplate(ctx context.Context, body sdk.AddTemplateJSONRequestBody, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return m.mock("AddTemplate", 0, body)
}
func (m *mockSDK) UpdateTemplate(ctx context.Context, id int64, body sdk.UpdateTemplateJSONRequestBody, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return m.mock("UpdateTemplate", id, body)
}
func (m *mockSDK) DelTemplate(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return m.mock("DelTemplate", id, nil)
}

func (m mockSDK) ExpectationsWereMet(t *testing.T) {
	t.Helper()
	for k, v := range m.mocks {
		if v.count != len(v.params) {
			t.Errorf("expected %d calls to %s, got %d", len(v.params), k, v.count)
		}
	}
}

func ptr[T any](val T) *T {
	return &val
}

func compare(t *testing.T, actual, expected any, name string, run int) {
	t.Helper()
	if !cmp.Equal(actual, expected) {
		t.Errorf("%s (%d): unexpected value, got: %#v, expected: %#v", name, run, actual, expected)
	}
}

func fastedgeMockProvider(mock sdk.ClientInterface) map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"gcore": func() (*schema.Provider, error) {
			return &schema.Provider{
				ResourcesMap: map[string]*schema.Resource{
					"gcore_fastedge_binary":   resourceFastEdgeBinary(),
					"gcore_fastedge_app":      resourceFastEdgeApp(),
					"gcore_fastedge_template": resourceFastEdgeTemplate(),
				},
				ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
					config := Config{FastEdgeClient: &sdk.ClientWithResponses{ClientInterface: mock}}
					return &config, nil
				},
			}, nil
		},
	}
}
