// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_registry_repository_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_registry_repository"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudRegistryRepositoriesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_registry_repository.CloudRegistryRepositoriesDataSourceModel)(nil)
	schema := cloud_registry_repository.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
