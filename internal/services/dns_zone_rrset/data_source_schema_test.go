// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_rrset_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/dns_zone_rrset"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestDNSZoneRrsetDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*dns_zone_rrset.DNSZoneRrsetDataSourceModel)(nil)
	schema := dns_zone_rrset.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
