// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_discovery_scan_result_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_api_discovery_scan_result"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapDomainAPIDiscoveryScanResultsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain_api_discovery_scan_result.WaapDomainAPIDiscoveryScanResultsDataSourceModel)(nil)
	schema := waap_domain_api_discovery_scan_result.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
