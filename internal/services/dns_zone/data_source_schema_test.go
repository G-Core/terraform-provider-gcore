// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/dns_zone"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestDNSZoneDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*dns_zone.DNSZoneDataSourceModel)(nil)
	schema := dns_zone.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
