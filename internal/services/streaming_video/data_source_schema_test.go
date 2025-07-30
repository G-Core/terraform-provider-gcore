// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_video_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_video"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingVideoDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_video.StreamingVideoDataSourceModel)(nil)
	schema := streaming_video.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
