// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_stream_overlay_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_stream_overlay"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingStreamOverlayModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_stream_overlay.StreamingStreamOverlayModel)(nil)
	schema := streaming_stream_overlay.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
