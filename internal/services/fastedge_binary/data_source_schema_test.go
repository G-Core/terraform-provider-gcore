// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_binary_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_binary"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestFastedgeBinaryDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*fastedge_binary.FastedgeBinaryDataSourceModel)(nil)
	schema := fastedge_binary.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
