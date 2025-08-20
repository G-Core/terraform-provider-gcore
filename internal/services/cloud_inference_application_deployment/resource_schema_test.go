// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_application_deployment_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_application_deployment"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudInferenceApplicationDeploymentModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_inference_application_deployment.CloudInferenceApplicationDeploymentModel)(nil)
	schema := cloud_inference_application_deployment.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
