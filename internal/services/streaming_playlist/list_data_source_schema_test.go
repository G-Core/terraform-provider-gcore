// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_playlist_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_playlist"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingPlaylistsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_playlist.StreamingPlaylistsDataSourceModel)(nil)
	schema := streaming_playlist.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
