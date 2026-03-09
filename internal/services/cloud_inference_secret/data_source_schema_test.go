// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_secret_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_inference_secret"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudInferenceSecretDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_inference_secret.CloudInferenceSecretDataSourceModel)(nil)
	schema := cloud_inference_secret.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
