// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_app_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_app"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestFastedgeAppDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*fastedge_app.FastedgeAppDataSourceModel)(nil)
	schema := fastedge_app.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
