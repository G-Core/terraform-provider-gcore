// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_api_token_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/iam_api_token"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestIamAPITokenDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*iam_api_token.IamAPITokenDataSourceModel)(nil)
	schema := iam_api_token.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
