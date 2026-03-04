// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_network_mapping_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/dns_network_mapping"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestDNSNetworkMappingModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*dns_network_mapping.DNSNetworkMappingModel)(nil)
	schema := dns_network_mapping.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
