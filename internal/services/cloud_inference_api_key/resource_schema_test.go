// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_api_key_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_api_key"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudInferenceAPIKeyModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_inference_api_key.CloudInferenceAPIKeyModel)(nil)
	schema := cloud_inference_api_key.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
