// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_deployment_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_deployment"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudInferenceDeploymentModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_inference_deployment.CloudInferenceDeploymentModel)(nil)
	schema := cloud_inference_deployment.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
