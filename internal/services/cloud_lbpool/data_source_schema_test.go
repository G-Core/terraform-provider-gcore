// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_lbpool_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_lbpool"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudLbpoolDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_lbpool.CloudLbpoolDataSourceModel)(nil)
	schema := cloud_lbpool.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
