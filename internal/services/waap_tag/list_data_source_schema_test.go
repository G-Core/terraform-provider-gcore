// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_tag_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_tag"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapTagsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_tag.WaapTagsDataSourceModel)(nil)
	schema := waap_tag.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
