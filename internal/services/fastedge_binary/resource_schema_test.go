// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_binary_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/fastedge_binary"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestFastedgeBinaryModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*fastedge_binary.FastedgeBinaryModel)(nil)
	schema := fastedge_binary.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
