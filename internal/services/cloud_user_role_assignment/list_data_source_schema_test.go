// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_user_role_assignment_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_user_role_assignment"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudUserRoleAssignmentsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_user_role_assignment.CloudUserRoleAssignmentsDataSourceModel)(nil)
	schema := cloud_user_role_assignment.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
