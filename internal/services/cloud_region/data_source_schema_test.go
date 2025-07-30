// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_region_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_region"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudRegionDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_region.CloudRegionDataSourceModel)(nil)
	schema := cloud_region.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
