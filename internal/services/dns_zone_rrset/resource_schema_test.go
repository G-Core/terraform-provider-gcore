// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_rrset_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/dns_zone_rrset"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestDNSZoneRrsetModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*dns_zone_rrset.DNSZoneRrsetModel)(nil)
	schema := dns_zone_rrset.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
