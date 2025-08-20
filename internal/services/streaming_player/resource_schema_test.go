// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_player_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_player"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingPlayerModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_player.StreamingPlayerModel)(nil)
	schema := streaming_player.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
