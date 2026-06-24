// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_client_config_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_client_config"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNClientConfigDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_client_config.CDNClientConfigDataSourceModel)(nil)
	schema := cdn_client_config.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
