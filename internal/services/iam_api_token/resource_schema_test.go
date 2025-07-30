// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_api_token_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/iam_api_token"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestIamAPITokenModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*iam_api_token.IamAPITokenModel)(nil)
	schema := iam_api_token.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
