// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_floating_ip_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_floating_ip"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudFloatingIPsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_floating_ip.CloudFloatingIPsDataSourceModel)(nil)
	schema := cloud_floating_ip.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
