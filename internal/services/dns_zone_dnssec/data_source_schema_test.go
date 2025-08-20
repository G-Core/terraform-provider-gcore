// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_dnssec_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/dns_zone_dnssec"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestDNSZoneDnssecDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*dns_zone_dnssec.DNSZoneDnssecDataSourceModel)(nil)
	schema := dns_zone_dnssec.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
