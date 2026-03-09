// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_reserved_fixed_ip_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_reserved_fixed_ip"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudReservedFixedIPsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_reserved_fixed_ip.CloudReservedFixedIPsDataSourceModel)(nil)
	schema := cloud_reserved_fixed_ip.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
