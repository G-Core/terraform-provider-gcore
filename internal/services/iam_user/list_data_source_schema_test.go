// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package iam_user_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/iam_user"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestIamUsersDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*iam_user.IamUsersDataSourceModel)(nil)
	schema := iam_user.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
