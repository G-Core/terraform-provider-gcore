// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_security_group"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudSecurityGroupsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_security_group.CloudSecurityGroupsDataSourceModel)(nil)
	schema := cloud_security_group.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
