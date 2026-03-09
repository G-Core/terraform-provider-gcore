// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/waap_domain"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestWaapDomainsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain.WaapDomainsDataSourceModel)(nil)
	schema := waap_domain.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
