// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_broadcast_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_broadcast"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingBroadcastsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_broadcast.StreamingBroadcastsDataSourceModel)(nil)
	schema := streaming_broadcast.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
