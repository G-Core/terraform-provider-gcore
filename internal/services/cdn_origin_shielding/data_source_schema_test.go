// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_origin_shielding_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_origin_shielding"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNOriginShieldingDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_origin_shielding.CDNOriginShieldingDataSourceModel)(nil)
	schema := cdn_origin_shielding.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
