// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_application_template_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_application_template"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudInferenceApplicationTemplateDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_inference_application_template.CloudInferenceApplicationTemplateDataSourceModel)(nil)
	schema := cloud_inference_application_template.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
