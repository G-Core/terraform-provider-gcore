// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package fastedge_template_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/fastedge_template"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestFastedgeTemplateDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*fastedge_template.FastedgeTemplateDataSourceModel)(nil)
	schema := fastedge_template.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
