// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_registry"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudRegistriesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_registry.CloudRegistriesDataSourceModel)(nil)
	schema := cloud_registry.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
