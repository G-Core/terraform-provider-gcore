// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_organization_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_organization"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapOrganizationsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_organization.WaapOrganizationsDataSourceModel)(nil)
	schema := waap_organization.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
