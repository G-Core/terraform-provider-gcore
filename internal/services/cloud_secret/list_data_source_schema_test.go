// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_secret_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_secret"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudSecretsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_secret.CloudSecretsDataSourceModel)(nil)
	schema := cloud_secret.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
