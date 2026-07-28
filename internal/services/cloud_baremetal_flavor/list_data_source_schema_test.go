// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_flavor_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_baremetal_flavor"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudBaremetalFlavorsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_baremetal_flavor.CloudBaremetalFlavorsDataSourceModel)(nil)
	schema := cloud_baremetal_flavor.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
