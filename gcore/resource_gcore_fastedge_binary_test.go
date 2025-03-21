package gcore

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

type mockBin struct {
	sdk.ClientSDKmock
}

func (mockBin) GetBinary(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	compare(contextTesting(ctx), id, contextExpected[int64](ctx))
	return contextReturn(ctx), nil
}

func (mockBin) StoreBinaryWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return contextReturn(ctx), nil
}

var mockBinClient *sdk.ClientWithResponses = &sdk.ClientWithResponses{ClientInterface: &mockBin{}}

func TestReadBinary(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceFastEdgeBinary().Schema, nil)
	resourceData.SetId("42")

	ctx := contextWithTesting(context.Background(), t)
	ctx = contextWithExpected(ctx, int64(42))
	ctx = contextWithReturn(ctx,
		sdk.NewJsonHttpResponse(
			http.StatusOK,
			`{
				"id": 42,
				"checksum":"1234567890"
			}`,
		),
	)

	diag := resourceFastEdgeBinaryRead(ctx, resourceData, &Config{FastEdgeClient: mockBinClient})

	if diag.HasError() {
		t.Error("unexpected error", diag[0].Summary)
	} else {
		compare(t, resourceData.Get("checksum"), "1234567890")
	}
}

func TestReadBinaryNotFound(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceFastEdgeBinary().Schema, nil)
	resourceData.SetId("42")

	ctx := contextWithTesting(context.Background(), t)
	ctx = contextWithExpected(ctx, int64(42))
	ctx = contextWithReturn(ctx, sdk.NewJsonHttpResponse(http.StatusNotFound, "Not found"))

	diag := resourceFastEdgeBinaryRead(ctx, resourceData, &Config{FastEdgeClient: mockBinClient})

	if diag.HasError() {
		t.Error("unexpected error", diag[0].Summary)
	} else {
		compare(t, resourceData.Id(), "")
	}
}

func TestReadBinaryError(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceFastEdgeBinary().Schema, nil)
	resourceData.SetId("42")

	ctx := contextWithTesting(context.Background(), t)
	ctx = contextWithExpected(ctx, int64(42))
	ctx = contextWithReturn(ctx, sdk.NewJsonHttpResponse(http.StatusBadRequest, `{"error": "Something is wrong"}`))

	diag := resourceFastEdgeBinaryRead(ctx, resourceData, &Config{FastEdgeClient: mockBinClient})

	if diag.HasError() {
		if diag[0].Summary != "Something is wrong" {
			t.Error("unexpected error", diag[0].Summary)
		}
	}
	compare(t, resourceData.Id(), "42") // State is not cleaned
}

func TestUpload(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceFastEdgeBinary().Schema, nil)
	resourceData.Set("filename", os.Args[0])

	ctx := contextWithTesting(context.Background(), t)
	ctx = contextWithReturn(ctx,
		sdk.NewJsonHttpResponse(
			http.StatusOK,
			`{
				"id": 42
			}`,
		),
	)

	diag := resourceFastEdgeBinaryUpload(ctx, resourceData, &Config{FastEdgeClient: mockBinClient})

	if diag.HasError() {
		t.Error("unexpected error", diag[0].Summary)
	} else {
		compare(t, resourceData.Id(), "42")
	}
}
