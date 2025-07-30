// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_video_subtitle_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_video_subtitle"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingVideoSubtitleDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_video_subtitle.StreamingVideoSubtitleDataSourceModel)(nil)
	schema := streaming_video_subtitle.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
