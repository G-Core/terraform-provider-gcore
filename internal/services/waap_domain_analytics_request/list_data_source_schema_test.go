// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_analytics_request_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_analytics_request"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapDomainAnalyticsRequestsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain_analytics_request.WaapDomainAnalyticsRequestsDataSourceModel)(nil)
	schema := waap_domain_analytics_request.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
