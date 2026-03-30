// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_server_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_baremetal_server"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudBaremetalServerModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_baremetal_server.CloudBaremetalServerModel)(nil)
	schema := cloud_baremetal_server.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
