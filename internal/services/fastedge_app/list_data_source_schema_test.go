// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_app"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestFastedgeAppsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*fastedge_app.FastedgeAppsDataSourceModel)(nil)
	schema := fastedge_app.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
