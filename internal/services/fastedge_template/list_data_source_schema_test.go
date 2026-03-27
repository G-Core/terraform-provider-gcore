// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_template_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/fastedge_template"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestFastedgeTemplatesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*fastedge_template.FastedgeTemplatesDataSourceModel)(nil)
	schema := fastedge_template.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
