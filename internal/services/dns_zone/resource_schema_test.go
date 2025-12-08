// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_zone_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/dns_zone"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestDNSZoneModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*dns_zone.DNSZoneModel)(nil)
	schema := dns_zone.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
