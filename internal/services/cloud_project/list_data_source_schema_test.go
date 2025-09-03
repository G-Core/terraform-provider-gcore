// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_project_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_project"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudProjectsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_project.CloudProjectsDataSourceModel)(nil)
	schema := cloud_project.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
