// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_ai_task_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_ai_task"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingAITaskModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_ai_task.StreamingAITaskModel)(nil)
	schema := streaming_ai_task.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
