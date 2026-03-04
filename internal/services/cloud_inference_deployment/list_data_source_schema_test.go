// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_deployment_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_inference_deployment"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudInferenceDeploymentsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_inference_deployment.CloudInferenceDeploymentsDataSourceModel)(nil)
	schema := cloud_inference_deployment.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
