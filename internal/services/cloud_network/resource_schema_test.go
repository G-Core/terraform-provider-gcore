// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_network_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_network"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudNetworkModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_network.CloudNetworkModel)(nil)
	schema := cloud_network.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
