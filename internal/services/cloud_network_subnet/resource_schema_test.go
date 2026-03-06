// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_subnet_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_network_subnet"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudNetworkSubnetModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_network_subnet.CloudNetworkSubnetModel)(nil)
	schema := cloud_network_subnet.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
