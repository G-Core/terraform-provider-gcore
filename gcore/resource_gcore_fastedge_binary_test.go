package gcore

import (
	"context"
	"net/http"
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
