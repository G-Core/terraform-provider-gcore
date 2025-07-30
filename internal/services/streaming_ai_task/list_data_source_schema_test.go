// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_ai_task_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_ai_task"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingAITasksDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_ai_task.StreamingAITasksDataSourceModel)(nil)
	schema := streaming_ai_task.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
