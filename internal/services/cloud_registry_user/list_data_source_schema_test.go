// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_user_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_registry_user"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudRegistryUsersDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_registry_user.CloudRegistryUsersDataSourceModel)(nil)
	schema := cloud_registry_user.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
