// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_user_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/iam_user"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestIamUserModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*iam_user.IamUserModel)(nil)
	schema := iam_user.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
