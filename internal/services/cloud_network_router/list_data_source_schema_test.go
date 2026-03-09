// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_router_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_network_router"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudNetworkRoutersDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_network_router.CloudNetworkRoutersDataSourceModel)(nil)
	schema := cloud_network_router.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
