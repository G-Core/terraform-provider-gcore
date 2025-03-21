package gcore

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

type mockAcc struct {
	sdk.ClientSDKmock
}

func (mockAcc) GetApp(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	compare(contextTesting(ctx), id, contextExpected[int64](ctx))
	return contextReturn(ctx), nil
}

func (mockAcc) AddApp(ctx context.Context, body sdk.AddAppJSONRequestBody, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	compare(contextTesting(ctx), body, contextExpected[sdk.App](ctx))
	return contextReturn(ctx), nil
}

var mockAccClient *sdk.ClientWithResponses = &sdk.ClientWithResponses{ClientInterface: &mockAcc{}}

func TestReadApp(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceFastEdgeApp().Schema, nil)
	resourceData.SetId("42")

	ctx := contextWithTesting(context.Background(), t)
	ctx = contextWithExpected(ctx, int64(42))
	ctx = contextWithReturn(ctx,
		sdk.NewJsonHttpResponse(
			http.StatusOK,
			`{
				"id": 42,
				"status": 1,
				"name": "test_app",
				"binary": 314,
				"template": 15,
				"env": {
					"key": "value"
				},
				"rsp_headers": {
					"key1": "value1"
				},
				"debug": true,
				"comment": "test application"
			}`,
		),
	)

	diag := resourceFastEdgeAppRead(ctx, resourceData, &Config{FastEdgeClient: mockAccClient})

	if diag.HasError() {
		t.Error("unexpected error", diag[0].Summary)
	} else {
		compare(t, resourceData.Get("status"), "enabled")
		compare(t, resourceData.Get("name"), "test_app")
		compare(t, resourceData.Get("binary"), 314)
		compare(t, resourceData.Get("template"), 15)
		compare(t, resourceData.Get("env").(map[string]interface{}), map[string]interface{}{"key": "value"})
		compare(t, resourceData.Get("rsp_headers").(map[string]interface{}), map[string]interface{}{"key1": "value1"})
		compare(t, resourceData.Get("comment"), "test application")
		compare(t, resourceData.Get("debug"), true)
	}
}

func TestReadAppNotFound(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceFastEdgeApp().Schema, nil)
	resourceData.SetId("42")

	ctx := contextWithTesting(context.Background(), t)
	ctx = contextWithExpected(ctx, int64(42))
	ctx = contextWithReturn(ctx, sdk.NewJsonHttpResponse(http.StatusNotFound, "Not found"))

	diag := resourceFastEdgeAppRead(ctx, resourceData, &Config{FastEdgeClient: mockAccClient})

	if diag.HasError() {
		t.Error("unexpected error", diag[0].Summary)
	} else {
		compare(t, resourceData.Id(), "")
	}
}

func TestReadAppError(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceFastEdgeApp().Schema, nil)
	resourceData.SetId("42")

	ctx := contextWithTesting(context.Background(), t)
	ctx = contextWithExpected(ctx, int64(42))
	ctx = contextWithReturn(ctx, sdk.NewJsonHttpResponse(http.StatusBadRequest, `{"error": "Something is wrong"}`))

	diag := resourceFastEdgeAppRead(ctx, resourceData, &Config{FastEdgeClient: mockAccClient})

	if diag.HasError() {
		if diag[0].Summary != "Something is wrong" {
			t.Error("unexpected error", diag[0].Summary)
		}
	}
	compare(t, resourceData.Id(), "42") // State is not cleaned
}

func TestAddAppWithBinary(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceFastEdgeApp().Schema, nil)
	resourceData.Set("name", "test_app")
	resourceData.Set("binary", 314)
	resourceData.Set("status", "enabled")
	resourceData.Set("env", map[string]interface{}{"key": "value"})
	resourceData.Set("rsp_headers", map[string]interface{}{"key1": "value1"})
	resourceData.Set("comment", "test application")
	resourceData.Set("debug", true)

	ctx := contextWithTesting(context.Background(), t)
	ctx = contextWithExpected(ctx, sdk.App{
		Name:   ptr("test_app"),
		Binary: ptr[int64](314),
		Env: ptr(map[string]string{
			"key": "value",
		}),
		Comment: ptr("test application"),
		Status:  ptr(1),
		Debug:   ptr(true),
		RspHeaders: ptr(map[string]string{
			"key1": "value1",
		}),
	})
	ctx = contextWithReturn(ctx,
		sdk.NewJsonHttpResponse(
			http.StatusOK,
			`{
				"id": 42,
				"status": 1,
				"name": "test_app",
				"binary": 314,
				"env": {
					"key": "value"
				},
				"rsp_headers": {
					"key1": "value1"
				},
				"debug": true,
				"comment": "test application"
			}`,
		),
	)

	diag := resourceFastEdgeAppCreate(ctx, resourceData, &Config{FastEdgeClient: mockAccClient})

	if diag.HasError() {
		t.Error("unexpected error", diag[0].Summary)
	} else {
		compare(t, resourceData.Id(), "42")
		compare(t, resourceData.Get("status"), "enabled")
		compare(t, resourceData.Get("name"), "test_app")
		compare(t, resourceData.Get("binary"), 314)
		compare(t, resourceData.Get("env").(map[string]interface{}), map[string]interface{}{"key": "value"})
		compare(t, resourceData.Get("rsp_headers").(map[string]interface{}), map[string]interface{}{"key1": "value1"})
		compare(t, resourceData.Get("comment"), "test application")
		compare(t, resourceData.Get("debug"), true)
	}
}

func TestAddAppWithTemplate(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceFastEdgeApp().Schema, nil)
	resourceData.Set("name", "test_app")
	resourceData.Set("template", 15)
	resourceData.Set("status", "enabled")
	resourceData.Set("env", map[string]interface{}{"key": "value"})
	resourceData.Set("rsp_headers", map[string]interface{}{"key1": "value1"})
	resourceData.Set("comment", "test application")
	resourceData.Set("debug", true)

	ctx := contextWithTesting(context.Background(), t)
	ctx = contextWithExpected(ctx, sdk.App{
		Name:     ptr("test_app"),
		Template: ptr[int64](15),
		Env: ptr(map[string]string{
			"key": "value",
		}),
		Comment: ptr("test application"),
		Status:  ptr(1),
		Debug:   ptr(true),
		RspHeaders: ptr(map[string]string{
			"key1": "value1",
		}),
	})
	ctx = contextWithReturn(ctx,
		sdk.NewJsonHttpResponse(
			http.StatusOK,
			`{
				"id": 42,
				"status": 1,
				"name": "test_app",
				"binary": 314,
				"template": 15,
				"env": {
					"key": "value"
				},
				"rsp_headers": {
					"key1": "value1"
				},
				"debug": true,
				"comment": "test application"
			}`,
		),
	)

	diag := resourceFastEdgeAppCreate(ctx, resourceData, &Config{FastEdgeClient: mockAccClient})

	if diag.HasError() {
		t.Error("unexpected error", diag[0].Summary)
	} else {
		compare(t, resourceData.Id(), "42")
		compare(t, resourceData.Get("status"), "enabled")
		compare(t, resourceData.Get("name"), "test_app")
		compare(t, resourceData.Get("binary"), 314)
		compare(t, resourceData.Get("template"), 15)
		compare(t, resourceData.Get("env").(map[string]interface{}), map[string]interface{}{"key": "value"})
		compare(t, resourceData.Get("rsp_headers").(map[string]interface{}), map[string]interface{}{"key1": "value1"})
		compare(t, resourceData.Get("comment"), "test application")
		compare(t, resourceData.Get("debug"), true)
	}
}

func ptr[T any](val T) *T {
	return &val
}
