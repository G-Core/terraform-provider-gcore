// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_inference_registry_credential_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_inference_registry_credential"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudInferenceRegistryCredentialsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_inference_registry_credential.CloudInferenceRegistryCredentialsDataSourceModel)(nil)
	schema := cloud_inference_registry_credential.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
