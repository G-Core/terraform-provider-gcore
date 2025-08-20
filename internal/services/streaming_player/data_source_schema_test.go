// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_player_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_player"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingPlayerDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_player.StreamingPlayerDataSourceModel)(nil)
	schema := streaming_player.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
