package gcore

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	sdk "github.com/G-Core/FastEdge-client-sdk-go"
)

type mock struct {
	sdk.ClientSDKmock
}

func (mock) GetBinary(ctx context.Context, id int64, reqEditors ...sdk.RequestEditorFn) (*http.Response, error) {
	return sdk.NewJsonHttpResponse(http.StatusOK, `{"id": 1, "checksum":"1234567890"}`), nil
}

var mockClient *sdk.ClientWithResponses = &sdk.ClientWithResponses{ClientInterface: &mock{}}

func TestReadBinary(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceFastEdgeBinary().Schema, nil)
	resourceData.SetId("1")

	diag := resourceFastEdgeBinaryRead(context.Background(), resourceData, &Config{FastEdgeClient: mockClient})

	if diag.HasError() {
		t.Error("unexpected error", diag[0].Summary)
	} else if resourceData.Get("checksum").(string) != "1234567890" {
		t.Error("unexpected checksum", resourceData.Get("checksum"))
	}
}
