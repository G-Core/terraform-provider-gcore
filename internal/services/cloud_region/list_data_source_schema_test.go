// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_region_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_region"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudRegionsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_region.CloudRegionsDataSourceModel)(nil)
	schema := cloud_region.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
