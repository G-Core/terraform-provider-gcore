// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_model_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_model"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudInferenceModelDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_inference_model.CloudInferenceModelDataSourceModel)(nil)
	schema := cloud_inference_model.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
