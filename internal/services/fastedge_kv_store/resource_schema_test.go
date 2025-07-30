// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_kv_store_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_kv_store"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestFastedgeKvStoreModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*fastedge_kv_store.FastedgeKvStoreModel)(nil)
	schema := fastedge_kv_store.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
