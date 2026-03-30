// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_baremetal_server_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_baremetal_server"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudBaremetalServerDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_baremetal_server.CloudBaremetalServerDataSourceModel)(nil)
	schema := cloud_baremetal_server.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
