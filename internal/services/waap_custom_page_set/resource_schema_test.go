// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_custom_page_set_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_custom_page_set"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapCustomPageSetModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_custom_page_set.WaapCustomPageSetModel)(nil)
	schema := waap_custom_page_set.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
