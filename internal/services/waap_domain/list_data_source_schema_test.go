// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapDomainsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain.WaapDomainsDataSourceModel)(nil)
	schema := waap_domain.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
