// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_flavor_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_inference_flavor"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudInferenceFlavorDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_inference_flavor.CloudInferenceFlavorDataSourceModel)(nil)
	schema := cloud_inference_flavor.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
