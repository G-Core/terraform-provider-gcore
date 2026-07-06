// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/dns_zone"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestDNSZonesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*dns_zone.DNSZonesDataSourceModel)(nil)
	schema := dns_zone.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
