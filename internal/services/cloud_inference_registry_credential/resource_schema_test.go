// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_registry_credential_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_inference_registry_credential"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudInferenceRegistryCredentialModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_inference_registry_credential.CloudInferenceRegistryCredentialModel)(nil)
	schema := cloud_inference_registry_credential.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
