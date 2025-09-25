// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_subnet_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_network_subnet"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudNetworkSubnetDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_network_subnet.CloudNetworkSubnetDataSourceModel)(nil)
	schema := cloud_network_subnet.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
