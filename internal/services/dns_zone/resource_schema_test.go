// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/dns_zone"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestDNSZoneModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*dns_zone.DNSZoneModel)(nil)
	schema := dns_zone.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
