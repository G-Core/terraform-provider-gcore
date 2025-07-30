// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_router_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_network_router"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudNetworkRouterDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_network_router.CloudNetworkRouterDataSourceModel)(nil)
	schema := cloud_network_router.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
