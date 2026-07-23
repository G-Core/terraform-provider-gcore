// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_instance_flavor_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_instance_flavor"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudInstanceFlavorsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_instance_flavor.CloudInstanceFlavorsDataSourceModel)(nil)
	schema := cloud_instance_flavor.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
