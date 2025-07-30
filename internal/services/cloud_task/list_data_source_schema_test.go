// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_task_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_task"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudTasksDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_task.CloudTasksDataSourceModel)(nil)
	schema := cloud_task.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
