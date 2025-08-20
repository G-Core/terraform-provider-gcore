// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_api_path_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_api_path"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapDomainAPIPathDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain_api_path.WaapDomainAPIPathDataSourceModel)(nil)
	schema := waap_domain_api_path.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
