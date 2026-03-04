// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package dns_network_mapping_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/dns_network_mapping"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestDNSNetworkMappingDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*dns_network_mapping.DNSNetworkMappingDataSourceModel)(nil)
	schema := dns_network_mapping.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
