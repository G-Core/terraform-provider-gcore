// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_shielding_location_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_shielding_location"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNShieldingLocationsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_shielding_location.CDNShieldingLocationsDataSourceModel)(nil)
	schema := cdn_shielding_location.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
