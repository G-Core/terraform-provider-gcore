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
				"status":1,
				"name": "test_app",
				"binary": 314,
				"env": {
					"key": "value"
				},
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
		compare(t, resourceData.Get("env").(map[string]interface{}), map[string]interface{}{"key": "value"})
		compare(t, resourceData.Get("comment"), "test application")
	}
}
