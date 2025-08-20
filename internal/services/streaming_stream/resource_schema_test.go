// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_stream_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_stream"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingStreamModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_stream.StreamingStreamModel)(nil)
	schema := streaming_stream.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
