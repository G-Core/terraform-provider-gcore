// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_reserved_fixed_ip_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_reserved_fixed_ip"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudReservedFixedIPDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_reserved_fixed_ip.CloudReservedFixedIPDataSourceModel)(nil)
	schema := cloud_reserved_fixed_ip.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
