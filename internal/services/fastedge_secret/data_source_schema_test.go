// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_secret_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_secret"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestFastedgeSecretDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*fastedge_secret.FastedgeSecretDataSourceModel)(nil)
	schema := fastedge_secret.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
