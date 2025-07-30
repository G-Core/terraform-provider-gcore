// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_restream_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_restream"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingRestreamModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_restream.StreamingRestreamModel)(nil)
	schema := streaming_restream.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
