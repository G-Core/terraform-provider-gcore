// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_registry"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudRegistryDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_registry.CloudRegistryDataSourceModel)(nil)
	schema := cloud_registry.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
