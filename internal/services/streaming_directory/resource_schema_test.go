// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package streaming_directory_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/streaming_directory"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestStreamingDirectoryModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*streaming_directory.StreamingDirectoryModel)(nil)
	schema := streaming_directory.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
